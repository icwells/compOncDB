// Compares parse output with nlp predictions

package predictor

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/go-tools/strarray"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type predictor struct {
	col		  string
	columns   []string
	diagnosis bool
	dir		  string
	infile	  string
	logger	  *log.Logger
	mass	  string
	mindiag	  float64
	minmass   float64
	neoplasia bool
	outfile	  string
	records	  *dataframe.Dataframe
	script	  string
}

func newPredictor(infile string, neoplasia, diagnosis bool) *predictor {
	// Return initialized struct
	var err error
	p := new(predictor)
	p.setMode(neoplasia, diagnosis)
	p.col = "Comments"
	p.dir = path.Join(iotools.GetGOPATH(), "src/github.com/icwells/compOncDB/scripts/nlpModel/")
	p.infile = "nlpInput.csv"
	p.logger = codbutils.GetLogger()
	p.mass = "Masspresent"
	p.mindiag = 0.995
	p.minmass = 0.5
	p.outfile = "nlpOutput.csv"
	if p.records, err = dataframe.FromFile(infile, 0); err != nil {
		p.logger.Fatal(err)
	}
	p.script = "nlpModel.py"
	p.removeNA()
	p.alterColumns()
	return p
}

func (p *predictor) setMode(neoplasia, diagnosis bool) {
	// Determines whether to run neoplasia comparison, diagnosis comparison, or both
	p.columns = []string{"MassVerified", "TypeVerified", "LocationVerified"}
	p.neoplasia = neoplasia
	p.diagnosis = diagnosis
	if !p.neoplasia && !p.diagnosis {
		p.neoplasia = true
		p.diagnosis = true
	} else if p.neoplasia && !p.diagnosis {
		p.columns = p.columns[:1]
	} else if !p.neoplasia && p.diagnosis {
		p.columns = p.columns[1:]
	}
}

func (p *predictor) removeNA() {
	// Removes rows where comments == NA since no prediction can be made
	var rm []string
	p.logger.Println("Removing NA comments...")
	for i := range p.records.Iterate() {
		if comments, _ := i.GetCell("Comments"); comments == "NA" {
			rm = append(rm, i.Name)
		}
	}
	for _, i := range rm {
		p.records.DeleteRow(i)
	}
}

func (p *predictor) alterColumns() {
	// Removes extra columns and adds columns for verifications
	columns := []string{"ID", "Comments", "Masspresent", "Type", "Location"}
	if !p.diagnosis {
		columns = columns[:3]
	}
	for k := range p.records.Header {
		if !strarray.InSliceStr(columns, k) {
			p.records.DeleteColumn(k)
		}
	}
	for _, i := range p.columns {
		p.records.AddColumn(i, "")
	}
}

func (p *predictor) callScript(diagnosis bool) {
	// Configures command for python script and calls
	var cmd *exec.Cmd
	p.logger.Println("Calling prediction script...")
	dir, _ := os.Getwd()
	os.Chdir(p.dir)
	infile := fmt.Sprintf("-i%s", p.infile)
	outfile := fmt.Sprintf("-o%s", p.outfile)
	if diagnosis {
		cmd = exec.Command("python", p.script, "--diagnosis", infile, outfile)
	} else {
		cmd = exec.Command("python", p.script, infile, outfile)
	}
	if err := cmd.Run(); err != nil {
		p.logger.Fatalf("Prediction script failed. %v\n", err)
	}
	p.logger.Println("Prediction script complete.")
	os.Chdir(dir)
}

func (p *predictor) writeInfile(diagnosis bool) {
	// Writes records to input file for script
	p.logger.Println("Writing input file for prediction script...")
	out := iotools.CreateFile(path.Join(p.dir, p.infile))
	defer out.Close()
	for i := range p.records.Iterate() {
		mp, err := i.GetCellInt(p.mass)
		if err != nil {
			p.logger.Fatal(err)
		}
		if !diagnosis || mp == 1 {
			// Only examine cancer records if diagnosis is true
			v, err := i.GetCell(p.col)
			if err == nil {
				out.WriteString(fmt.Sprintf("%s,%s\n", i.Name, v))
			} else {
				p.logger.Fatal(err)
			}
		}
	}
}

func (p *predictor) compareNeopasia() {
	// Compares neoplasia results to parse output
	p.logger.Println("Comparing neoplasia results...")
	reader, _ := iotools.YieldFile(path.Join(p.dir, p.outfile), false)
	for i := range reader {
		id := i[0]
		if score, err := strconv.ParseFloat(i[2], 64); err == nil {
			mp, _ := p.records.GetCellInt(id, p.mass)
			if score >= p.minmass && mp != 1 {
				p.records.UpdateCell(id, p.columns[0], "1")
			} else if score <= 1 - p.minmass && mp != 0 {
				p.records.UpdateCell(id, p.columns[0], "0")
			}
		}

	}
}

func (p *predictor) predictMass() {
	// Calls nlp model to predict mass
	p.logger.Println("Predicting neoplasia diagnoses...")
	p.writeInfile(false)
	p.callScript(false)
	p.compareNeopasia()
}

func (p *predictor) compareDiagnoses() {
	// Compares type and location results to parse output
	p.logger.Println("Comparing type and location results...")
	reader, header := iotools.YieldFile(path.Join(p.dir, p.outfile), true)
	for i := range reader {
		id := i[header["ID"]]
		typ, _ := p.records.GetCell(id, "Type")
		loc, _ := p.records.GetCell(id, "Location")
		if score, err := strconv.ParseFloat(i[header["Lscore"]], 64); err == nil {
			if score >= p.mindiag && strings.ToLower(loc) != i[header["Location"]] {
				p.records.UpdateCell(id, p.columns[2], i[header["Location"]])
			}
		}
		if score, err := strconv.ParseFloat(i[header["Tscore"]], 64); err == nil {
			if score >= p.mindiag && strings.ToLower(typ) != i[header["Type"]] {
				p.records.UpdateCell(id, p.columns[1], i[header["Type"]])
			}
		}
	}
}

func (p *predictor) predictDiagnoses() {
	// Calls nlp model to predict type and location
	p.logger.Println("Predicting type and location diagnoses...")
	p.writeInfile(true)
	p.callScript(true)
	p.compareDiagnoses()
}

func (p *predictor) removePasses() {
	// Removes rows which don't need to be  updated
	p.logger.Println("Removing approved records...")
	var rm []string
	var count int
	for i := range p.records.Iterate() {
		var mp, t, l string
		if p.neoplasia {
			mp, _ = i.GetCell("MassVerified")
		}
		if p.diagnosis {
			t, _ = i.GetCell("TypeVerified")
			l, _ = i.GetCell("LocationVerified")
		}
		if p.neoplasia && !p.diagnosis && mp == "" {
			rm = append(rm, i.Name)
		} else if !p.neoplasia && t == "" && l == "" {
			rm = append(rm, i.Name)
		} else if mp == "" && t == "" && l == "" {
			rm = append(rm, i.Name)
		}
	}
	for _, i := range rm {
		p.records.DeleteRow(i)
		fmt.Printf("\tRemoved %d of %d verified records.\r", count, len(rm))
		count++
	}
	fmt.Println()
	p.logger.Printf("Identified %d records to review...", p.records.Length())
}

func (p *predictor) cleanup() {
	// Removes infiile and outfile after use
	for _, i := range []string{p.infile, p.outfile} {
		f := path.Join(p.dir, i)
		if iotools.Exists(f) {
			os.Remove(f)
		}
	}
}

func ComparePredictions(infile string, neoplasia, diagnosis bool) *dataframe.Dataframe {
	// Compares parse output with nlp predictions
	p := newPredictor(infile, neoplasia, diagnosis)
	defer p.cleanup()
	if p.neoplasia {
		p.predictMass()
	}
	if p.diagnosis {
		p.predictDiagnoses()
	}
	p.removePasses()
	return p.records
}

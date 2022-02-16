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
	col		string
	columns []string
	dir		string
	infile	string
	logger	*log.Logger
	mass	string
	min		float64
	outfile	string
	records	*dataframe.Dataframe
	script	string
}

func newPredictor(infile string) *predictor {
	// Return initialized struct
	var err error
	p := new(predictor)
	p.col = "Comments"
	p.columns = []string{"MassVerified", "TypeVerified", "LocationVerified"}
	p.dir = path.Join(iotools.GetGOPATH(), "src/github.com/icwells/compOncDB/scripts/nlpModel/")
	p.infile = "nlpInput.csv"
	p.logger = codbutils.GetLogger()
	p.mass = "Masspresent"
	p.min = 0.9
	p.outfile = "nlpOutput.csv"
	if p.records, err = dataframe.FromFile(infile, 0); err != nil {
		p.logger.Fatal(err)
	}
	p.script = "nlpModel.py"
	for k := range p.records.Header {
		if !strarray.InSliceStr([]string{"ID", "Comments", "Masspresent", "Type", "Location"}, k) {
			p.records.DeleteColumn(k)
		}
	}
	for _, i := range p.columns {
		p.records.AddColumn(i, "")
	}
	return p
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
		verified := "0"
		id := i[0]
		if score, err := strconv.ParseFloat(i[2], 64); err == nil {
			mp, _ := p.records.GetCellInt(id, p.mass)
			if score >= p.min && mp == 1 {
				verified = "1"
			} else if score <= 1 - p.min && mp == 0 {
				verified = "1"
			}
		}
		p.records.UpdateCell(id, p.columns[0], verified)
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
			if score >= p.min && strings.ToLower(loc) != i[header["Location"]] {
				p.records.UpdateCell(id, p.columns[2], i[header["Location"]])
			}
		}
		if score, err := strconv.ParseFloat(i[header["Tscore"]], 64); err == nil {
			if score >= p.min && strings.ToLower(typ) != i[header["Type"]] {
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
		mp, _ := i.GetCellInt("MassVerified")
		t, _ := i.GetCell("TypeVerified")
		l, _ := i.GetCell("LocationVerified")
		if mp == 1 && t == "" && l == "" {
			rm = append(rm, i.Name)
		}
	}
	//p.logger.Printf("Removing %d records...", len(rm))
	for _, i := range rm {
		if err := p.records.DeleteRow(i); err != nil {
			panic(err)
		}
		fmt.Printf("\tRemoved %d of %d verified records.\r", count, len(rm))
		count++
	}
	fmt.Println()
	p.logger.Printf("Identified %d records to review...", p.records.Length())
}

func (p *predictor) cleanup() {
	// Removes infiile and outfile after use
	os.Remove(path.Join(p.dir, p.infile))
	os.Remove(path.Join(p.dir, p.outfile))
}

func ComparePredictions(infile string) *dataframe.Dataframe {
	// Compares parse output with nlp predictions
	p := newPredictor(infile)
	defer p.cleanup()
	p.predictMass()
	p.predictDiagnoses()
	p.removePasses()
	return p.records
}

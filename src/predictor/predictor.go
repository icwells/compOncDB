// Defines predictor struct and common methods

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
)

type predictor struct {
	col         string
	columns     []string
	diagnosis   bool
	dir         string
	hcol		string
	hyperplasia string
	infile      string
	lcol        string
	logger      *log.Logger
	mass        string
	mcol        string
	mindiag     float64
	minmass     float64
	neoplasia   bool
	outfile     string
	records     *dataframe.Dataframe
	results     *dataframe.Dataframe
	script      string
	tcol        string
}

func newPredictor(infile string, neoplasia, diagnosis bool) *predictor {
	// Return initialized struct
	var err error
	p := new(predictor)
	p.setMode(neoplasia, diagnosis)
	p.col = "Comments"
	p.dir = path.Join(iotools.GetGOPATH(), "src/github.com/icwells/compOncDB/scripts/nlpModel/")

	p.hyperplasia = "Hyperplasia"
	p.infile = "nlpInput.csv"
	p.logger = codbutils.GetLogger()
	p.mass = "Masspresent"
	p.mindiag = 0.99
	p.minmass = 0.99
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
	p.hcol = "HyperplasiaVerified"
	p.lcol = "LocationVerified"
	p.mcol = "MassVerified"
	p.tcol = "TypeVerified"
	p.columns = []string{p.mcol, p.hcol, p.tcol, p.lcol}
	p.neoplasia = neoplasia
	p.diagnosis = diagnosis
	if !p.neoplasia && !p.diagnosis {
		p.neoplasia = true
		p.diagnosis = true
	} else if p.neoplasia && !p.diagnosis {
		p.columns = p.columns[:2]
	} else if !p.neoplasia && p.diagnosis {
		p.columns = p.columns[2:]
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
	columns := []string{"ID", "Comments", p.mass, p.hyperplasia, "Type", "Location"}
	if !p.diagnosis {
		columns = columns[:4]
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
	var count int
	p.logger.Println("Writing input file for prediction script...")
	out := iotools.CreateFile(path.Join(p.dir, p.infile))
	defer out.Close()
	for i := range p.records.Iterate() {
		mp, _ := i.GetCellInt(p.mass)
		hyp, _ := i.GetCellInt(p.hyperplasia)
		if !diagnosis || mp == 1 || hyp == 1 {
			// Only examine cancer records if diagnosis is true
			v, err := i.GetCell(p.col)
			if err == nil {
				out.WriteString(fmt.Sprintf("%s,%s\n", i.Name, v))
				count++
			} else {
				p.logger.Fatal(err)
			}
		}
	}
	p.logger.Printf("Identified %d records for verification.", count)
}

func (p *predictor) cleanup() {
	// Removes infile and outfile after use
	for _, i := range []string{p.infile, p.outfile} {
		f := path.Join(p.dir, i)
		if iotools.Exists(f) {
			os.Remove(f)
		}
	}
}

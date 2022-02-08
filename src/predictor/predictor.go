// Compares parse output with nlp predictions

package predictor

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/go-tools/iotools"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
)

type predictor struct {
	col		string
	columns []string
	infile	string
	logger	*log.Logger
	mass	string
	outfile	string
	records	*dataframe.Dataframe
	script	string
}

func newPredictor(infile string) *predictor {
	// Return initialized struct
	var err error
	p := new(predictor)
	p.col = "Comments"
	p.columns = []string{"MassVerified"}//, "TypeVerified", "LocationVerified"}
	p.infile = "/tmp/nlpInput.csv"
	p.logger = codbutils.GetLogger()
	p.mass = "Masspresent"
	p.outfile = "/tmp/nlpOutput.csv"
	if p.records, err = dataframe.FromFile(infile, 0); err != nil {
		p.logger.Fatal(err)
	}
	p.script = path.Join(iotools.GetGOPATH(), "src/github.com/icwells/compOncDB/scripts/nlpModel/nlpModel.py")
	for _, i := range p.columns {
		p.records.AddColumn(i, "")
	}
	return p
}

func (p *predictor) callScript(diagnosis bool) {
	// Configures command for python script and calls
	var cmd *exec.Cmd
	p.logger.Println("Calling prediction script...")
	infile := fmt.Sprintf("-i%s", p.infile)
	outfile := fmt.Sprintf("-o%s", p.outfile)
	if diagnosis {
		cmd = exec.Command("python", p.script, "--diagnosis", infile, outfile)
	} else {
		cmd = exec.Command("python", p.script, infile, outfile)
	}
	err := cmd.Run()
	if err == nil {
		p.logger.Println("Prediction script complete.")
	} else {
		p.logger.Fatalf("Prediction script failed. %v\n", err)
	}
}

func (p *predictor) writeInfile(diagnosis bool) {
	// Writes records to input file for script
	p.logger.Println("Writing input file for prediction script...")
	out := iotools.CreateFile(p.infile)
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
	reader, _ := iotools.YieldFile(p.outfile, false)
	for i := range reader {
		verified := "0"
		id := i[0]
		if score, err := strconv.ParseFloat(i[1], 64); err == nil {
			mp, _ := p.records.GetCellInt(id, p.columns[0])
			if score >= 0.8 && mp == 1 {
				verified = "1"
			} else if score <= 0.2 && mp == 1 {
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

/*func (p *predictor) compareDiagnoses() {
	// Compares type and location results to parse output
	p.logger.Println("Comparing type and location results...")
	reader, header := iotools.YieldFile(p.outfile, true)
	for i := range reader {
		id := i[0]
		
	}
}

func (p *predictor) predictDiagnoses() {
	// Calls nlp model to predict type and location
	p.logger.Println("Predicting type and location diagnoses...")
	p.writeInfile(true)
	p.callScript(true)
	p.compareDiagnoses()
}*/

func (p *predictor) cleanup() {
	// Removes infiile and outfile after use
	os.Remove(p.infile)
	os.Remove(p.outfile)
}

func ComparePredictions(infile string) *dataframe.Dataframe {
	// Compares parse output with nlp predictions
	p := newPredictor(infile)
	defer p.cleanup()
	p.predictMass()
	//p.predictDiagnoses()
	return p.records
}

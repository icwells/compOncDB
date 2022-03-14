// Compares parse output with nlp predictions

package predictor

import (
	"fmt"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/go-tools/iotools"
	"path"
	"strconv"
	"strings"
)

func (p *predictor) compareNeopasia() {
	// Compares neoplasia results to parse output
	p.logger.Println("Comparing neoplasia results...")
	reader, _ := iotools.YieldFile(path.Join(p.dir, p.outfile), false)
	for i := range reader {
		id := i[0]
		if score, err := strconv.ParseFloat(i[2], 64); err == nil {
			mp, _ := p.records.GetCellInt(id, p.mass)
			if score >= p.minmass && mp != 1 {
				p.records.UpdateCell(id, p.mcol, "1")
			} else if score <= 1-p.minmass && mp != 0 {
				p.records.UpdateCell(id, p.mcol, "0")
			}
		}
		if score, err := strconv.ParseFloat(i[3], 64); err == nil {
			hyp, _ := p.records.GetCellInt(id, p.hyperplasia)
			if score >= p.minmass && hyp != 1 {
				p.records.UpdateCell(id, p.hcol, "1")
			} else if score <= 1-p.minmass && hyp != 0 {
				p.records.UpdateCell(id, p.hcol, "0")
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

func (p *predictor) diagnosisContains(diag, val string) bool {
	// Returns true if diagnosis value has already been identified
	d := ";"
	val = strings.ToLower(val)
	if strings.Contains(diag, d) {
		for _, i := range strings.Split(diag, d) {
			if i == val {
				return true
			}
		}
	} else if diag == val {
		return true
	}
	return false
}

func (p *predictor) compareDiagnoses() {
	// Compares type and location results to parse output
	p.logger.Println("Comparing type and location results...")
	reader, header := iotools.YieldFile(path.Join(p.dir, p.outfile), true)
	for i := range reader {
		id := i[header["ID"]]
		typ, _ := p.records.GetCell(id, "Type")
		loc, _ := p.records.GetCell(id, "Location")
		if score, err := strconv.ParseFloat(i[header["Tscore"]], 64); err == nil {
			if score >= p.mindiag && !p.diagnosisContains(typ, i[header["Type"]]) {
				p.records.UpdateCell(id, p.tcol, i[header["Type"]])
			}
		}
		if score, err := strconv.ParseFloat(i[header["Lscore"]], 64); err == nil {
			if score >= p.mindiag && !p.diagnosisContains(loc, i[header["Location"]]) {
				p.records.UpdateCell(id, p.lcol, i[header["Location"]])
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

func (p *predictor) subsetUnverified() {
	// Removes rows which don't need to be  updated
	p.logger.Println("Removing approved records...")
	p.results = p.records.Clone()
	for i := range p.records.Iterate() {
		var save bool
		var mp, t, l string
		if p.neoplasia {
			mp, _ = i.GetCell(p.mcol)
		}
		if p.diagnosis {
			t, _ = i.GetCell(p.tcol)
			l, _ = i.GetCell(p.lcol)
		}
		if p.neoplasia && !p.diagnosis && mp != "" {
			save = true
		} else if !p.neoplasia && p.diagnosis {
			if t != "" || l != "" {
				save = true
			}
		} else if mp != "" || t != "" || l != "" {
			save = true
		}
		if save {
			p.results.AddRow(i.ToSlice())
		}
	}
	p.logger.Printf("Identified %d records to review.", p.results.Length())
}

func ComparePredictions(infile string, neoplasia, diagnosis bool) *dataframe.Dataframe {
	// Compares parse output with nlp predictions
	fmt.Println()
	p := newPredictor(infile, neoplasia, diagnosis)
	defer p.cleanup()
	if p.neoplasia {
		p.predictMass()
	}
	if p.diagnosis {
		p.predictDiagnoses()
	}
	p.subsetUnverified()
	return p.results
}

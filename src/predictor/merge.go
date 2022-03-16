// Merges verification output with parse output

package predictor

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/diagnoses"
	"github.com/icwells/go-tools/dataframe"
	"log"
)

type merger struct {
	col      []string
	logger   *log.Logger
	records  *dataframe.Dataframe
	tissues  map[string]string
	verified *dataframe.Dataframe
}

func newMerger(infile, parsefile string) *merger {
	// Returns initialized struct
	var err error
	m := new(merger)
	m.logger = codbutils.GetLogger()
	m.logger.Println("Reading input files...")
	if m.records, err = dataframe.FromFile(parsefile, 0); err != nil {
		m.logger.Fatal(err)
	}
	matcher := diagnoses.NewMatcher(m.logger)
	m.tissues = matcher.GetTissues()
	if m.verified, err = dataframe.FromFile(infile, 0); err != nil {
		m.logger.Fatal(err)
	}
	for _, i := range []string{"Masspresent", "Hyperplasia", "Type", "Location"} {
		// Only add columns present in verification file
		if _, ex := m.verified.Header[i]; ex {
			m.col = append(m.col, i)
		}
	}
	return m
}

func (m *merger) mergeRecords() {
	// Updates records with masspresent, types, and locations from verification file
	m.logger.Println("Merging files...")
	for row := range m.verified.Iterate() {
		for _, i := range m.col {
			val, _ := row.GetCell(i)
			m.records.UpdateCell(row.Name, i, val)
			if i == "Type" && val != "" && val != "NA" {
				if mp, _ := m.records.GetCellInt(row.Name, "Masspresent"); mp != 1 {
					m.records.UpdateCell(row.Name, "Masspresent", "1")
				}
			} else if i == "Location" {
				if val == "" || val == "NA" {
					m.records.UpdateCell(row.Name, "Tissue", "NA")
				} else if tis, ex := m.tissues[val]; ex {
					if tissue, _ := m.records.GetCell(row.Name, "Tissue"); tissue != tis {
						m.records.UpdateCell(row.Name, "Tissue", tis)
					}
				}
			}
		}
	}
}

func MergePredictions(infile, parsefile string) *dataframe.Dataframe {
	// Merges verification output with parse output
	m := newMerger(infile, parsefile)
	m.mergeRecords()
	return m.records
}

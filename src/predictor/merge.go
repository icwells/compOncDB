// Merges verification output with parse output

package predictor

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/go-tools/dataframe"
	"log"
)

type merger struct {
	col      []string
	logger	 *log.Logger
	records	 *dataframe.Dataframe
	verified *dataframe.Dataframe
}

func newMerger(infile, parsefile string) *merger {
	// Returns initialized struct
	var err error
	m := new(merger)
	m.col = []string{"Masspresent", "Type", "Location"}
	m.logger = codbutils.GetLogger()
	m.logger.Println("Reading input files...")
	if m.records, err = dataframe.FromFile(parsefile, 0); err != nil {
		m.logger.Fatal(err)
	}
	if m.verified, err = dataframe.FromFile(infile, 0); err != nil {
		m.logger.Fatal(err)
	}
	return m
}

func (m *merger) mergeRecords() {
	// Updates records with masspresent, types, and locations from verification file
	m.logger.Println("Merging files...")
	for row := range m.records.Iterate() {
		for _, i := range m.col {
			val, _ := row.GetCell(i)
			m.records.UpdateCell(row.Name, i, val)
		}
	}
}

func MergePredictions(infile, parsefile string) *dataframe.Dataframe {
	// Merges verification output with parse output
	m := newMerger(infile, parsefile)
	m.mergeRecords()
	return m.records
}

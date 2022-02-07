// Compares parse output with nlp predictions

package parse

import (

	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/go-tools/dataframe"
	"log"
)

type predictor struct {
	log		*log.Logger
	records	*dataframe.Dataframe
	rows	[][]string
}

func newPredictor(infile string) *predictor {
	// Return initialized struct
	var err error
	p := new(predictor)
	
	p.logger = codbutils.GetLogger()
	p.outfile = outfile
	if p.records, err = database.FromFile(infile, 0); err != nil {
		logger.Fatal(err)
	}
	for _, i := range []string{"MassVerified", "TypeVerified", "LocationVerified"} {
		p.records.AddRow(i, "NA")
	}
	return p
}

func predictMass() {
	// Calls nlp model to predict mass
}

func ComparePredictions(infile string) *dataframe.Dataframe {
	// Compares parse output with nlp predictions
	p := newPredictor(infile)
	p.predictMass()
	return p.records
}

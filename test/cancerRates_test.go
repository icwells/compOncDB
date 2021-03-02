// Tests cancer rate calculations

package main

import (
	"flag"
	"github.com/icwells/compOncDB/src/cancerrates"
	"github.com/icwells/compOncDB/src/codbutils"
	"testing"
)

func TestCancerRates(t *testing.T) {
	// Tests taxonomy search output
	db := connectToDatabase()
	rates := cancerrates.GetCancerRates(db, 1, 0, false, false, false, "", "")
	compareTables(t, "Cancer Rates", getExpectedRates(), rates)
}

func TestNecropsies(t *testing.T) {
	// Tests necropsy filtering with full database
	flag.Parse()
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), *password)
	for _, val := range []int{-1, 1} {
		rates := cancerrates.GetCancerRates(db, 1, val, false, false, false, "", "")
		if rates.Length() == 0 {
			t.Error("Necropsy dataframe length is 0.")
			break
		}
		for idx := range rates.Rows {
			total, err := rates.GetCellInt(idx, "RecordsWithDenominators")
			nec, er := rates.GetCellInt(idx, "Necropsies")
			if err != nil {
				t.Error(err)
				break
			} else if er != nil {
				t.Error(er)
				break
			}
			if val == 1 && total != nec {
				t.Errorf("Total records %d does not equal necropsies: %d.", total, nec)
				break
			} else if nec != 0 {
				t.Errorf("%d necropsies found in non-necropsies records.", nec)
				break
			}
		}
	}
}

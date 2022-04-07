// Tests cancer rate calculations

package main

import (
	"flag"
	"github.com/icwells/compOncDB/src/cancerrates"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"testing"
)

func getRates(db *dbIO.DBIO, nec int) *dataframe.Dataframe {
	// Initialized cancerrates struct and returns dataframe
	c := cancerrates.NewCancerRates(db, 1, false, "", "")
	c.SearchSettings(nec, false, false, "all")
	c.OutputSettings(true, false, true, true)
	ret, _ := c.GetCancerRates("")
	return ret
}

func TestCancerRates(t *testing.T) {
	// Tests taxonomy search output
	db := connectToDatabase()
	rates := getRates(db, 0)
	compareTables(t, "Cancer Rates", getExpectedRates(), rates)
}

func TestPrevlenceTotals(t *testing.T) {
	// Compares total records with RecordsWithDenominators
	var count int
	flag.Parse()
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), *password)
	rates := getRates(db, 0)
	for idx := range rates.Rows {
		total, _ := rates.GetCellInt(idx, "TotalRecords")
		den, _ := rates.GetCellInt(idx, "RecordsWithDenominators")
		if den > total {
			count++
		}
	}
	if count > 0 {
		t.Errorf("Found %d species with more RecordsWithDenominators than TotalRecords.", count)
	}
}

func TestNecropsies(t *testing.T) {
	// Tests necropsy filtering with full database
	flag.Parse()
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), *password)
	for _, val := range []int{-1, 1} {
		var count int
		name := "necropsies"
		if val == -1 {
			name = "non-" + name
		}
		rates := getRates(db, val)
		if rates.Length() == 0 {
			t.Error("Necropsy dataframe length is 0.")
			break
		}
		for idx := range rates.Rows {
			total, _ := rates.GetCellInt(idx, "RecordsWithDenominators")
			nec, _ := rates.GetCellInt(idx, "Necropsies")
			if val == 1 && total != nec {
				count++
				//t.Errorf("Total records %d does not equal necropsies: %d.", total, nec)
			} else if count == -1 && nec != 0 {
				count++
				//t.Errorf("%d necropsies found in non-necropsies records.", nec)
			}
		}
		if count > 0 {
			t.Errorf("Found %d species with incorrect number of records for %s.", count, name)
		}
	}
}

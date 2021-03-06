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
	rates := cancerrates.GetCancerRates(db, 1, 0, false, false, false, false, "all", "", "")
	compareTables(t, "Cancer Rates", getExpectedRates(), rates)
}

func TestPrevlenceTotals(t *testing.T) {
	// Compares total records with RecordsWithDenominators
	var count int
	flag.Parse()
	db := codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), *password)
	rates := cancerrates.GetCancerRates(db, 1, 0, false, false, false, false, "all", "", "")
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
		rates := cancerrates.GetCancerRates(db, 1, val, false, false, false, false, "all", "", "")
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

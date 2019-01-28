// This script will upload patient data to the comparative oncology database

package dbupload

import (
	"bufio"
	"fmt"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"math"
	"strconv"
	"strings"
)

func sizeOf(list [][]string) int {
	// Returns size of array in bytes
	ret := 0
	for _, i := range list {
		for _, j := range i {
			ret += len([]byte(j))
		}
	}
	return ret * 8
}

func getDenominator(size int) int {
	// Returns denominator for subsetting upload slice (size in bytes / 16Mb)
	max := 15000000.0
	return int(math.Floor(float64(size) / max))
}

func uploadPatients(db *dbIO.DBIO, table string, list [][]string) {
	// Uploads patient entries to db
	l := len(list)
	den := getDenominator(sizeOf(list))
	fmt.Println(den, sizeOf(list))
	if den <= 1 {
		// Upload slice at once
		vals, l := dbIO.FormatSlice(list)
		db.UpdateDB(table, vals, l)
	} else {
		// Upload in chunks
		var set [][][]string
		idx := l / den
		ind := 0
		for i := 0; i < den; i++ {
			if ind+idx > l {
				// Get last less than idx rows
				idx = l - ind + 1
			}
			sub := list[ind : ind+idx]
			set = append(set, sub)
			ind = ind + idx
		}
		for _, i := range set {
			vals, ln := dbIO.FormatSlice(i)
			db.UpdateDB(table, vals, ln)
		}
	}
}

func tumorPairs(typ, loc string) [][]string {
	// Returns slice of pairs of type, location
	var ret [][]string
	types := strings.Split(typ, ";")
	locations := strings.Split(loc, ";")
	for idx, i := range types {
		if idx < len(locations) {
			ret = append(ret, []string{strings.TrimSpace(i), strings.TrimSpace(locations[idx])})
		}
	}
	return ret
}

//----------------------------------------------------------------------------

type entries struct {
	count    int
	p        [][]string
	d        [][]string
	t        [][]string
	s        [][]string
	accounts map[string]map[string]string
	taxa  map[string]string
	col      map[string]int
	length   int
}

func newEntries(count int) *entries {
	// Initializes new struct
	e := new(entries)
	e.count = count
	e.accounts = make(map[string]map[string]string)
	e.taxa = make(map[string]string)
	return e
}

func (e *entries) addTumors(id string, row []string) {
	// Assign ID to all tumor, location pairs tumorPairs
	t := []string{id, row[e.col["Primary"]], row[e.col["Malignant"]]}
	pairs := tumorPairs(row[e.col["Type"]], row[e.col["Location"]])
	for _, i := range pairs {
		tumor := append(t, i...)
		e.t = append(e.t, tumor)
	}
}

func (e *entries) addDiagnosis(id string, row []string) {
	// Diagnosis entry: ID, masspresent, hyperplasia, necropsy, metastasis_id
	d := []string{id, row[e.col["MassPresent"]], row[e.col["Hyperplasia"]], row[e.col["Necropsy"]], row[e.col["Metastasis"]]}
	e.d = append(e.d, d)
}

func (e *entries) addSource(id, aid string, row []string) {
	// ID, service, account_id
	e.s = append(e.s, []string{id, row[e.col["Service"]], aid})
}

func (e *entries) addPatient(id, taxaid string, row []string) {
	// Formats patient data for upload
	if strings.Contains(row[e.col["ID"]], "NA") == true {
		// Make sure source ID is an integer
		row[e.col["ID"]] = "-1"
	} else if len(row[e.col["Age"]]) > 6 {
		// Make sure age does not exceed decimal precision
		row[e.col["Age"]] = row[e.col["Age"]][:7]
	}
	// ID, Sex, Age, Castrated, taxa_id, source_id, Species, Date, Comments
	p := []string{id, row[e.col["Sex"]], row[e.col["Age"]], row[e.col["Castrated"]], taxaid, row[e.col["ID"]], row[e.col["Date"]], row[e.col["Comments"]]}
	e.p = append(e.p, p)
}

func (e *entries) evaluateRow(row []string) int {
	// Appends data to relevent slice
	miss := 1
	if strings.ToUpper(row[4]) != "NA" {
		t := getTaxon(row[e.col["Family"]], row[e.col["Genus"]], row[e.col["Species"]])
		taxaid, exists := e.taxa[t]
		ac, ex := e.accounts[row[e.col["Account"]]]
		if len(row) == e.length && exists == true && ex == true {
			// Skip entries without valid species and source data
			aid, inmap := ac[row[e.col["Submitter"]]]
			if inmap == true {
				e.count++
				id := strconv.Itoa(e.count)
				e.addPatient(id, taxaid, row)
				e.addSource(id, aid, row)
				e.addDiagnosis(id, row)
				e.addTumors(id, row)
				miss--
			}
		}
	}
	return miss
}

func (e *entries) extractPatients(infile string) {
	// Assigns patient data to appropriate slices with unique entry IDs
	missed := 0
	first := true
	start := e.count
	fmt.Printf("\n\tExtracting patient data from %s\n", infile)
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := bufio.NewScanner(f)
	for input.Scan() {
		spl := strings.Split(string(input.Text()), ",")
		if first == false {
			missed += e.evaluateRow(spl)
		} else {
			e.col = getColumns(spl)
			e.length = len(spl)
			first = false
		}
	}
	fmt.Printf("\tExtracted %d records.\n", e.count-start)
	if missed > 0 {
		fmt.Printf("\t[Warning] Count not find taxa ID or source ID for %d records.\n", missed)
	}
}

func LoadPatients(db *dbIO.DBIO, infile string) {
	// Loads unique patient info to appropriate tables
	e := newEntries(db.GetMax("Patient", "ID"))
	e.accounts = MapOfMaps(db.GetTable("Accounts"))
	e.taxa = getTaxaIDs(db, false)
	// Get entry slices and upload to db
	e.extractPatients(infile)
	uploadPatients(db, "Patient", e.p)
	uploadPatients(db, "Diagnosis", e.d)
	uploadPatients(db, "Tumor", e.t)
	uploadPatients(db, "Source", e.s)
	// Recacluate species totals
	SpeciesTotals(db)
}

// This script will upload patient data to the comparative oncology database

package dbupload

import (
	"bufio"
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"log"
	"os"
	"strconv"
	"strings"
)

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

type entries struct {
	accounts  map[string]map[string]string
	col       map[string]int
	count     int
	d         [][]string
	dbset     bool
	ex        *Existing
	infant    *Infancy
	length    int
	logger    *log.Logger
	p         [][]string
	s         [][]string
	t         [][]string
	taxa      map[string]string
	unmatched [][]string
}

func newEntries(db *dbIO.DBIO, test bool) *entries {
	// Initializes new struct
	e := new(entries)
	e.logger = codbutils.GetLogger()
	if db != nil {
		e.count = db.GetMax("Patient", "ID")
		e.accounts = codbutils.MapOfMaps(db.GetTable("Accounts"))
		e.infant = NewInfancy(db)
		e.taxa = GetTaxaIDs(db, false)
		e.dbset = true
	} else {
		e.accounts = make(map[string]map[string]string)
		e.taxa = make(map[string]string)
	}
	if test {
		e.ex = NewExisting(nil)
	} else {
		e.ex = NewExisting(db)
	}
	return e
}

func checkInt(val string) string {
	// Makes sure value is an integer
	if _, err := strconv.Atoi(val); err != nil {
		val = "-1"
	}
	return val
}

func (e *entries) addUnmatched(row []string) {
	// Adds row elements to unmatched
	rec := []string{row[e.col["ID"]]}
	rec = append(rec, row[e.col["Name"]])
	rec = append(rec, row[e.col["Sex"]])
	rec = append(rec, row[e.col["Age"]])
	rec = append(rec, row[e.col["Date"]])
	rec = append(rec, checkInt(row[e.col["Masspresent"]]))
	rec = append(rec, checkInt(row[e.col["Necropsy"]]))
	rec = append(rec, row[e.col["Comments"]])
	rec = append(rec, row[e.col["Service"]])
	e.unmatched = append(e.unmatched, rec)
}

func (e *entries) addTumors(id string, row []string) {
	// Assign ID to all tumor, location pairs tumorPairs
	t := []string{id, row[e.col["Primary"]], row[e.col["Malignant"]]}
	pairs := tumorPairs(row[e.col["TumorType"]], row[e.col["Location"]])
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
	e.s = append(e.s, []string{id, row[e.col["Service"]], row[e.col["Zoo"]], row[e.col["AZA"]], row[e.col["Institute"]], "-1", aid})
}

func formatAge(age string) string {
	// Returns age formatted for sql upload
	ret := "-1.0"
	if _, err := strconv.ParseFloat(age, 64); err == nil {
		if strings.Contains(age, ".") {
			s := strings.Split(age, ".")
			for len(s[1]) < 2 {
				s[1] += "0"
			}
			age = fmt.Sprintf("%s.%s", s[0], s[1][:2])
			if len(age) <= 7 {
				ret = age
			}
		} else {
			ret = age + ".00"
		}
	}
	return ret
}

func (e *entries) addPatient(id, taxaid, age string, row []string) {
	// Formats patient data for upload
	infant := "-1"
	if strings.Contains(row[e.col["ID"]], "NA") == true {
		// Make sure source ID is an integer
		row[e.col["ID"]] = "-1"
	}
	if e.dbset {
		infant = e.infant.SetInfant(taxaid, age, row[e.col["Comments"]])
	}
	// ID, Sex, Age, Castrated, taxa_id, source_id, Species, Date, Comments
	p := []string{id, row[e.col["Sex"]], age, infant, row[e.col["Castrated"]], taxaid, row[e.col["ID"]], row[e.col["Name"]], row[e.col["Date"]], row[e.col["Year"]], row[e.col["Comments"]]}
	e.p = append(e.p, p)
}

func (e *entries) evaluateRow(row []string) {
	// Appends data to relevent slice
	t := getTaxon(row[e.col["Genus"]], row[e.col["Species"]])
	taxaid, exists := e.taxa[t]
	if len(row) == e.length && exists == true {
		// Skip entries without valid species/genus data
		aid := "-1"
		ac, ex := e.accounts[row[e.col["Account"]]]
		if ex == true {
			id, inmap := ac[row[e.col["Submitter"]]]
			if inmap == true {
				aid = id
			}
		}
		age := formatAge(row[e.col["Age"]])
		if !e.ex.Exists(aid, row[e.col["ID"]], age, taxaid, row[e.col["Date"]]) {
			e.count++
			id := strconv.Itoa(e.count)
			e.addPatient(id, taxaid, age, row)
			e.addSource(id, aid, row)
			e.addDiagnosis(id, row)
			e.addTumors(id, row)
		}
	} else if !e.ex.Exists("", row[e.col["ID"]], row[e.col["Age"]], "", row[e.col["Date"]]) {
		e.addUnmatched(row)
	}
}

func (e *entries) extractPatients(infile string) {
	// Assigns patient data to appropriate slices with unique entry IDs
	first := true
	e.logger.Printf("Extracting patient data from %s\n", infile)
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := iotools.GetScanner(f)
	for input.Scan() {
		spl := strings.Split(string(input.Text()), ",")
		if first == false {
			e.evaluateRow(spl)
		} else {
			e.col = iotools.GetHeader(spl)
			e.length = len(spl)
			first = false
		}
	}
	e.logger.Printf("Extracted %d records.\n", len(e.p))
	e.logger.Printf("Found %d unmatched records.", len(e.unmatched))
}

func LoadPatients(db *dbIO.DBIO, infile string, test, proceed bool) {
	// Loads unique patient info to appropriate tables
	e := newEntries(db, test)
	// Get entry slices and upload to db
	e.extractPatients(infile)
	if !proceed {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\tProceed with upload?")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(strings.ToLower(text))
		if text == "y" || text == "yes" {
			proceed = true
		}
	}
	if proceed {
		e.logger.Println("Proceeding with upload...")
		if len(e.p) > 0 {
			db.UploadSlice("Patient", e.p)
			db.UploadSlice("Diagnosis", e.d)
			db.UploadSlice("Tumor", e.t)
			db.UploadSlice("Source", e.s)
		}
		if len(e.unmatched) > 0 {
			db.UploadSlice("Unmatched", e.unmatched)
		}
	} else {
		e.logger.Println("Aborting upload.")
	}
}

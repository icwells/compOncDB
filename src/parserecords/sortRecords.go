// This script defines functions for sorting entries data

package parserecords

import (
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/go-tools/iotools"
	"os"
	"strings"
	"sync"
)

func subsetLine(idx int, line []string) string {
	// Returns line[idx]/NA
	ret := "NA"
	if idx > -1 && idx < len(line) {
		ret = strings.TrimSpace(line[idx])
		if len(ret) <= 0 {
			ret = "NA"
		}
	}
	return ret
}

func (e *entries) getSubmitter(line []string) []string {
	// Gets corrected submitter name and type from map
	var ret []string
	key := subsetLine(e.col.submitter, line)
	if key != "NA" {
		val, ex := e.accounts[key]
		if ex == true {
			ret = val
		}
	}
	return ret
}

func (e *entries) sortLine(wg *sync.WaitGroup, mut *sync.RWMutex, debug bool, out *os.File, line []string) {
	// Write or stores formatted string
	defer wg.Done()
	write := false
	rec := newRecord()
	var idx int
	if e.col.common >= 0 && line[e.col.common] != "NA" {
		// Get common name if present
		idx = e.col.common
	} else if e.col.species >= 0 {
		idx = e.col.species
	} else {
		e.logger.Fatalf("Cannot determine species column. Exiting.\n")
	}
	if len(line) >= e.col.max && len(line[idx]) >= 3 && strings.ToUpper(line[idx]) != "N/A" {
		// Proceed if line is properly formatted and species is present
		id := subsetLine(e.col.id, line)
		rec.setID(id)
		// Set diagnosis info
		e.parseLine(&rec, line)
		// Replace entry with scientific name
		sp, ex := e.taxa[line[idx]]
		if ex == true {
			rec.setSpecies(sp)
		}
		rec.name = subsetLine(idx, line)
		rec.setDate(subsetLine(e.col.date, line))
		rec.setYear(subsetLine(e.col.year, line))
		rec.setComments(subsetLine(e.col.comments, line))
		rec.service = e.service
		rec.setAccount(subsetLine(e.col.account, line))
		rec.setSubmitter(e.getSubmitter(line))
		if e.col.code >= 0 && len(line[e.col.code]) > 0 {
			// Store code for debugging
			rec.code = line[e.col.code]
		}
		if e.dupsPresent == true {
			rec.setPatient(line, e.col)
			if e.inDuplicates(rec) == true {
				// Resolve duplicate records and write when done
				e.resolveDuplicates(rec)
			} else {
				write = true
			}
		} else {
			write = true
		}
	}
	if write == true {
		mut.Lock()
		out.WriteString(rec.String(debug) + "\n")
		mut.Unlock()
		e.extracted++
	}
}

func (e *entries) SortRecords(debug bool, infile, outfile string) {
	// Sorts data and merges if necessary
	first := true
	var wg sync.WaitGroup
	var mut sync.RWMutex
	var total int
	e.logger.Println("Parsing input records...")
	f := iotools.OpenFile(infile)
	defer f.Close()
	out := e.getOutputFile(outfile, codbutils.ParseHeader(debug))
	defer out.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		line := string(scanner.Text())
		if first == false {
			total++
			if e.d != "," && strings.Contains(line, ",") {
				line = strings.Replace(line, ",", ";", -1)
			}
			s := strings.Split(strings.Replace(line, "\"", "", -1), e.d)
			wg.Add(1)
			go e.sortLine(&wg, &mut, debug, out, s)
		} else {
			e.parseHeader(line)
			first = false
		}
	}
	wg.Wait()
	if e.dupsPresent == true {
		for _, val := range e.dups.records {
			// Write each stored record before closing
			for _, v := range val {
				out.WriteString(v.String(debug) + "\n")
				e.extracted++
			}
		}
	}
	e.logger.Printf("Extracted %d records from %d total records.\n", e.extracted, total)
	e.logger.Printf("Found diagnosis data for %d records.\n", e.found)
	e.logger.Printf("Found complete information for %d records.\n", e.complete)
}

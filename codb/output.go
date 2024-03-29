// Contains functions for interacting with slq database

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/cancerrates"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/compOncDB/src/dbextract"
	"github.com/icwells/compOncDB/src/search"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/dataframe"
	"github.com/icwells/go-tools/iotools"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
)

type TableRow struct {
	Cells []string
}

type HTMLTable struct {
	Header []string
	Body   []TableRow
}

func newFlash(w http.ResponseWriter, msg string) *Output {
	// Returns output with flash error message
	var o Output
	o.w = w
	o.Flash = msg
	return &o
}

type Output struct {
	User      string
	Update    string
	Flash     string
	File      string
	Outfile   string
	Pathfile  string
	Pathology string
	Plot      string
	Table     HTMLTable
	Count     string
	db        *dbIO.DBIO
	pw        string
	w         http.ResponseWriter
	r         *http.Request
}

func newOutput(w http.ResponseWriter, r *http.Request, user, pw, ut string) (*Output, error) {
	// Returns empty output struct
	var err error
	o := new(Output)
	o.w = w
	o.r = r
	o.User = user
	o.Update = ut
	if pw != "" {
		o.pw = pw
		o.db, err = dbIO.Connect(C.config.Host, C.config.Database, o.User, o.pw)
		if err == nil {
			o.db.GetTableColumns()
		}
	}
	return o, err
}

func (o *Output) formatTable(header []string, table [][]string) {
	// Formats slice into table for display in a browser
	o.Table.Header = header
	for _, i := range table {
		var c TableRow
		c.Cells = i
		o.Table.Body = append(o.Table.Body, c)
	}
}

func (o *Output) getTempFile(name string) {
	// Stores path to named file in tmp directory
	o.File = fmt.Sprintf("%s.%s.csv", name, codbutils.GetTimeStamp())
	o.Outfile = fmt.Sprintf("/tmp/%s", o.File)
}

func (o *Output) extractTable() {
	// Extracts given table from the database
	name := strings.TrimSpace(o.r.PostForm.Get("Table"))
	if name != "" && name != "Empty" {
		table := o.db.GetTable(name)
		o.getTempFile(name)
		codbutils.WriteResults(o.Outfile, o.db.Columns[name], table)
		C.renderTemplate(C.temp.result, o)
	} else {
		o.Flash = "Please select a table to extract."
		C.renderTemplate(C.temp.menu, o)
	}
}

func (o *Output) lifeHistorySummary() {
	// Returns summary of life history table
	opt := setOptions(o.r)
	o.getTempFile("lifeHistorySummary")
	res := dbextract.LifeHistorySummary(o.db, opt.AllTaxa)
	res.ToCSV(o.Outfile)
	o.Count = fmt.Sprintf("\tFound %d records matching search criteria.\n", res.Length())
	C.renderTemplate(C.temp.result, o)
}

func (o *Output) neoplasiaPrevalence() {
	// Performs cancer rate calculations
	var eval string
	var necropsy int
	var res, pathology *dataframe.Dataframe
	opt := setOptions(o.r)
	if opt.Taxa != "" && opt.Operation != "" && opt.Value != "" {
		eval = opt.Taxa + opt.Operation + opt.Value
	}
	switch opt.Necropsy {
	case "necropsyonly":
		necropsy = 1
	case "nonnecropsy":
		necropsy = -1
	}
	c := cancerrates.NewCancerRates(o.db, opt.Min, opt.Keepall, opt.Tissue, opt.Location)
	c.SearchSettings(necropsy, opt.Infant, opt.Wild, opt.Source)
	c.OutputSettings(opt.Agecol, opt.Lifehistory, opt.Sexcol, opt.Taxacol)
	res, pathology = c.GetCancerRates(eval)
	if opt.Pathology {
		o.Pathfile = fmt.Sprintf("pathologyRecords.min%d.%s.csv", opt.Min, codbutils.GetTimeStamp())
		o.Pathology = fmt.Sprintf("/tmp/%s", o.Pathfile)
		pathology.ToCSV(o.Pathology)
	}
	if opt.Location == "" {
		// Use location as file name stem
		opt.Location = "neoplasiaPrevalence"
	}
	o.renderResults(opt, res, fmt.Sprintf("%s.min%d", opt.Location, opt.Min))
}

func (o *Output) referenceTaxonomy() {
	// Returns merged common name and taxonomy tables
	table := dbextract.GetReferenceTaxonomy(o.db)
	o.getTempFile("mergedTaxonomy")
	table.ToCSV(o.Outfile)
	C.renderTemplate(C.temp.result, o)
}

func (o *Output) renderResults(opt *Options, res *dataframe.Dataframe, name string) {
	// Renders search and cancer rate results
	if res.Length() >= 1 {
		// Format link for download whether or not results are printed to screen
		o.getTempFile(name)
		res.ToCSV(o.Outfile)
		if opt.Print {
			o.formatTable(strings.Split(res.FormatHeader(","), ","), res.ToSlice())
		}
	}
	if o.Count == "" {
		o.Count = fmt.Sprintf("\tFound %d records matching search criteria.\n", res.Length())
	}
	C.renderTemplate(C.temp.result, o)
}

func (o *Output) searchDB() {
	// Performs searches
	var res *dataframe.Dataframe
	var eval [][]codbutils.Evaluation
	opt := setOptions(o.r)
	eval, o.Flash = checkEvaluations(o.r, o.db.Columns)
	if o.Flash == "" {
		res, o.Count = search.SearchColumns(o.db, codbutils.GetLogger(), eval, opt.Infant, false)
	}
	if o.Flash == "" {
		o.renderResults(opt, res, o.User)
	} else {
		// Return to search page with flash message
		C.renderTemplate(C.temp.menu, o)
	}
}

func (o *Output) getTable(outdir, name string) string {
	// Saves table and returns name
	table := o.db.GetTable(name)
	ret := filepath.Join(outdir, name+".csv")
	codbutils.WriteResults(ret, o.db.Columns[name], table)
	return ret
}

func (o *Output) barplot() {
	// Calls barplot.py and returns result file
	tmp := "/tmp/"
	// Make outdir
	outdir, _ := iotools.FormatPath(fmt.Sprintf("%sbarplots%s", tmp, codbutils.GetTimeStamp()), true)
	d := o.getTable(outdir, "Diagnosis")
	p := o.getTable(outdir, "Patient")
	s := o.getTable(outdir, "Source")
	cmd := exec.Command("python", "barplot.py", d, p, s, outdir)
	cmd.Run()
	if f, err := filepath.Glob(outdir + "*.svg"); err == nil {
		if len(f) > 0 {
			o.Plot = f[0]
			fmt.Println(o.Plot)
			fs := http.FileServer(http.Dir(outdir))
			http.Handle(outdir, http.StripPrefix(outdir, fs))
		}
	}
}

func (o *Output) summary() {
	// Returns general database summary
	header := []string{"Field", "Total", "%"}
	o.formatTable(header, dbextract.GetSummary(o.db))
	o.barplot()
	C.renderTemplate(C.temp.result, o)
}

func (o *Output) tissueLeaderBoard() {
	// Returns database tissue leaderboard
	var res *dataframe.Dataframe
	opt := setOptions(o.r)
	if opt.TumorType != "" {
		opt.Print = true
		res = search.LeaderBoard(o.db, opt.TumorType)
	} else {
		res = search.TypesPerSpecies(o.db, opt.Min)
	}
	o.renderResults(opt, res, o.User)
}

func (o *Output) routePost(source string) {
	// Sends post data to appropriate function
	o.r.ParseForm()
	switch source {
	case C.u.summary:
		o.summary()
	case C.u.tissue:
		o.tissueLeaderBoard()
	case C.u.reftaxa:
		o.referenceTaxonomy()
	case C.u.table:
		o.extractTable()
	case C.u.lifehist:
		o.lifeHistorySummary()
	case C.u.prevalence:
		o.neoplasiaPrevalence()
	case C.u.output:
		o.searchDB()
	}
}

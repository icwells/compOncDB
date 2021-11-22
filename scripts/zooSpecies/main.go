// Extracts zoo names for target species

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/simpleset"
	"gopkg.in/alecthomas/kingpin.v2"
	"sort"
	"strings"
	"time"
)

var (
	approved = kingpin.Flag("approved", "Extract only approved sources.").Default("false").Bool()
	infile   = kingpin.Flag("infile", "Name of input csv.").Short('i').Required().String()
	outfile  = kingpin.Flag("outfile", "Name of output file.").Short('o').Required().String()
	user     = kingpin.Flag("user", "MySQL username.").Short('u').Required().String()
)

type zoos struct {
	accounts map[string]string
	all      *simpleset.Set
	db       *dbIO.DBIO
	ids      map[string][]string
	sources  map[string]string
	taxa     map[string]*simpleset.Set
}

func newZoos() *zoos {
	// Returns new struct
	z := new(zoos)
	z.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	z.accounts = codbutils.EntryMap(z.db.GetColumns("Accounts", []string{"submitter_name", "account_id"}))
	z.all = simpleset.NewStringSet()
	z.sources = codbutils.EntryMap(z.db.GetColumns("Source", []string{"account_id", "ID"}))
	z.taxa = make(map[string]*simpleset.Set)
	return z
}

func (z *zoos) setAccounts() {
	// Stores patient ids by taxa
	fmt.Println("\n\tStoring patient IDs by species...")
	var species []string
	records, header := iotools.ReadFile(*infile, true)
	for _, i := range records {
		species = append(species, i[header["Species"]])
	}
	for _, i := range z.db.GetRows("Records", "Species", strings.Join(species, ","), "Species,Zoo,Institute,Approved,ID") {
		if !*approved || i[3] == "1" {
			if i[1] == "1" || i[2] == "1" {
				// Select approved zoo/institute records
				sp := i[0]
				if aid, exists := z.sources[i[4]]; exists {
					if v, ex := z.accounts[aid]; ex {
						if _, ex := z.taxa[sp]; !ex {
							z.taxa[sp] = simpleset.NewStringSet()
						}
						z.taxa[sp].Add(v)
						z.all.Add(v)
					}
				}
			}
		}
	}
}

func (z *zoos) sortAccounts(s *simpleset.Set) string {
	// Converts set to string slice, sorts slice, and returns joined string
	ret := s.ToStringSlice()
	sort.Strings(ret)
	return strings.Join(ret, ";")
}

func (z *zoos) write() {
	// Writes names to file
	out := iotools.CreateFile(*outfile)
	out.WriteString("Genus,Species,Zoos\n")
	all := []string{"", "all", z.sortAccounts(z.all)}
	out.WriteString(strings.Join(all, ",") + "\n")
	for k, v := range z.taxa {
		genus := strings.Split(k, " ")[0]
		row := []string{genus, k, z.sortAccounts(v)}
		out.WriteString(strings.Join(row, ",") + "\n")
	}
}

func main() {
	start := time.Now()
	kingpin.Parse()
	z := newZoos()
	z.setAccounts()
	z.write()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

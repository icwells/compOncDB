// Extracts zoo names for target species

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/dbIO"
	"github.com/icwells/simpleset"
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
	"time"
)

var (
	infile   = kingpin.Flag("infile", "Name of input csv.").Short('i').Required().String()
	outfile  = kingpin.Flag("outfile", "Name of output file.").Short('o').Required().String()
	user     = kingpin.Flag("user", "MySQL username (default is root).").Short('u').Required().String()
)

type zoos struct {
	accounts map[string]string
	all		*simpleset.Set
	db		*dbIO.DBIO
	header  map[string]int
	ids		map[string]string
	records [][]string
	sources map[string][]string
	species map[string]*simpleset.Set
	taxa	map[string]*simpleset.Set
}

func newZoos() *zoos {
	// Returns new struct
	z := new(zoos)
	z.db = codbutils.ConnectToDatabase(codbutils.SetConfiguration(*user, false), "")
	z.accounts = codbutils.EntryMap(z.db.GetColumns("Accounts", []string{"submitter_name", "account_id"}))
	z.all = simpleset.NewStringSet()
	z.records, z.header = iotools.ReadFile(*infile, true)
	z.ids = codbutils.EntryMap(z.db.GetColumns("Patient", []string{"ID", "taxa_id"}))
	z.sources = codbutils.ToMap(z.db.GetColumns("Source", []string{"ID", "Zoo", "Institute", "account_id"}))
	z.species = make(map[string]*simpleset.Set)
	z.taxa = make(map[string]*simpleset.Set)
	return z
}

func (z *zoos) setTaxa() {
	// Stores patient ids by taxa
	fmt.Println("\n\tStoring patient IDs by species...")
	taxa := codbutils.EntryMap(z.db.GetColumns("Taxonomy", []string{"taxa_id", "Species"}))
	for _, i := range z.records {
		sp := i[z.header["Species"]]
		if _, ex := z.taxa[sp]; !ex {
			z.taxa[sp] = simpleset.NewStringSet()
		}
		if tid, ex := taxa[sp]; ex {
			if pid, ex := z.ids[tid]; ex {
				z.taxa[sp].Add(pid)
			}
		}
	}
}

func (z *zoos) setNames() {
	// Stores zoo names for each species
	fmt.Println("\tStoring account names by species...")
	for k, v := range z.taxa {
		z.species[k] = simpleset.NewStringSet()
		for _, i := range v.ToStringSlice() {
			if s, ex := z.sources[i]; ex {
				if s[0] == "1" || s[1] == "1" {
					// Ignore private records
					if name, ex := z.accounts[s[2]]; ex {
						z.species[k].Add(name)
					}
				}
			}
		}
	}
}

func (z *zoos) write() {
	// Writes names to file
	out := iotools.CreateFile(*outfile)
	out.WriteString("Genus,Species,Necropsies,Zoos\n")
	for _, i := range z.records {
		sp := i[z.header["Species"]]
		if v, ex := z.species[sp]; ex {
			row := append(i, strings.Join(v.ToStringSlice(), ";"))
			out.WriteString(strings.Join(row, ",") + "\n")
		}
	}
}

func main() {
	start := time.Now()
	kingpin.Parse()
	z := newZoos()
	z.setTaxa()
	z.setNames()
	z.write()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

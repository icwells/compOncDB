// Contains functions for convertng slice of string slices to map

package codbutils

import (
	"bufio"
	"fmt"
	"github.com/icwells/dbIO"
	"github.com/icwells/go-tools/iotools"
	"github.com/icwells/simpleset"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

func GetLogger() *log.Logger {
	// Returns logger
	return log.New(os.Stdout, "compOnDB: ", log.Ldate|log.Ltime)
}

func Getutils() string {
	// Returns path to utils directory
	return path.Join(iotools.GetGOPATH(), "src/github.com/icwells/compOncDB/utils")
}

func getAbsPath(f string) string {
	// Prepends GOPATH to file name if needed
	if !strings.Contains(f, string(os.PathSeparator)) {
		f = path.Join(Getutils(), f)
	}
	if iotools.Exists(f) == false {
		GetLogger().Fatalf("Cannot find %s file. Exiting.\n", f)
		os.Exit(1)
	}
	return f
}

type Configuration struct {
	Host     string
	Database string
	User     string
	Testdb   string
	Tables   string
	Test     bool
}

func SetConfiguration(user string, test bool) Configuration {
	// Gets setting from config.txt
	var c Configuration
	c.Test = test
	c.User = user
	f := iotools.OpenFile(getAbsPath("config.txt"))
	defer f.Close()
	scanner := iotools.GetScanner(f)
	for scanner.Scan() {
		s := strings.Split(string(scanner.Text()), "=")
		for idx, i := range s {
			s[idx] = strings.TrimSpace(i)
		}
		switch s[0] {
		case "host":
			c.Host = s[1]
		case "database":
			c.Database = s[1]
		case "test_database":
			c.Testdb = s[1]
		case "table_columns":
			c.Tables = getAbsPath(s[1])
		}
	}
	return c
}

func ConnectToDatabase(c Configuration) *dbIO.DBIO {
	// Manages call to Connect and GetTableColumns
	d := c.Database
	if c.Test == true {
		d = c.Testdb
	}
	db, err := dbIO.Connect(c.Host, d, c.User, "")
	if err != nil {
		GetLogger().Fatal(err)
	}
	db.GetTableColumns()
	return db
}

func ReadList(infile string, idx int) []string {
	// Reads list of queries from file
	set := simpleset.NewStringSet()
	var d string
	first := true
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := iotools.GetScanner(f)
	for input.Scan() {
		line := string(input.Text())
		if first == false {
			s := strings.Split(line, d)
			if len(s) > idx {
				// Replace underscores if present
				name := strings.Replace(s[idx], "_", " ", -1)
				name = strings.TrimSpace(name)
				if len(name) > 1 {
					set.Add(name)
				}
			}
		} else {
			d, _ = iotools.GetDelim(line)
			first = false
		}
	}
	return set.ToStringSlice()
}

func printArray(header string, table [][]string) {
	// Prints slice of string slcies to screen
	head := strings.Split(header, ",")
	// Wrap in newlines
	fmt.Println()
	fmt.Println(strings.Join(head, "\t"))
	for _, row := range table {
		fmt.Println(strings.Join(row, "\t"))
	}
	fmt.Println()
}

func WriteResults(outfile, header string, table [][]string) {
	// Wraps calls to writeCSV/printArray
	if len(table) > 0 {
		if outfile != "nil" {
			iotools.WriteToCSV(outfile, header, table)
		} else {
			printArray(header, table)
		}
	}
}

func DeleteEntries(d *dbIO.DBIO, table, column, value string) {
	// Deletes matches from appropriate tables
	var input string
	reader := bufio.NewReader(os.Stdin)
	if d.Database == "testDataBase" {
		input = "Y"
	} else {
		fmt.Printf("\tAre you sure you want to delete all records from %s where %s equals %s (Enter Y for 'yes')? ", table, column, value)
		input, _ = reader.ReadString('\n')
	}
	if strings.TrimSpace(strings.ToUpper(input)) == "Y" {
		fmt.Println("\tProceeding with deletion...")
		d.DeleteRow(table, column, value)
	} else {
		fmt.Println("\tSkipping deletion.")
	}
}

func GetUpdateTime(d *dbIO.DBIO) string {
	// Returns most recent update time
	ret := d.GetColumnText("Update_time", "Time")
	return ret[len(ret) - 1]
}

func UpdateTimeStamp(d *dbIO.DBIO) {
	// Stores current time stamp in update_time table
	t := time.Now().Format(time.RFC822)
	cmd := fmt.Sprintf("INSERT INTO Update_time(Time) VALUES('%s');", t)
	d.Insert("Update_time", cmd)
}

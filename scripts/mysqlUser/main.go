// Adds new user to MySQL database

package main

import (
	"fmt"
	"github.com/icwells/compOncDB/src/codbutils"
	"github.com/icwells/dbIO"
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
	"time"
)

var (
	app  = kingpin.New("mysqlUser", "Adds new user to MySQL database. Must use root password.")
	all  = kingpin.Flag("all", "Grant all privileges to user. User will only have remote access if all is false.").Default("false").Bool()
	user = kingpin.Flag("user", "Name of new user. Will also be used as a temporary password.").Required().Short('u').String()
)

type mysqluser struct {
	all        bool
	database   string
	db         *dbIO.DBIO
	host       string
	name       string
	permission string
	records    string
	root       string
	//tables     []string
}

func newUser() *mysqluser {
	// Returns new struct
	m := new(mysqluser)
	m.all = *all
	m.database = "comparativeOncology"
	m.host = "localhost"
	m.name = *user
	m.permission = "SELECT"
	if m.all {
		m.permission = "ALL PRIVILEGES"
	}
	m.records = "Records"
	m.root = "root"
	//m.tables = []string{"Patient", "Unmatched", "Taxonomy", "Common", "Totals", "Denominators", "Source", "Diagnosis", "Tumor", "Life_history", "Update_time"}
	m.connect()
	return m
}

func (m *mysqluser) connect() {
	// Connects to database
	c := codbutils.SetConfiguration(m.root, false)
	//c.Host = m.host
	m.db = codbutils.ConnectToDatabase(c, "")
}

func (m *mysqluser) execute(command string) {
	// Executes given command
	cmd, err := m.db.DB.Prepare(command)
	if err != nil {
		panic(fmt.Sprintf("[Error] Formatting command %s: %v\n", command, err))
	} else {
		if _, err = cmd.Exec(); err != nil {
			panic(fmt.Sprintf("[Error] Executing %s: %v\n", command, err))
		}
	}
}

func (m *mysqluser) setPrivileges() {
	// Grants access to user
	// Insert % seperately to avoid formatting error
	cmd := fmt.Sprintf("GRANT %s ON %s.* TO '%s'@'%s';", m.permission, m.database, m.name, "%")
	if m.all {
		m.execute(cmd)
		m.execute(strings.Replace(cmd, "%", m.host, 1))
	} else {
		m.execute(strings.Replace(cmd, "*", m.records, 1))
	}
}

func (m *mysqluser) createUser() {
	// Executes create user command
	cmd := fmt.Sprintf("CREATE USER '%s'@'%s' IDENTIFIED BY '%s';", m.name, "%", m.name)
	m.execute(cmd)
	if m.all {
		// Grant remote access
		m.execute(strings.Replace(cmd, "%", m.host, 1))
	}
}

func main() {
	start := time.Now()
	kingpin.Parse()
	m := newUser()
	fmt.Println("\n\tAdding new user...")
	m.createUser()
	m.setPrivileges()
	m.execute("FLUSH PRIVILEGES;")
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}

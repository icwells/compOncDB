// Adds new user to MySQL database

package codbutils

import (
	"fmt"
	"github.com/icwells/dbIO"
	"log"
	"strings"
)

type mysqluser struct {
	all        bool
	db         *dbIO.DBIO
	host       string
	logger     *log.Logger
	name       string
	permission string
	records    string
	updatetime string
}

func (m *mysqluser) execute(command string) {
	// Executes given command
	cmd, err := m.db.DB.Prepare(command)
	if err != nil {
		m.logger.Fatalf("[Error] Formatting command %s: %v\n", command, err)
	} else {
		if _, err = cmd.Exec(); err != nil {
			m.logger.Fatalf("[Error] Executing %s: %v\n", command, err)
		}
	}
}

func (m *mysqluser) setPrivileges() {
	// Grants access to user
	// Insert % seperately to avoid formatting error
	cmd := fmt.Sprintf("GRANT %s ON %s.* TO '%s'@'%s';", m.permission, m.db.Database, m.name, "%")
	if m.all {
		m.execute(cmd)
		m.execute(strings.Replace(cmd, "%", m.host, 1))
	} else {
		m.execute(strings.Replace(cmd, "*", m.records, 1))
		m.execute(strings.Replace(cmd, "*", m.updatetime, 1))
	}
}

func (m *mysqluser) createUser() {
	// Executes create user command
	cmd := fmt.Sprintf("CREATE USER '%s'@'%s' IDENTIFIED BY '%s';", m.name, "%", m.name)
	m.execute(cmd)
	if m.all {
		// Grant local access
		m.execute(strings.Replace(cmd, "%", m.host, 1))
	}
}

func NewUser(db *dbIO.DBIO, username string, admin bool) {
	// Returns new struct
	m := new(mysqluser)
	m.all = admin
	m.db = db
	m.host = "localhost"
	m.logger = GetLogger()
	m.name = username
	m.permission = "SELECT"
	if m.all {
		m.permission = "ALL PRIVILEGES"
	}
	m.records = "Records"
	m.updatetime = "Update_time"
	m.logger.Println("Adding new user...")
	m.createUser()
	m.setPrivileges()
	m.execute("FLUSH PRIVILEGES;")
}

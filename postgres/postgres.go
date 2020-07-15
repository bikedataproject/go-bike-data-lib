package postgres

import (
	"database/sql"
	"fmt"

	// Import the Posgres driver for the database/sql package
	_ "github.com/lib/pq"
)

// Database : The configuration for Postgres
type Database struct {
	PostgresHost       string `required:"true"`
	PostgresUser       string `required:"true"`
	PostgresPassword   string `required:"true"`
	PostgresPort       int    `default:"5432"`
	PostgresDb         string `required:"true"`
	PostgresRequireSSL string `default:"require"`

	Connection *sql.DB
}

// Connect : Connect to the remote database
func (db Database) Connect() (err error) {
	conStr, err := db.getDBConnectionString()
	if err != nil {
		return
	}
	dbTmp, err := sql.Open("postgres", conStr)
	if err != nil {
		return
	}
	db.Connection = dbTmp
	return
}

// CheckConnection : Check if the connection to the database is succesfull
func (db Database) CheckConnection() (result string, err error) {
	if err = db.Connection.Ping(); err == nil {
		result = "Pong"
	}
	return
}

// CloseConnection : Disconnect the connector from the database
func (db Database) CloseConnection() {
	defer db.Connection.Close()
}

// getDBConnectionString : Generate a connectionstring for the database
func (db Database) getDBConnectionString() (result string, err error) {
	result = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%v", db.PostgresHost, db.PostgresPort, db.PostgresUser, db.PostgresPassword, db.PostgresDb, db.PostgresRequireSSL)
	return
}

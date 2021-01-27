package db_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	db "github.com/chutified/simple-bank/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:simplebankpassword@localhost:5433/simple_bank?sslmode=disable"
)

var testQueries *db.Queries
var testDB *sql.DB

// TestMain connects to the database and initializes testQueries.
func TestMain(m *testing.M) {
	var err error

	// get db connection
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	testQueries = db.New(testDB)

	os.Exit(m.Run())
}

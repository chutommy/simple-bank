package db_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	db "github.com/chutified/simple-bank/db/sqlc"
	_ "github.com/lib/pq"
)

var (
	dbDriver = "postgres"
	dbSource = "postgres:simplebankpassword@localhost:5433/simple_bank?sslmode=disable"
)

var testQueries *db.Queries

// TestMain connects to the database and initializes testQueries.
func TestMain(m *testing.M) {
	// get db connection
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	testQueries = db.New(conn)

	os.Exit(m.Run())
}

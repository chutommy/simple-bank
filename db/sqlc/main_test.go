package db_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/chutommy/simple-bank/config"
	db "github.com/chutommy/simple-bank/db/sqlc"
	_ "github.com/lib/pq"
)

var (
	testQueries *db.Queries
	testDB      *sql.DB
)

// TestMain connects to the database and initializes testQueries.
func TestMain(m *testing.M) {
	cfg, _, err := config.LoadConfig("../..")
	if err != nil {
		log.Fatalf("cannot load configuration file: %v", err)
	}

	// get db connection
	testDB, err = sql.Open(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	testQueries = db.New(testDB)

	os.Exit(m.Run())
}

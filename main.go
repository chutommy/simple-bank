package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chutified/simple-bank/api"
	db "github.com/chutified/simple-bank/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:simplebankpassword@localhost:5432/simple_bank?sslmode=disable"

	serverAddress = "0.0.0.0:8080"
)

func main() {
	dbConn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal(fmt.Errorf("cannot connect to db: %w", err))
	}

	store := db.NewStore(dbConn)

	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to start the server: %w", err))
	}
}

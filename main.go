package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chutified/simple-bank/api"
	"github.com/chutified/simple-bank/config"
	db "github.com/chutified/simple-bank/db/sqlc"
	_ "github.com/lib/pq"
)

func main() {
	// load config
	cfg, upd, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal(fmt.Errorf("can not load config file: %w", err))
	}

	for {
		// connect to db
		dbConn, err := sql.Open(cfg.DBDriver, cfg.DBSource)
		if err != nil {
			log.Fatal(fmt.Errorf("cannot connect to db: %w", err))
		}
		// check db connection
		for i := 0; i <= 5; i++ { // 5 attemps
			if err := dbConn.Ping(); err != nil {
				if i == 5 {
					log.Fatal(err) // failed at last attemp
				}
				time.Sleep(2 * time.Second)
			} else {
				break // successfully connection
			}
		}

		store := db.NewStore(dbConn)

		server := api.NewServer(store)

		// run the server concurrently
		go func() {
			fmt.Println("Starting the server...")

			if err := server.Start(cfg.ServerAddress); err != nil {
				log.Fatal(fmt.Errorf("failed to start the server: %w", err))
			}
		}()

		// init quit channel
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-quit:
			break // exit the loop
		case <-upd:
			fmt.Println("Config file updated...")
			fmt.Println("Restating the server...")
		}

		if err := server.Stop(); err != nil {
			log.Fatal(fmt.Errorf("unsuccessfully closed server: %w", err))
		}

		// close database connection
		_ = dbConn.Close()
	}
}

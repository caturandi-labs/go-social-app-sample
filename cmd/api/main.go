package main

import (
	"database/sql"
	"fmt"
	db "github.com/caturandi-labs/go-social/internal/db"
	"github.com/caturandi-labs/go-social/internal/env"
	"github.com/caturandi-labs/go-social/internal/store"
	"log"
)

const version = "0.0.1"

func main() {
	cfg := config{
		addr: env.GetString("API_ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/go_social?sslmode=disable"),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 5),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("API_ENV", "development"),
	}

	dbConn, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		log.Panic(err)
	}

	defer func(dbConn *sql.DB) {
		err := dbConn.Close()
		if err != nil {
			log.Panic(err)
		}
	}(dbConn)

	fmt.Println("Database connection established")

	pgStore := store.NewPostgresStorage(dbConn)

	app := &application{
		config: cfg,
		store:  pgStore,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}

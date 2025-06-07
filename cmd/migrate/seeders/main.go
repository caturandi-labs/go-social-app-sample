package main

import (
	"database/sql"
	"github.com/caturandi-labs/go-social/internal/db"
	"github.com/caturandi-labs/go-social/internal/env"
	storePsql "github.com/caturandi-labs/go-social/internal/store"
	"log"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/go_social?sslmode=disable")
	log.Println(addr)
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	store := storePsql.NewPostgresStorage(conn)
	db.Seed(store)
}

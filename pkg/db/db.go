package db

import (
	"database/sql"
	"log"

	"github.com/kevingdc/pulley/pkg/config"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect(config *config.Config) {
	db, err := sql.Open("postgres", config.DBUrl)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	DB = db

	log.Println("Successfully connected to the database.")
}

func Close() {
	DB.Close()
}

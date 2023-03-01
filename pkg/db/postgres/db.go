package postgres

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect(dbURL string) {
	db, err := sql.Open("postgres", dbURL)
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

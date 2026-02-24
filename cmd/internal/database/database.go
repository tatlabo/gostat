package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const (
	connStr = "user=golang password=golang dbname=webapp host=localhost port=5432 sslmode=disable"
)

func DBConnection() {

	db := New()
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

}

var DB *sql.DB

func New() *sql.DB {
	if DB != nil {
		return DB
	}

	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return DB
}

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var (
	user     = "golang"
	password = "golang"
	dbname   = "webapp"
	host     = "localhost"
	port     = "5432"
)

var postgresUrl = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

func main() {
	// Define flag
	direction := flag.String("direction", "up", "Migration direction: up or down")
	flag.Parse()

	db, err := sql.Open("postgres", postgresUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal("Connection failed:", err)
	}

	m, err := migrate.New(
		"file://./migrations",
		postgresUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Run based on direction flag
	switch *direction {
	case "up":
		if err := m.Up(); err != nil && err.Error() != "no change" {
			log.Fatal(err)
		}
		log.Println("Migrations up completed")
	case "down":
		if err := m.Down(); err != nil && err.Error() != "no change" {
			log.Fatal(err)
		}
		log.Println("Migrations down completed")
	default:
		log.Fatal("Invalid direction. Use 'up' or 'down'")
	}
}

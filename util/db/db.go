package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func InitDB() *sql.DB {
	var err error
	// PostgreSQL connection string format: "postgres://username:password@host:port/dbname?sslmode=disable"
	connStr := "postgres://wishtournament:root@localhost:5433/wishtournament?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Optionally, ping the database to check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot connect to the database:", err)
	}

	return db
}

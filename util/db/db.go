package db

import (
	"database/sql"
	"log"
)

func InitDB() *sql.DB {
	var err error
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/wishticket")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

package dev

import (
	"database/sql"
	"log"
	"net/http"
	"wishtournament/util/hashing"
)

func RegisterTicketRoute(router *http.ServeMux, db *sql.DB) {
	baseRoute := "/dev"
	router.HandleFunc(baseRoute, func(w http.ResponseWriter, r *http.Request) {
		handleEncryptPassword(w, r, db)
	})

}

func handleEncryptPassword(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println(hashing.HashPassword("admin"))
}
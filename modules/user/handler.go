package user

import (
	"database/sql"
	"log"
	"net/http"
)

func RegisterUserRoute(router *http.ServeMux, db *sql.DB) {

	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		handleUsers(w, r, db)
	})
	router.HandleFunc("/users/login", func(w http.ResponseWriter, r *http.Request) {
		handleLogin(w, r, db)
	})
}

func handleUsers(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println(r.Body)

	if r.Method == http.MethodPost {
		CreateNewUser(w, r, db)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		SignIn(w, r, db)
	}
}

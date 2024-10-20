package user

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterUserRoute(router *mux.Router, db *sql.DB) {
	registerUserRoutes(router, db)
}

func registerUserRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		CreateNewUser(w, r, db)
	}).Methods(http.MethodPost)

	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		// GetUserById(w, r, db)
	}).Methods(http.MethodGet)

	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		// UpdateUser(w, r, db)
	}).Methods(http.MethodPut)

	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		// DeleteUser(w, r, db)
	}).Methods(http.MethodDelete)
}

func registerAuthRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		// DeleteUser(w, r, db)
	}).Methods(http.MethodDelete)
	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		DeleteUser(w, r, db)
	}).Methods(http.MethodDelete)
}

func handleLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		SignIn(w, r, db)
	}
}

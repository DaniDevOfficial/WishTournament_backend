package user

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterUserRoute(router *gin.Engine, db *sql.DB) {
	registerUserRoutes(router, db)
	registerAuthRoutes(router, db)
}

func registerUserRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/users/:id", func(c *gin.Context) {
		CreateNewUser()
	})
	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		// UpdateUser(w, r, db)
	}).Methods(http.MethodPut)

	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		// DeleteUser(w, r, db)
	}).Methods(http.MethodDelete)
}

func registerAuthRoutes(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/auth/signup", func(w http.ResponseWriter, r *http.Request) {
		CreateNewUser(w, r, db)
	}).Methods(http.MethodPost)
	router.HandleFunc("/auth/signin", func(w http.ResponseWriter, r *http.Request) {
		SignIn(w, r, db)
	}).Methods(http.MethodPost)
}

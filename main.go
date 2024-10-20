package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"wishtournament/modules/dev"
	"wishtournament/modules/user"
	"wishtournament/util/auth"
	"wishtournament/util/db"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	dbConnection := db.InitDB()
	router := mux.NewRouter()

	user.RegisterUserRoute(router, dbConnection)
	dev.RegisterTicketRoute(router, dbConnection)

	handler := corsMiddleware(router)

	fmt.Println("Server is listening on http://localhost:8000/")
	log.Fatal(http.ListenAndServe("localhost:8000", handler))
}

// corsMiddleware sets the CORS headers to allow all origins
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// goofy ah request for checking if the server supports CORS or sum bs
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		log.Println("New Request Started")
		jwt, err := auth.GetJWTPayloadFromHeader(r)
		log.Println("Request made By: ")
		if err == nil {
			log.Println(jwt)
		} else {
			log.Println("Anonymous")
		}
		log.Println(r.Method)
		next.ServeHTTP(w, r)
	})
}

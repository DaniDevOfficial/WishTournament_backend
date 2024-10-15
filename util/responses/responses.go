package responses

import (
	"encoding/json"
	"log"
	"net/http"
	"wishtournament/util/error"
)

func ResponseWithJSON(w http.ResponseWriter, response interface{}, statusCode int) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("Error converting to JSON:", err)
		error.HttpResponse(w, "Error converting to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Println("Error sending response body:", err)
		error.HttpResponse(w, "Error sending response body", http.StatusInternalServerError)
	}
}

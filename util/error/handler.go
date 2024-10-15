package error

import (
	"log"
	"net/http"
)

func HttpResponse(w http.ResponseWriter, message string, statusCode int) {
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}
	if message == "" {
		message = "An unexpected error occurred"
	}
	log.Println("Error: " + message)
	http.Error(w, message, statusCode)
}

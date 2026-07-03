package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// respond for arbitrary error
func responseWithError(w http.ResponseWriter, statusCode int, msg string) {
	if statusCode > 499 {
		log.Println("responding with 5xx errors: ", msg)
	}
	type errResponse struct {
		Error string `json:"error"` // this line says hey for this var Error you use the json ley called error
	}
	responseWithJSON(w, statusCode, errResponse{
		Error: msg,
	})
}

// this responds for arbitrary JSON
func responseWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	data, err := json.Marshal(payload) //converts go-value of payload into JSON byte slice
	if err != nil {
		log.Println("Failed to marshall JSON response: % v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)

}

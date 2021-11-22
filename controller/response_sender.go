package controller

import (
	"encoding/json"
	"log"
	"net/http"
)

func sendResponseData(w http.ResponseWriter, status int, message string, data interface{}) {
	log.Println(message)
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func sendUnAuthorizedResponse(w http.ResponseWriter) {
	message := "Unauthorized Access"
	log.Println(message)
	w.WriteHeader(401)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

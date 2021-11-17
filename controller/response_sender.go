package controller

import (
	"encoding/json"
	"net/http"
	"tutuplapak-api/model"
)

func sendResponseData(w http.ResponseWriter, status int, message string, data interface{}) {
	var response model.Response
	response.Status = status
	response.Message = message
	response.Data = data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sendUnAuthorizedResponse(w http.ResponseWriter) {
	var response model.Response
	response.Status = 401
	response.Message = "Unauthorized Access"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

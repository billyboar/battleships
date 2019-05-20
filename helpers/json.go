package helpers

import (
	"encoding/json"
	"net/http"
)

// RenderJSON adds JSON content-type header and writes body
func RenderJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	jsonBody, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(statusCode)
	w.Write(jsonBody)
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   error  `json:"error"`
}

func RenderError(w http.ResponseWriter, message string, err error, statusCode int) {
	RenderJSON(w, ErrorResponse{Message: message, Error: err}, statusCode)
}

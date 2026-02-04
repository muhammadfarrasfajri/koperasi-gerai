package utils

import (
	"encoding/json"
	"net/http"

	"github.com/muhammadfarrasfajri/koperasi-gerai/models"
)

// Helper untuk mengirim respon Error
func WriteError(w http.ResponseWriter, statusCode int, errType string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := models.APIResponse{
		Error:   true,
		Type:    errType,
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}

// Helper untuk mengirim respon Sukses
func WriteSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := models.APIResponse{
		Error:   false,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}
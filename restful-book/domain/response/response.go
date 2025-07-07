package response

import (
	"encoding/json"
	"net/http"
)

type StandardResponse struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func JSON(w http.ResponseWriter, data any, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp := StandardResponse{
		Data:    data,
		Message: message,
	}
	json.NewEncoder(w).Encode(resp)
}

func Error(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := StandardResponse{
		Message: "Failed",
		Error:   err.Error(),
	}
	json.NewEncoder(w).Encode(resp)
}

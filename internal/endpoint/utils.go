package endpoint

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func returnJSON(w http.ResponseWriter, status int, res any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		slog.Error("Could Not Encode JSON Response", err)
	}
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func returnError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(ErrorResponse{message})
	if err != nil {
		slog.Error("Could Not Encode Error Response", err)
	}
}

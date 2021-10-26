package handlers

import (
	"encoding/json"
	"net/http"
)

func encodeError(w http.ResponseWriter, status int, err string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse{
		Error: err,
	})
}

func encodeSuccess(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse{
		Success: true,
		Data:    data,
	})
}

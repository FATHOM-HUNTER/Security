package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/form3tech-oss/jwt-go"
)

func encodeError(w http.ResponseWriter, status int, err string) {
	log.Println(err)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse{
		Error: err,
	})
}

func encodeSuccess(w http.ResponseWriter, data interface{}) {
	encodeSuccessHeader(w)
	json.NewEncoder(w).Encode(successResponse{
		Success: true,
		Data:    data,
	})
}

func encodeSuccessHeader(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func parseJWTToken(r *http.Request) string {
	return r.Context().Value("user").(*jwt.Token).Claims.(jwt.MapClaims)["UID"].(string)
}

package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/form3tech-oss/jwt-go"
)

type jwtClaims struct {
	ExpiresAt int64
	Issuer    string
	UID       string
	jwt.Claims
}

func encodeError(w http.ResponseWriter, status int, err string) {
	log.Println(err)
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

func parseJWTToken(r *http.Request) string {
	return r.Context().Value("user").(*jwt.Token).Claims.(jwtClaims).UID
}

package middleware

import (
	"log"
	"net/http"
	"time"

	middleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

const secret = "My Secret"

var jwtMiddleware = middleware.New(middleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	},
	ErrorHandler: errHandler,
	// When set, the middleware verifies that tokens are signed with the specific signing algorithm
	// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
	// Important to avoid security issues described here: https://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries/
	SigningMethod: jwt.SigningMethodHS256,
})

func errHandler(w http.ResponseWriter, r *http.Request, err string) {
	log.Println(err)
}

// GetJWTWrappedNegroni returns a negroni instance wrapping a router
func GetJWTWrappedNegroni(mux *mux.Router) *negroni.Negroni {
	return negroni.New(negroni.HandlerFunc(jwtMiddleware.HandlerWithNext), negroni.Wrap(mux))
}

// ValidateToken validates jwt auth token
func ValidateToken(token string) bool {
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return false
	}
	return true
}

// GenerateToken generates a new jwt token
func GenerateToken(uid string) (string, error) {
	claims := make(jwt.MapClaims, 0)
	claims["ExpiresAt"] = time.Now().Add(time.Hour * 72).Unix()
	claims["Issuer"] = "nameOfWebsiteHere"
	claims["UID"] = uid
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

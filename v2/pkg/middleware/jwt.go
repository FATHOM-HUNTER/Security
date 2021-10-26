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

func customHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Println(r.Header.Get("Authorization"))

	jwtMiddleware.HandlerWithNext(w, r, next)
}

// GetJWTWrappedNegroni returns a negroni instance wrapping a router
func GetJWTWrappedNegroni(mux *mux.Router) *negroni.Negroni {

	return negroni.New(negroni.HandlerFunc(customHandler), negroni.Wrap(mux))
}

// GenerateToken generates a new jwt token
func GenerateToken() (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		Issuer:    "nameOfWebsiteHere",
	}).SignedString([]byte(secret))
}

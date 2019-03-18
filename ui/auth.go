package ui

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	// Time, after which a user gets automatically logged out.
	authCookieTTL = 1 * time.Hour
)

// HMAC hash key for signing authentication cookies. Is randomly generated at web server start.
var hmacKey = make([]byte, 128)

type claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func verifyUser(w http.ResponseWriter, r *http.Request) (string, bool) {
	cookie, err := r.Cookie(authCookieName)
	if err != nil {
		log.Printf("authentication failed: %v\n", err)
		return "", false
	}

	tokenContent := &claims{}

	token, err := jwt.ParseWithClaims(cookie.Value, tokenContent, func(t *jwt.Token) (interface{}, error) { return hmacKey, nil })
	if err != nil {
		log.Printf("authentication failed: %v\n", err)
		return "", false
	}
	if token.Valid {
		username := tokenContent.Username
		renewJWT(w, username)
		return username, true
	}
	log.Println("authentication failed: jwt is invalid")
	return "", false
}

func renewJWT(w http.ResponseWriter, username string) {
	timeout := time.Now().Add(authCookieTTL)
	jwtContent := &claims{username, jwt.StandardClaims{ExpiresAt: timeout.Unix()}}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, jwtContent).SignedString(hmacKey)

	if err != nil {
		log.Printf("creation of json web token not possible: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	setAuthCookie(w, token, timeout)
}

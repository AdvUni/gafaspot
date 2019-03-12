package ui

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	authCookieTTL  = 1 * time.Minute
	authCookieName = "token"
)

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
		setNewAuthCookie(w, username)
		return username, true
	}
	log.Println("authentication failed: jwt is invalid")
	return "", false
}

func setNewAuthCookie(w http.ResponseWriter, username string) {
	timeout := time.Now().Add(authCookieTTL)
	jwtContent := &claims{username, jwt.StandardClaims{ExpiresAt: timeout.Unix()}}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, jwtContent).SignedString(hmacKey)

	if err != nil {
		log.Printf("creation of json web token not possible: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cookie := &http.Cookie{
		Name:     authCookieName,
		Value:    token,
		Expires:  timeout,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

}

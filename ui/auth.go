package ui

import (
	"html/template"
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

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Banner struct {
	Wrong bool
}

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("referer: %v\n", r.Referer())
	log.Printf("path: %v\n", r.URL.Path)
	log.Printf("raw path: %v\n", r.URL.RawPath)
	log.Printf("ur: %v\n", r.RequestURI)

	banner := false
	if r.Referer() == indexpage {
		banner = true
	}
	t, err := template.ParseFiles(indexpageTmpl, topTmpl, bottomTmpl)
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, Banner{banner})
	if err != nil {
		log.Println(err)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("referer: %v\n", r.Referer())
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	username := r.Form.Get("name")
	pass := r.Form.Get("pass")

	//if vault.DoLdapAuthentication(username, pass) {
	if username == "hello" && pass == "world" {
		setNewAuthCookie(w, username)
		http.Redirect(w, r, mainview, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, indexpage, http.StatusSeeOther)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("referer: %v\n", r.Referer())
	_, ok := verifyUser(w, r)
	if ok {
		// return a new, empty cookie which overwrites previous ones and expires immediately
		cookie := &http.Cookie{
			Name:   authCookieName,
			Value:  "",
			MaxAge: -1,
		}
		http.SetCookie(w, cookie)
	}

	// redirect to login page
	http.Redirect(w, r, indexpage, http.StatusSeeOther)
}

func verifyUser(w http.ResponseWriter, r *http.Request) (string, bool) {
	cookie, err := r.Cookie(authCookieName)
	if err != nil {
		log.Printf("authentication failed: %v\n", err)
		return "", false
	}

	tokenContent := &Claims{}

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
	jwtContent := &Claims{username, jwt.StandardClaims{ExpiresAt: timeout.Unix()}}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, jwtContent).SignedString(hmacKey)

	if err != nil {
		log.Printf("creation of json web token not possible: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cookie := &http.Cookie{
		Name:    authCookieName,
		Value:   token,
		Expires: timeout,
	}
	http.SetCookie(w, cookie)

}

package ui

import (
	"crypto/rand"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	mainview     = "/mainview"
	mainviewTmpl = "ui/templates/mainview.html"

	indexpage     = "/"
	indexpageTmpl = "ui/templates/index.html"

	loginpath  = "/login"
	logoutpath = "/logout"
)

var router = mux.NewRouter()

func RunWebserver(addr string) {

	// at webserver start, generate a random key for signing json web tokens for authentication
	// save it to global var (file authentication)
	_, err := rand.Read(hmacKey)
	if err != nil {
		log.Fatalf("could not create key for jwt signing: %v\n", err)
	}

	router.HandleFunc(indexpage, indexPageHandler)
	router.HandleFunc(mainview, mainPageHandler)

	router.HandleFunc(loginpath, loginHandler).Methods(http.MethodPost)
	router.HandleFunc(logoutpath, logoutHandler).Methods(http.MethodPost)

	http.Handle(indexpage, router)
	err = http.ListenAndServe(addr, nil)

	// cause entire program to stop if the server crashes for any reason
	log.Fatalf("webserver crashed: %v\n", err)
}

type Mainviewcontent struct {
	Username string
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := verifyUser(w, r)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	t, err := template.ParseFiles(mainviewTmpl)
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(w, Mainviewcontent{username})
	if err != nil {
		log.Println(err)
	}
}

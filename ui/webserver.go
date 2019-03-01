package ui

import (
	"crypto/rand"
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	topTmpl    = "ui/templates/top.html"
	bottomTmpl = "ui/templates/bottom.html"
	navTmpl    = "ui/templates/nav.html"

	indexpage     = "/"
	indexpageTmpl = "ui/templates/index.html"

	loginpath  = "/login"
	logoutpath = "/logout"

	mainview     = "/mainview"
	mainviewTmpl = "ui/templates/mainview.html"

	personalview     = "/personal"
	personalviewTmpl = "ui/templates/personalview.html"

	credsview = "/personal/creds"

	reservationform     = "/newreservation"
	reservationformTmpl = "ui/templates/newreservation.html"
)

var (
	router = mux.NewRouter()
	db     *sql.DB
)

func RunWebserver(database *sql.DB, addr string) {

	db = database

	// at webserver start, generate a random key for signing json web tokens for authentication
	// save it to global var (file authentication)
	_, err := rand.Read(hmacKey)
	if err != nil {
		log.Fatalf("could not create key for jwt signing: %v\n", err)
	}

	router.HandleFunc(indexpage, indexPageHandler)

	router.HandleFunc(loginpath, loginHandler).Methods(http.MethodPost)
	router.HandleFunc(logoutpath, logoutHandler).Methods(http.MethodPost)

	router.HandleFunc(mainview, mainPageHandler)
	router.HandleFunc(personalview, personalPageHandler)
	//router.HandleFunc(credsview, credsPageHandler)

	router.HandleFunc(reservationform, reservationHandler)

	http.Handle(indexpage, router)
	err = http.ListenAndServe(addr, nil)

	// cause entire program to stop if the server crashes for any reason
	log.Fatalf("webserver crashed: %v\n", err)
}

type Mainviewcontent struct {
	Username string
}

func reservationHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := verifyUser(w, r)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	t, err := template.ParseFiles(reservationformTmpl, topTmpl, bottomTmpl, navTmpl)
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(w, Mainviewcontent{username})
	if err != nil {
		log.Println(err)
	}

}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := verifyUser(w, r)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	t, err := template.ParseFiles(mainviewTmpl, topTmpl, bottomTmpl, navTmpl)
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(w, Mainviewcontent{username})
	if err != nil {
		log.Println(err)
	}
}

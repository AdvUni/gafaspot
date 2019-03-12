package ui

import (
	"crypto/rand"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/gorilla/mux"
	"gitlab-vs.informatik.uni-ulm.de/gafaspot/database"
	"gitlab-vs.informatik.uni-ulm.de/gafaspot/util"
)

const (
	loginpage        = "/"
	login            = "/login"
	logout           = "/logout"
	mainview         = "/mainview"
	personalview     = "/personal"
	credsview        = "/personal/creds"
	reservationform  = "/newreservation/{env}"
	reserve          = "/reserve"
	abortreservation = "/abortreservation"
)

var (
	environments []util.Environment
	envHasSSHMap map[string]bool

	loginformTmpl       *template.Template
	mainviewTmpl        *template.Template
	personalviewTmpl    *template.Template
	reservationformTmpl *template.Template
	reservesuccessTmpl  *template.Template
)

func init() {
	// generate a random key for signing json web tokens for authentication. Save it to global var (file authentication)
	_, err := rand.Read(hmacKey)
	if err != nil {
		log.Fatalf("could not create key for jwt signing: %v\n", err)
	}

	// pre-assembling and caching of all the page templates
	const (
		topTmplFile             = "ui/templates/top.html"
		bottomTmplFile          = "ui/templates/bottom.html"
		navTmplFile             = "ui/templates/nav.html"
		loginformTmplFile       = "ui/templates/login.html"
		mainviewTmplFile        = "ui/templates/mainview.html"
		personalviewTmplFile    = "ui/templates/personalview.html"
		reservationformTmplFile = "ui/templates/newreservation.html"
		reservesuccessTmplFile  = "ui/templates/reservesuccess.html"
	)
	loginformTmpl, err = template.ParseFiles(loginformTmplFile, topTmplFile, bottomTmplFile)
	if err != nil {
		log.Fatal(err)
	}
	mainviewTmpl, err = template.New(path.Base(mainviewTmplFile)).Funcs(template.FuncMap{
		"formatDatetime": func(t time.Time) string { return t.Format(util.TimeLayout) },
	}).ParseFiles(mainviewTmplFile, topTmplFile, bottomTmplFile, navTmplFile)
	if err != nil {
		log.Fatal(err)
	}
	personalviewTmpl, err = template.New(path.Base(personalviewTmplFile)).Funcs(template.FuncMap{
		"formatDatetime": func(t time.Time) string { return t.Format(util.TimeLayout) },
	}).ParseFiles(personalviewTmplFile, topTmplFile, bottomTmplFile, navTmplFile)
	if err != nil {
		log.Fatal(err)
	}
	reservationformTmpl, err = template.ParseFiles(reservationformTmplFile, topTmplFile, bottomTmplFile, navTmplFile)
	if err != nil {
		log.Fatal(err)
	}
	reservesuccessTmpl, err = template.ParseFiles(reservesuccessTmplFile, topTmplFile, bottomTmplFile, navTmplFile)
	if err != nil {
		log.Fatal(err)
	}
}

func RunWebserver(addr string) {

	// fetch static information about environments from database
	environments, envHasSSHMap = database.GetEnvironments()

	// create router and register all paths
	router := mux.NewRouter()

	router.HandleFunc(loginpage, loginPageHandler)
	router.HandleFunc(login, loginHandler).Methods(http.MethodPost)
	router.HandleFunc(logout, logoutHandler).Methods(http.MethodPost)
	router.HandleFunc(mainview, mainPageHandler)
	router.HandleFunc(personalview, personalPageHandler)
	router.HandleFunc(credsview, credsPageHandler)
	router.HandleFunc(reservationform, newreservationPageHandler)
	router.HandleFunc(reserve, reserveHandler)
	router.HandleFunc(abortreservation, abortreservationHandler)
	//router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("ui/templates/js"))))

	// start web server
	http.Handle(loginpage, router)
	err := http.ListenAndServe(addr, nil)
	// cause entire program to stop if the server crashes for any reason
	log.Fatalf("webserver crashed: %v\n", err)
}

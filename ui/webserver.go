package ui

import (
	"crypto/rand"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

const (
	loginpage       = "/"
	login           = "/login"
	logout          = "/logout"
	mainview        = "/mainview"
	personalview    = "/personal"
	credsview       = "/personal/creds"
	reservationform = "/newreservation/{env}"
	reserve         = "/reserve"
)

var (
	db      *sql.DB
	envList []env
	envMap  map[string]env

	loginformTmpl       *template.Template
	mainviewTmpl        *template.Template
	personalviewTmpl    *template.Template
	reservationformTmpl *template.Template
)

type env struct {
	Name        string
	NamePlain   string
	HasSSH      bool
	Description string
}

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
	)
	loginformTmpl, err = template.ParseFiles(loginformTmplFile, topTmplFile, bottomTmplFile)
	if err != nil {
		log.Fatal(err)
	}
	mainviewTmpl, err = template.ParseFiles(mainviewTmplFile, topTmplFile, bottomTmplFile, navTmplFile)
	if err != nil {
		log.Fatal(err)
	}
	personalviewTmpl, err = template.ParseFiles(personalviewTmplFile, topTmplFile, bottomTmplFile, navTmplFile)
	if err != nil {
		log.Fatal(err)
	}
	reservationformTmpl, err = template.ParseFiles(reservationformTmplFile, topTmplFile, bottomTmplFile, navTmplFile)
	if err != nil {
		log.Fatal(err)
	}
}

func RunWebserver(database *sql.DB, addr string) {

	db = database

	// fetch static information about environments from database
	envList, envMap = getEnvironments(db)

	// create router and register all paths
	router := mux.NewRouter()

	router.HandleFunc(loginpage, loginPageHandler)
	router.HandleFunc(login, loginHandler).Methods(http.MethodPost)
	router.HandleFunc(logout, logoutHandler).Methods(http.MethodPost)
	router.HandleFunc(mainview, mainPageHandler)
	router.HandleFunc(personalview, personalPageHandler)
	//router.HandleFunc(credsview, credsPageHandler)
	router.HandleFunc(reservationform, newreservationPageHandler)
	router.HandleFunc(reserve, reserveHandler)
	//router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("ui/templates/js"))))

	// start web server
	http.Handle(loginpage, router)
	err := http.ListenAndServe(addr, nil)
	// cause entire program to stop if the server crashes for any reason
	log.Fatalf("webserver crashed: %v\n", err)
}

// createPlainIdentifier replaces all characters which are not ascii letters oder numbers through an underscore
func createPlainIdentifier(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	return strings.ToLower(re.ReplaceAllString(name, "_"))
}

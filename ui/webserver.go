package ui

import (
	"crypto/rand"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

const (
	indexpage = "/"

	loginpath  = "/login"
	logoutpath = "/logout"

	mainview = "/mainview"

	personalview = "/personal"

	credsview = "/personal/creds"

	reservationform = "/newreservation/{env}"
)

var (
	db   *sql.DB
	envs = []string{"DEMO 0", "DEMO 1", "DEMO 2", "DEMO 3", "DEMO 4", "DEMO 5"}

	loginformTmpl       *template.Template
	mainviewTmpl        *template.Template
	personalviewTmpl    *template.Template
	reservationformTmpl *template.Template
)

func RunWebserver(database *sql.DB, addr string) {

	db = database

	// at webserver start, generate a random key for signing json web tokens for authentication
	// save it to global var (file authentication)
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
	mainviewTmpl, err = template.New(path.Base(mainviewTmplFile)).Funcs(template.FuncMap{
		"plainString": createPlainIdentifier,
	}).ParseFiles(mainviewTmplFile, topTmplFile, bottomTmplFile, navTmplFile)
	if err != nil {
		log.Fatal(err)
	}
	personalviewTmpl, err = template.ParseFiles(personalviewTmplFile, topTmplFile, bottomTmplFile, navTmplFile)
	if err != nil {
		log.Fatal(err)
	}
	reservationformTmpl, err = template.New(path.Base(reservationformTmplFile)).Funcs(template.FuncMap{
		"plainString": createPlainIdentifier,
	}).ParseFiles(reservationformTmplFile, topTmplFile, bottomTmplFile, navTmplFile)
	if err != nil {
		log.Fatal(err)
	}

	// create router and register all paths
	router := mux.NewRouter()

	router.HandleFunc(indexpage, indexPageHandler)

	router.HandleFunc(loginpath, loginHandler).Methods(http.MethodPost)
	router.HandleFunc(logoutpath, logoutHandler).Methods(http.MethodPost)

	router.HandleFunc(mainview, mainPageHandler)
	router.HandleFunc(personalview, personalPageHandler)
	//router.HandleFunc(credsview, credsPageHandler)

	router.HandleFunc(reservationform, reservationHandler)

	//router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("ui/templates/js"))))

	http.Handle(indexpage, router)
	err = http.ListenAndServe(addr, nil)

	// cause entire program to stop if the server crashes for any reason
	log.Fatalf("webserver crashed: %v\n", err)
}

type Mainviewcontent struct {
	Username string
	Envs     []env
}

type Reservationcontent struct {
	Username    string
	SelectedEnv string
	Envs        []string
	SSHmissing  bool
}

type env struct {
	Name        string
	Description string
}

func reservationHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := verifyUser(w, r)
	if !ok {
		http.NotFound(w, r)
		return
	}

	selectedEnv := mux.Vars(r)["env"]

	// TODO: Check if ssh key is needed and if user has one
	err := reservationformTmpl.Execute(w, Reservationcontent{username, selectedEnv, envs, true})
	if err != nil {
		log.Println(err)
	}
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("referer: %v\n", r.Referer())
	username, ok := verifyUser(w, r)
	if !ok {
		http.NotFound(w, r)
		return
	}

	envs := []env{env{"DEMO 0", "zero"}, env{"DEMO 1", "first"}, env{"DEMO 2", "second"}, env{"DEMO 3", "third"}, env{"DEMO 4", "fourth"}, env{"DEMO 5", "last"}}
	err := mainviewTmpl.Execute(w, Mainviewcontent{username, envs})
	if err != nil {
		log.Println(err)
	}
}

// createPlainIdentifier replaces all characters which are not ascii letters oder numbers through an underscore
func createPlainIdentifier(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	return strings.ToLower(re.ReplaceAllString(name, "_"))
}

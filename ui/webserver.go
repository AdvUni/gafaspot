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
	envs map[string]env

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

type reservation struct {
	ID      int
	User    string
	Env     string
	Start   string
	End     string
	Subject string
	Labels  string
}

type envReservations struct {
	Env                  env
	ReservationsUpcoming []reservation
	ReservationsActive   []reservation
	ReservationsExpired  []reservation
}

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

	// fetch static information about environments from database
	// TODO: get data from database
	envs = make(map[string]env)
	envs[createPlainIdentifier("DEMO 0")] = env{"DEMO 0", createPlainIdentifier("DEMO 0"), true, "zero"}
	envs[createPlainIdentifier("DEMO 1")] = env{"DEMO 1", createPlainIdentifier("DEMO 1"), false, "first"}
	envs[createPlainIdentifier("DEMO 2")] = env{"DEMO 2", createPlainIdentifier("DEMO 2"), true, "second"}
	envs[createPlainIdentifier("DEMO 3")] = env{"DEMO 3", createPlainIdentifier("DEMO 3"), false, "third"}
	envs[createPlainIdentifier("DEMO 4")] = env{"DEMO 4", createPlainIdentifier("DEMO 4"), true, "fourth"}
	envs[createPlainIdentifier("DEMO 5")] = env{"DEMO 5", createPlainIdentifier("DEMO 5"), false, "last"}

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

func reservationHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := verifyUser(w, r)
	if !ok {
		http.NotFound(w, r)
		return
	}

	selectedEnv := mux.Vars(r)["env"]

	// TODO: Check if ssh key is needed and if user has one
	sshMissing := !envs[selectedEnv].HasSSH || false

	err := reservationformTmpl.Execute(w, map[string]interface{}{"Username": username, "Envs": envs, "Selected": selectedEnv, "SSHmissing": sshMissing})
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

	envReservationsList := []envReservations{}
	for _, env := range envs {
		// TODO: fetch reservations from database
		res := reservation{0, "user1", "demo0", "2000-01-01", "2000-01-01", "no subject", ""}
		list := []reservation{res, res, res}
		envReservations := envReservations{env, list, list, list}

		envReservationsList = append(envReservationsList, envReservations)
	}

	err := mainviewTmpl.Execute(w, map[string]interface{}{"Username": username, "Envcontent": envReservationsList})
	if err != nil {
		log.Println(err)
	}
}

// createPlainIdentifier replaces all characters which are not ascii letters oder numbers through an underscore
func createPlainIdentifier(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	return strings.ToLower(re.ReplaceAllString(name, "_"))
}

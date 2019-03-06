package ui

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

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

func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("referer: %v\n", r.Referer())
	log.Printf("path: %v\n", r.URL.Path)
	log.Printf("raw path: %v\n", r.URL.RawPath)
	log.Printf("ur: %v\n", r.RequestURI)

	banner := false
	if r.Referer() == loginpage {
		banner = true
	}
	err := loginformTmpl.Execute(w, map[string]interface{}{"ShowBanner": banner})
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

func personalPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("referer: %v\n", r.Referer())
	username, ok := verifyUser(w, r)
	if !ok {
		http.NotFound(w, r)
		return
	}

	res := reservation{0, "user1", "demo0", "2000-01-01", "2000-01-01", "no subject", ""}
	list := []reservation{res, res, res}
	err := personalviewTmpl.Execute(w, map[string]interface{}{"Username": username, "SSHkey": "blank", "ReservationsUpcoming": list, "ReservationsActive": list, "ReservationsExpired": list})
	if err != nil {
		log.Println(err)
	}
}

func credsPageHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := verifyUser(w, r)
	if !ok {
		http.NotFound(w, r)
		return
	}

	stmt, err := db.Prepare("SELECT env_name, start, end FROM reservations WHERE (status='active') AND (username=?);")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	resp, err := stmt.Query(username)
	log.Println(resp)
}

func newreservationPageHandler(w http.ResponseWriter, r *http.Request) {
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
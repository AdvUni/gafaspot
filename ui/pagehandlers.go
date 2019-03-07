package ui

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type reservation struct {
	ID      int
	Status  string
	User    string
	EnvName string
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

	for _, env := range envList {
		// TODO: fetch reservations from database
		reservations := getEnvReservations(db, env.Name)
		var upcoming, active, expired []reservation
		for _, r := range reservations {
			switch r.Status {
			case "upcoming":
				upcoming = append(upcoming, r)
			case "active":
				active = append(active, r)
			case "expired":
				expired = append(expired, r)
			}
		}
		envReservationsList = append(envReservationsList, envReservations{env, upcoming, active, expired})
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

	res := reservation{0, "somestatus", "user1", "demo0", "2000-01-01", "2000-01-01", "no subject", ""}
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

	selectedEnvPlainName := mux.Vars(r)["env"]
	selectedEnv, ok := envMap[selectedEnvPlainName]
	if !ok {
		fmt.Fprint(w, "environment in url does not exist")
		return
	}

	// TODO: Check if ssh key is needed and if user has one
	sshMissing := !selectedEnv.HasSSH || false

	err := reservationformTmpl.Execute(w, map[string]interface{}{"Username": username, "Envs": envList, "Selected": selectedEnvPlainName, "SSHmissing": sshMissing})
	if err != nil {
		log.Println(err)
	}
}

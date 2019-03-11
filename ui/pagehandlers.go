package ui

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"gitlab-vs.informatik.uni-ulm.de/gafaspot/database"
	"gitlab-vs.informatik.uni-ulm.de/gafaspot/util"
	"gitlab-vs.informatik.uni-ulm.de/gafaspot/vault"

	"github.com/gorilla/mux"
)

type envReservations struct {
	Env                  util.Environment
	ReservationsUpcoming []util.Reservation
	ReservationsActive   []util.Reservation
	ReservationsExpired  []util.Reservation
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

func sortReservations(reservations []util.Reservation) ([]util.Reservation, []util.Reservation, []util.Reservation) {
	// sort reservation list by start date
	sort.Slice(reservations, func(i, j int) bool {
		return reservations[i].Start.Before(reservations[j].Start)
	})
	// split list into three sub lists
	var upcoming, active, expired []util.Reservation
	for _, r := range reservations {
		switch r.Status {
		case "upcoming":
			upcoming = append(upcoming, r)
		case "active":
			active = append(active, r)
		case "expired":
			expired = append(expired, r)
		case "error":
			log.Println("reservation with status error is not displayed")
		}
	}
	return upcoming, active, expired
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("referer: %v\n", r.Referer())
	username, ok := verifyUser(w, r)
	if !ok {
		http.NotFound(w, r)
		return
	}
	envReservationsList := []envReservations{}

	for _, env := range environments {
		upcoming, active, expired := sortReservations(database.GetEnvReservations(env.NiceName))
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

	sshEntry, ok := database.GetUserSSH(username)
	if !ok {
		sshEntry = "no key yet"
	}

	upcoming, active, expired := sortReservations(database.GetUserReservations(username))
	err := personalviewTmpl.Execute(w, map[string]interface{}{"Username": username, "SSHkey": sshEntry, "ReservationsUpcoming": upcoming, "ReservationsActive": active, "ReservationsExpired": expired})
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
	envNames := database.GetUserActiveReservationEnv(username)
	log.Println(envNames)

	for _, env := range envNames {
		creds, err := vault.ReadCredentials(env, vault.CreateVaultToken())
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(w, "creds for environment '%v':\n%v\n", env, creds)
	}
}

func newreservationPageHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := verifyUser(w, r)
	if !ok {
		http.NotFound(w, r)
		return
	}

	selectedEnvPlainName := mux.Vars(r)["env"]
	envHasSSH, ok := envHasSSHMap[selectedEnvPlainName]
	if !ok {
		fmt.Fprint(w, "environment in url does not exist")
		return
	}
	sshMissing := envHasSSH && !database.UserHasSSH(username)

	err := reservationformTmpl.Execute(w, map[string]interface{}{"Username": username, "Envs": environments, "Selected": selectedEnvPlainName, "SSHmissing": sshMissing})
	if err != nil {
		log.Println(err)
	}
}

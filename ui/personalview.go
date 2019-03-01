package ui

import (
	"html/template"
	"log"
	"net/http"
)

func credsPageHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := verifyUser(w, r)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
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

type PersonalviewContent struct {
	Username             string
	SSH                  string
	ReservationsUpcoming []string
	ReservationsActive   []string
	ReservationsExpired  []string
}

func personalPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("referer: %v\n", r.Referer())
	username, ok := verifyUser(w, r)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	t, err := template.ParseFiles(personalviewTmpl, topTmpl, bottomTmpl, navTmpl)
	if err != nil {
		log.Fatal(err)
	}
	upcoming := []string{"first reservation", "second reservation", "third reservation"}
	active := []string{"first reservation", "second reservation", "third reservation"}
	expired := []string{"first reservation", "second reservation", "third reservation"}
	err = t.Execute(w, PersonalviewContent{username, "blank", upcoming, active, expired})
	if err != nil {
		log.Println(err)
	}
}

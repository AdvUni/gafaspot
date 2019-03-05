package ui

import (
	"html/template"
	"log"
	"net/http"
)

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

type reservation struct {
	ID      int
	User    string
	Env     string
	Start   string
	End     string
	Subject string
	Labels  string
}

type PersonalviewContent struct {
	Username             string
	SSH                  string
	ReservationsUpcoming []reservation
	ReservationsActive   []reservation
	ReservationsExpired  []reservation
}

func personalPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("referer: %v\n", r.Referer())
	username, ok := verifyUser(w, r)
	if !ok {
		http.NotFound(w, r)
		return
	}

	t, err := template.ParseFiles(personalviewTmpl, topTmpl, bottomTmpl, navTmpl)
	if err != nil {
		log.Fatal(err)
	}
	res := reservation{0, "user1", "demo0", "2000-01-01", "2000-01-01", "no subject", ""}
	list := []reservation{res, res, res}
	err = t.Execute(w, PersonalviewContent{username, "blank", list, list, list})
	if err != nil {
		log.Println(err)
	}
}

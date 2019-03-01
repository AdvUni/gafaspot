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
	Username     string
	SSH          string
	Reservations string
}

func personalPageHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := verifyUser(w, r)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	t, err := template.ParseFiles(personalviewTmpl, topTmpl, bottomTmpl, navTmpl)
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, PersonalviewContent{username, "blank", string(template.HTML("<li class='list-group-item'>blank</li>"))})
	if err != nil {
		log.Println(err)
	}
}

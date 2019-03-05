package ui

import (
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

package ui

import (
	"html/template"
	"log"
	"net/http"

	"gitlab-vs.informatik.uni-ulm.de/gafaspot/vault"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("referer: %v\n", r.Referer())
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	username := r.Form.Get("name")
	pass := r.Form.Get("pass")

	if vault.DoLdapAuthentication(username, pass) {
		setNewAuthCookie(w, username)
		http.Redirect(w, r, mainview, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, loginpage, http.StatusSeeOther)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("referer: %v\n", r.Referer())
	_, ok := verifyUser(w, r)
	if ok {
		// return a new, empty cookie which overwrites previous ones and expires immediately
		cookie := &http.Cookie{
			Name:   authCookieName,
			Value:  "",
			MaxAge: -1,
		}
		http.SetCookie(w, cookie)
	}

	// redirect to login page
	http.Redirect(w, r, loginpage, http.StatusSeeOther)
}

func reserveHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := verifyUser(w, r)
	if !ok {
		http.NotFound(w, r)
		return
	}
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(username)
	env := template.HTMLEscapeString(r.Form.Get("env"))
	startdate := template.HTMLEscapeString(r.Form.Get("startdate"))
	starttime := template.HTMLEscapeString(r.Form.Get("starttime"))
	subject := template.HTMLEscapeString(r.Form.Get("sub"))
	log.Printf("env: %v, startdate: %v, starttime: %v, subject: %v\n", env, startdate, starttime, subject)

	//CreateReservation(db, username)
}

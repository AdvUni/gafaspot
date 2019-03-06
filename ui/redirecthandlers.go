package ui

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"gitlab-vs.informatik.uni-ulm.de/gafaspot/constants"

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
	env, ok := envs[template.HTMLEscapeString(r.Form.Get("env"))]
	if !ok {
		fmt.Fprintf(w, "environment invalid")
		return
	}
	envName := env.Name

	startstring := template.HTMLEscapeString(r.Form.Get("startdate")) + " " + template.HTMLEscapeString(r.Form.Get("starttime"))
	start, err := time.ParseInLocation(constants.TimeLayout, startstring, time.Local)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "start date/time malformed")
		return
	}

	endstring := template.HTMLEscapeString(r.Form.Get("enddate")) + " " + template.HTMLEscapeString(r.Form.Get("endtime"))
	end, err := time.ParseInLocation(constants.TimeLayout, endstring, time.Local)
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, "end date/time malformed")
		return
	}

	subject := template.HTMLEscapeString(r.Form.Get("sub"))

	log.Printf("env: %v, start: %v, end: %v, subject: %v\n", env, start, end, subject)

	err = CreateReservation(db, username, envName, subject, "", start, end)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	fmt.Fprint(w, "success!")
}
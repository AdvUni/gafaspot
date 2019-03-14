package ui

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/ssh"

	"gitlab-vs.informatik.uni-ulm.de/gafaspot/database"
	"gitlab-vs.informatik.uni-ulm.de/gafaspot/vault"
)

func redirectNotAuthenticated(w http.ResponseWriter, r *http.Request) {
	redirectShowLoginError(w, r, "You are not (longer) logged in")
}

func redirectShowLoginError(w http.ResponseWriter, r *http.Request, errormessage string) {
	setErrorCookie(w, errormessage)
	http.Redirect(w, r, loginpage, http.StatusSeeOther)
}

func redirectLogoutSuccessful(w http.ResponseWriter, r *http.Request) {
	setInfoCookie(w, "Successfully logged out")
	http.Redirect(w, r, loginpage, http.StatusSeeOther)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	username := r.Form.Get("name")
	pass := r.Form.Get("pass")

	if !vault.DoLdapAuthentication(username, pass) {
		redirectShowLoginError(w, r, "Invalid credentials")
		return
	}

	// each time a user logs in, update the TTL for his database entry
	database.RefreshDeletionDate(username)

	renewJWT(w, username)
	http.Redirect(w, r, mainview, http.StatusSeeOther)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	_, ok := verifyUser(w, r)
	if ok {
		invalidateCookie(w, authCookieName)
	}

	// redirect to login page
	redirectLogoutSuccessful(w, r)
}

func redirectInvalidSubmission(w http.ResponseWriter, r *http.Request, errormessage string) {
	setErrorCookie(w, errormessage)
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func abortreservationHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := verifyUser(w, r)
	if !ok {
		redirectNotAuthenticated(w, r)
		return
	}
	err := r.ParseForm()
	if err != nil {
		log.Printf("could not get parameter id from abort reservation request: %v\n", err)
		return
	}

	reservationID, err := strconv.Atoi(template.HTMLEscapeString(r.Form.Get("id")))
	if err != nil {
		log.Printf("abortreservation request passes an id which is not comparable to int: %v\n", template.HTMLEscapeString(r.Form.Get("id")))
		return
	}
	database.AbortReservation(username, reservationID)
	// return to personal view
	http.Redirect(w, r, personalview, http.StatusSeeOther)
}

func uploadkeyHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := verifyUser(w, r)
	if !ok {
		redirectNotAuthenticated(w, r)
		return
	}
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	sshPubkey := []byte(r.Form.Get("ssh-paste-field"))
	_, _, _, _, err = ssh.ParseAuthorizedKey(sshPubkey)
	if err != nil {
		redirectInvalidSubmission(w, r, fmt.Sprintf("You did not enter a valid key (%v)", err))
		return
	}
	fmt.Fprint(w, "Key is valid!")

	database.SaveUserSSH(username, sshPubkey)

	log.Println(username)
}

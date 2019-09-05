// Copyright 2019, Advanced UniByte GmbH.
// Author Marie Lohbeck.
//
// This file is part of Gafaspot.
//
// Gafaspot is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Gafaspot is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Gafaspot.  If not, see <https://www.gnu.org/licenses/>.

package ui

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/AdvUni/gafaspot/database"
	"github.com/AdvUni/gafaspot/vault"
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
		logger.Warning(err)
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
		logger.Warningf("could not get parameter id from abort reservation request: %v\n", err)
		return
	}

	reservationID, err := strconv.Atoi(template.HTMLEscapeString(r.Form.Get("id")))
	if err != nil {
		logger.Warningf("abortreservation request passes an id which is not comparable to int: %v\n", template.HTMLEscapeString(r.Form.Get("id")))
		return
	}
	database.AbortReservation(username, reservationID)
	// return to personal view
	http.Redirect(w, r, personalview, http.StatusSeeOther)
}

func deletekeyHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := verifyUser(w, r)
	if !ok {
		redirectNotAuthenticated(w, r)
		return
	}
	database.DeleteUserSSH(username)
	http.Redirect(w, r, personalview, http.StatusSeeOther)
}

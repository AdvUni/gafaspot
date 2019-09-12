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
	"fmt"
	"html/template"
	"net/http"
	"net/mail"
	"sort"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/AdvUni/gafaspot/database"
	"github.com/AdvUni/gafaspot/email"
	"github.com/AdvUni/gafaspot/util"
	"github.com/AdvUni/gafaspot/vault"
	"github.com/gorilla/mux"
)

// envReservations is a struct used for passing data to main view
type envReservations struct {
	Env          util.Environment
	Reservations []util.Reservation
}

// reservationNiceName is a struct used for passing reservation data to personal view
type reservationNiceName struct {
	EnvNiceName  string
	ID           int
	Status       string
	User         string
	EnvPlainName string
	Start        time.Time
	End          time.Time
	Subject      string
	Labels       string
}

func newReservationNiceName(r util.Reservation) reservationNiceName {
	return reservationNiceName{
		environmentsMap[r.EnvPlainName].NiceName,
		r.ID,
		r.Status,
		r.User,
		r.EnvPlainName,
		r.Start,
		r.End,
		r.Subject,
		r.Labels,
	}
}

func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	errormessage := readErrorCookie(w, r)
	infomessage := readInfoCookie(w, r)

	err := loginformTmpl.Execute(w, map[string]interface{}{"Error": errormessage, "Info": infomessage})
	if err != nil {
		logger.Error(err)
	}
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := verifyUser(w, r)
	if !ok {
		redirectNotAuthenticated(w, r)
		return
	}
	var envReservationsList []envReservations

	for _, env := range environments {

		reservations := database.GetEnvReservations(env.PlainName)
		// sort reservations
		sort.Slice(reservations, func(i, j int) bool {
			return reservations[i].Start.Before(reservations[j].Start)
		})

		envReservationsList = append(envReservationsList, envReservations{env, reservations})
	}

	err := mainviewTmpl.Execute(w, map[string]interface{}{"Username": username, "Envcontent": envReservationsList})
	if err != nil {
		logger.Error(err)
	}
}

func personalPageHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := verifyUser(w, r)
	if !ok {
		redirectNotAuthenticated(w, r)
		return
	}

	sshEntry, ok := database.GetUserSSH(username)
	if !ok {
		sshEntry = ""
	}

	mail, ok := database.GetUserEmail(username)
	if !ok {
		mail = ""
	}

	reservations := database.GetUserReservations(username)
	// sort reservations
	sort.Slice(reservations, func(i, j int) bool {
		return reservations[i].Start.Before(reservations[j].Start)
	})
	var resNice []reservationNiceName
	for _, r := range reservations {
		resNice = append(resNice, newReservationNiceName(r))
	}

	err := personalviewTmpl.Execute(w, map[string]interface{}{
		"Username":      username,
		"SSHkey":        sshEntry,
		"EmailDisabled": !email.MailingEnabled,
		"Email":         mail,
		"Reservations":  resNice,
	})
	if err != nil {
		logger.Error(err)
	}
}

func credsPageHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := verifyUser(w, r)
	if !ok {
		redirectNotAuthenticated(w, r)
		return
	}
	credsData := database.CollectUserCreds(username, vault.ReadCredentials)

	credsviewTmpl.Execute(w, map[string]interface{}{"Username": username, "CredsData": credsData})
}

func newreservationPageHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := verifyUser(w, r)
	if !ok {
		redirectNotAuthenticated(w, r)
		return
	}

	errormessage := readErrorCookie(w, r)

	selectedEnvPlainName := mux.Vars(r)["env"]
	env, ok := environmentsMap[selectedEnvPlainName]
	if !ok {
		fmt.Fprint(w, "environment in url does not exist")
		return
	}
	sshMissing := env.HasSSH && !database.UserHasSSH(username)
	emailMissing := !database.UserHasEmail(username)

	err := reservationformTmpl.Execute(w, map[string]interface{}{
		"Username":      username,
		"Envs":          environments,
		"Selected":      selectedEnvPlainName,
		"SSHmissing":    sshMissing,
		"EmailDisabled": !email.MailingEnabled,
		"EmailMissing":  emailMissing,
		"Error":         errormessage,
	})
	if err != nil {
		logger.Error(err)
	}
}

func reserveHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := verifyUser(w, r)
	if !ok {
		redirectNotAuthenticated(w, r)
		return
	}
	err := r.ParseForm()
	if err != nil {
		logger.Warning(err)
		return
	}

	var reservation util.Reservation

	reservation.User = username

	// get environment from form
	reservation.EnvPlainName = template.HTMLEscapeString(r.Form.Get("env"))
	if reservation.EnvPlainName == "" {
		logger.Debugf("reserve handler received reservation with invalid environment: %v", err)
		redirectInvalidSubmission(w, r, "environment invalid")
		return
	}

	// get start from form
	startstring := template.HTMLEscapeString(r.Form.Get("startdate")) + " " + template.HTMLEscapeString(r.Form.Get("starttime"))
	reservation.Start, err = time.ParseInLocation(util.TimeLayout, startstring, time.Local)
	if err != nil {
		logger.Debugf("reserve handler received reservation with malformed date/time submission: %v", err)
		redirectInvalidSubmission(w, r, "start date/time malformed")
		return
	}

	// get end from form
	endstring := template.HTMLEscapeString(r.Form.Get("enddate")) + " " + template.HTMLEscapeString(r.Form.Get("endtime"))
	reservation.End, err = time.ParseInLocation(util.TimeLayout, endstring, time.Local)
	if err != nil {
		logger.Debugf("reserve handler received reservation with malformed date/time submission: %v", err)
		redirectInvalidSubmission(w, r, "end date/time malformed")
		return
	}

	// get subject from form
	reservation.Subject = template.HTMLEscapeString(r.Form.Get("sub"))

	// get email checkboxes from form
	if r.Form.Get("startmail") != "" {
		reservation.SendStartMail = true
	}
	if r.Form.Get("endmail") != "" {
		reservation.SendEndMail = true
	}

	err = database.CreateReservation(reservation)
	if err != nil {
		logger.Debugf("reserve handler received invalid reservation: %v", err)
		redirectInvalidSubmission(w, r, err.Error())
		return
	}

	err = reservesuccessTmpl.Execute(w, newReservationNiceName(reservation))
	if err != nil {
		logger.Error(err)
	}
}

func addkeyPageHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := verifyUser(w, r)
	if !ok {
		redirectNotAuthenticated(w, r)
		return
	}

	errormessage := readErrorCookie(w, r)

	err := addkeyformTmpl.Execute(w, map[string]interface{}{"Username": username, "Error": errormessage})
	if err != nil {
		logger.Error(err)
	}
}

func uploadkeyHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := verifyUser(w, r)
	if !ok {
		redirectNotAuthenticated(w, r)
		return
	}
	err := r.ParseForm()
	if err != nil {
		logger.Warning(err)
		return
	}

	sshString := r.Form.Get("ssh-paste-field")
	sshPubkey := []byte(sshString)
	_, _, _, _, err = ssh.ParseAuthorizedKey(sshPubkey)
	if err != nil {
		redirectInvalidSubmission(w, r, "You did not enter a valid key")
		return
	}

	database.SaveUserSSH(username, sshPubkey)

	err = addkeysuccessTmpl.Execute(w, map[string]interface{}{"Username": username, "SSHkey": sshString})
	if err != nil {
		logger.Error(err)
	}
}

func addmailPageHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := verifyUser(w, r)
	if !ok {
		redirectNotAuthenticated(w, r)
		return
	}

	errormessage := readErrorCookie(w, r)

	err := addmailformTmpl.Execute(w, map[string]interface{}{"Username": username, "Error": errormessage})
	if err != nil {
		logger.Error(err)
	}
}

func uploadmailHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := verifyUser(w, r)
	if !ok {
		redirectNotAuthenticated(w, r)
		return
	}

	err := r.ParseForm()
	if err != nil {
		logger.Warning(err)
		return
	}

	email := template.HTMLEscapeString(r.Form.Get("email-field"))

	_, err = mail.ParseAddress(email)
	if err != nil {
		redirectInvalidSubmission(w, r, "You did not enter a valid e-mail address")
		return
	}

	database.SaveUserEmail(username, email)

	err = addmailsuccessTmpl.Execute(w, map[string]interface{}{"Username": username, "Email": email})
	if err != nil {
		logger.Error(err)
	}
}

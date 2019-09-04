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
	"sort"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/AdvUni/gafaspot/database"
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

	reservations := database.GetUserReservations(username)
	// sort reservations
	sort.Slice(reservations, func(i, j int) bool {
		return reservations[i].Start.Before(reservations[j].Start)
	})
	var resNice []reservationNiceName
	for _, r := range reservations {
		resNice = append(resNice, newReservationNiceName(r))
	}
	err := personalviewTmpl.Execute(w, map[string]interface{}{"Username": username, "SSHkey": sshEntry, "Reservations": resNice})
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

	envPlainNames := database.GetUserActiveReservationEnv(username)
	sort.Strings(envPlainNames)

	var credsData []util.ReservationCreds

	for _, plainName := range envPlainNames {

		env, ok := environmentsMap[plainName]
		if !ok {
			env = util.Environment{NiceName: plainName, PlainName: plainName, HasSSH: false, Description: ""}
		}
		creds := vault.ReadCredentials(plainName)

		credsData = append(credsData, util.ReservationCreds{Res: util.Reservation{}, Env: env, Creds: creds})
	}

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

	err := reservationformTmpl.Execute(w, map[string]interface{}{"Username": username, "Envs": environments, "Selected": selectedEnvPlainName, "SSHmissing": sshMissing, "Error": errormessage})
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

	reservation.EnvPlainName = template.HTMLEscapeString(r.Form.Get("env"))
	if reservation.EnvPlainName == "" {
		logger.Debugf("reserve handler received reservation with invalid environment: %v", err)
		redirectInvalidSubmission(w, r, "environment invalid")
		return
	}

	startstring := template.HTMLEscapeString(r.Form.Get("startdate")) + " " + template.HTMLEscapeString(r.Form.Get("starttime"))
	reservation.Start, err = time.ParseInLocation(util.TimeLayout, startstring, time.Local)
	if err != nil {
		logger.Debugf("reserve handler received reservation with malformed date/time submission: %v", err)
		redirectInvalidSubmission(w, r, "start date/time malformed")
		return
	}

	endstring := template.HTMLEscapeString(r.Form.Get("enddate")) + " " + template.HTMLEscapeString(r.Form.Get("endtime"))
	reservation.End, err = time.ParseInLocation(util.TimeLayout, endstring, time.Local)
	if err != nil {
		logger.Debugf("reserve handler received reservation with malformed date/time submission: %v", err)
		redirectInvalidSubmission(w, r, "end date/time malformed")
		return
	}

	reservation.Subject = template.HTMLEscapeString(r.Form.Get("sub"))

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
		redirectInvalidSubmission(w, r, fmt.Sprintf("You did not enter a valid key (%v)", err))
		return
	}

	database.SaveUserSSH(username, sshPubkey)

	err = addkeysuccessTmpl.Execute(w, map[string]interface{}{"Username": username, "SSHkey": sshString})
	if err != nil {
		logger.Error(err)
	}
}

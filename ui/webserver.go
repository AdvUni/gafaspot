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
	"crypto/rand"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/AdvUni/gafaspot/database"
	"github.com/AdvUni/gafaspot/util"
	logging "github.com/alexcesaro/log"
	"github.com/gorilla/mux"
)

const (
	loginpage        = "/"
	login            = "/login"
	logout           = "/logout"
	mainview         = "/mainview"
	personalview     = "/personal"
	credsview        = "/personal/creds"
	reservationform  = "/newreservation/{env}"
	reserve          = "/reserve"
	abortreservation = "/abortreservation"
	addkeyform       = "/personal/addkey"
	uploadkey        = "/personal/uploadkey"
	deletekey        = "/personal/deletekey"
)

var (
	logger logging.Logger

	// This list contains all environment information from database table "environments".
	// This table shouldn't change during runtime, so the list content can be fetched once at program start.
	environments []util.Environment
	// This maps associates the PlainName attribute of each environment with its HasSSH attribute for a fast lookup.
	envHasSSHMap map[string]bool
	// This maps associates the PlainName attribute of each environment with its NiceName attribute for a fast lookup.
	envNiceNameMap map[string]string

	// The following are the parsed templates for all the application's web pages, ready for execution with the right parameters.
	loginformTmpl       *template.Template
	mainviewTmpl        *template.Template
	personalviewTmpl    *template.Template
	reservationformTmpl *template.Template
	reservesuccessTmpl  *template.Template
	credsviewTmpl       *template.Template
	addkeyformTmpl      *template.Template
	addkeysuccessTmpl   *template.Template
)

// all initialization which does not need parameters from main routine.
func init() {
	// generate a random key for signing json web tokens for authentication.
	// Save it to global var (file authentication)
	_, err := rand.Read(hmacKey)
	if err != nil {
		log.Fatalf("could not create key for jwt signing: %v\n", err)
	}

	// pre-assembling and caching of all the page templates
	const (
		topTmplFile             = "ui/templates/top.html"
		bottomTmplFile          = "ui/templates/bottom.html"
		navTmplFile             = "ui/templates/nav.html"
		wordbreakTmplFile       = "ui/templates/wordbreak.html"
		loginformTmplFile       = "ui/templates/login.html"
		mainviewTmplFile        = "ui/templates/mainview.html"
		personalviewTmplFile    = "ui/templates/personalview.html"
		reservationformTmplFile = "ui/templates/newreservation.html"
		reservesuccessTmplFile  = "ui/templates/reservesuccess.html"
		credsviewTmplFile       = "ui/templates/credsview.html"
		addkeyformTmplFile      = "ui/templates/addkey.html"
		addkeysuccessTmplFile   = "ui/templates/addkeysuccess.html"
	)
	loginformTmpl, err = template.ParseFiles(loginformTmplFile, topTmplFile, bottomTmplFile)
	if err != nil {
		log.Fatal(err)
	}
	mainviewTmpl, err = template.New(path.Base(mainviewTmplFile)).Funcs(template.FuncMap{
		"formatDatetime": func(t time.Time) string { return t.Format(util.TimeLayout) },
		"past":           func(r util.Reservation) bool { return r.End.Before(time.Now()) },
	}).ParseFiles(mainviewTmplFile, topTmplFile, bottomTmplFile, navTmplFile)
	if err != nil {
		log.Fatal(err)
	}
	personalviewTmpl, err = template.New(path.Base(personalviewTmplFile)).Funcs(template.FuncMap{
		"formatDatetime": func(t time.Time) string { return t.Format(util.TimeLayout) },
		"past":           func(r reservationNiceName) bool { return r.End.Before(time.Now()) },
	}).ParseFiles(personalviewTmplFile, topTmplFile, bottomTmplFile, navTmplFile, wordbreakTmplFile)
	if err != nil {
		log.Fatal(err)
	}
	reservationformTmpl, err = template.ParseFiles(reservationformTmplFile, topTmplFile, bottomTmplFile, navTmplFile)
	if err != nil {
		log.Fatal(err)
	}
	reservesuccessTmpl, err = template.New(path.Base(reservesuccessTmplFile)).Funcs(template.FuncMap{
		"formatDatetime": func(t time.Time) string { return t.Format(util.TimeLayout) },
	}).ParseFiles(reservesuccessTmplFile, topTmplFile, bottomTmplFile, navTmplFile)
	if err != nil {
		log.Fatal(err)
	}
	credsviewTmpl, err = template.ParseFiles(credsviewTmplFile, topTmplFile, bottomTmplFile, navTmplFile, wordbreakTmplFile)
	if err != nil {
		log.Fatal(err)
	}
	addkeyformTmpl, err = template.ParseFiles(addkeyformTmplFile, topTmplFile, bottomTmplFile, navTmplFile)
	if err != nil {
		log.Fatal(err)
	}
	addkeysuccessTmpl, err = template.ParseFiles(addkeysuccessTmplFile, topTmplFile, bottomTmplFile, navTmplFile, wordbreakTmplFile)
	if err != nil {
		log.Fatal(err)
	}
}

// RunWebserver registers all page handlers to a router and then starts the web server.
func RunWebserver(l logging.Logger, addr string) {
	logger = l

	// fetch static information about environments from database
	environments, envHasSSHMap, envNiceNameMap = database.GetEnvironments()

	// create router and register all paths
	router := mux.NewRouter()

	router.HandleFunc(loginpage, loginPageHandler)
	router.HandleFunc(login, loginHandler).Methods(http.MethodPost)
	router.HandleFunc(logout, logoutHandler).Methods(http.MethodPost)
	router.HandleFunc(mainview, mainPageHandler)
	router.HandleFunc(personalview, personalPageHandler)
	router.HandleFunc(credsview, credsPageHandler)
	router.HandleFunc(reservationform, newreservationPageHandler)
	router.HandleFunc(reserve, reserveHandler)
	router.HandleFunc(abortreservation, abortreservationHandler)
	router.HandleFunc(addkeyform, addkeyPageHandler)
	router.HandleFunc(uploadkey, uploadkeyHandler)
	router.HandleFunc(deletekey, deletekeyHandler)

	// start web server
	http.Handle(loginpage, router)
	err := http.ListenAndServe(addr, nil)
	// cause entire program to stop if the server crashes for any reason
	logger.Emergencyf("webserver crashed: %v\n", err)
	os.Exit(1)
}

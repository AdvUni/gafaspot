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
	"net/http"
	"time"
)

const (
	authCookieName  = "token"
	errorCookieName = "errormessage"
	infoCookieName  = "infomessage"
)

type reservationFormData struct {
	startdateStr string
	starttimeStr string
	enddateStr   string
	endtimeStr   string
	subject      string
}

func setAuthCookie(w http.ResponseWriter, token string, timeout time.Time) {
	cookie := &http.Cookie{
		Name:     authCookieName,
		Value:    token,
		Expires:  timeout,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}

func setMessageCookie(w http.ResponseWriter, cookieName, message string) {
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    message,
		MaxAge:   10,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}

func setErrorCookie(w http.ResponseWriter, message string) {
	setMessageCookie(w, errorCookieName, message)
}

func setInfoCookie(w http.ResponseWriter, message string) {
	setMessageCookie(w, infoCookieName, message)
}

func setReservationFormCookies(w http.ResponseWriter, data reservationFormData) {
	setMessageCookie(w, "startdate", data.startdateStr)
	setMessageCookie(w, "starttime", data.starttimeStr)
	setMessageCookie(w, "enddate", data.enddateStr)
	setMessageCookie(w, "endtime", data.endtimeStr)
	setMessageCookie(w, "subject", data.subject)
}

func readMessageCookie(w http.ResponseWriter, r *http.Request, cookieName string) string {
	var message string
	cookie, err := r.Cookie(cookieName)
	if err == nil {
		message = cookie.Value
		invalidateCookie(w, cookieName)
	}
	return message
}

func readErrorCookie(w http.ResponseWriter, r *http.Request) string {
	return readMessageCookie(w, r, errorCookieName)
}

func readInfoCookie(w http.ResponseWriter, r *http.Request) string {
	return readMessageCookie(w, r, infoCookieName)
}

func readReservationFormCookies(w http.ResponseWriter, r *http.Request) reservationFormData {
	var data reservationFormData
	data.startdateStr = readMessageCookie(w, r, "startdate")
	data.starttimeStr = readMessageCookie(w, r, "starttime")
	data.enddateStr = readMessageCookie(w, r, "enddate")
	data.endtimeStr = readMessageCookie(w, r, "endtime")
	data.subject = readMessageCookie(w, r, "subject")

	return data
}

func invalidateCookie(w http.ResponseWriter, cookieName string) {
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}

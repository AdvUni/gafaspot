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

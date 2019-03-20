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

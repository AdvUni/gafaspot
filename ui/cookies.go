package ui

import (
	"net/http"
	"time"
)

const (
	authCookieName  = "token"
	errorCookieName = "errormessage"
)

func setAuthCookie(w http.ResponseWriter, token string, timeout time.Time) {
	cookie := &http.Cookie{
		Name:     authCookieName,
		Value:    token,
		Expires:  timeout,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func setErrorCookie(w http.ResponseWriter, message string) {
	cookie := &http.Cookie{
		Name:     errorCookieName,
		Value:    message,
		MaxAge:   10,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
}

func readErrorCookie(w http.ResponseWriter, r *http.Request) string {
	var errormessage string
	cookie, err := r.Cookie(errorCookieName)
	if err == nil {
		errormessage = cookie.Value
		invalidateErrorCookie(w)
	}
	return errormessage
}

func invalidateCookie(w http.ResponseWriter, cookieName string) {
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func invalidateErrorCookie(w http.ResponseWriter) {
	invalidateCookie(w, errorCookieName)
}

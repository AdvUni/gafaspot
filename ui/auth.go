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

	"github.com/dgrijalva/jwt-go"
)

const (
	// Time, after which a user gets automatically logged out.
	authCookieTTL = 1 * time.Hour
)

// HMAC hash key for signing authentication cookies. Is randomly generated at web server start.
var hmacKey = make([]byte, 128)

type claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func verifyUser(w http.ResponseWriter, r *http.Request) (string, bool) {
	cookie, err := r.Cookie(authCookieName)
	if err != nil {
		logger.Debugf("authentication failed: %v\n", err)
		return "", false
	}

	tokenContent := &claims{}

	token, err := jwt.ParseWithClaims(cookie.Value, tokenContent, func(t *jwt.Token) (interface{}, error) { return hmacKey, nil })
	if err != nil {
		logger.Debug("authentication failed: %v\n", err)
		return "", false
	}
	if token.Valid {
		username := tokenContent.Username
		renewJWT(w, username)
		return username, true
	}
	logger.Debug("authentication failed: jwt is invalid")
	return "", false
}

func renewJWT(w http.ResponseWriter, username string) {
	timeout := time.Now().Add(authCookieTTL)
	jwtContent := &claims{username, jwt.StandardClaims{ExpiresAt: timeout.Unix()}}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, jwtContent).SignedString(hmacKey)

	if err != nil {
		logger.Error("creation of json web token not possible: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	setAuthCookie(w, token, timeout)
}

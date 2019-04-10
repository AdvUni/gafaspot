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

package vault

import (
	"fmt"
	"log"
	"strings"
)

const (
	createTokenPath   = "auth/approle/login"
	ldapAuthBasicPath = "auth/ldap/login"
)

var ldapAuthBasicURL, ldapAuthPolicy string
var apprl approle

type approle struct {
	roleID      string
	secretID    string
	getTokenURL string
}

func initApprole(approleID, approleSecret, vaultAddress string) {
	getTokenURL := joinRequestPath(vaultAddress, createTokenPath)
	apprl = approle{
		approleID,
		approleSecret,
		getTokenURL,
	}
}

func createVaultToken() string {
	payload := fmt.Sprintf("{\"role_id\": \"%v\", \"secret_id\": \"%v\"}", apprl.roleID, apprl.secretID)
	token, err := sendVaultTokenRequest(apprl.getTokenURL, strings.NewReader(payload))
	if err != nil {
		log.Println(err)
	}
	return token
}

func initLDAP(authPolicy, vaultAddress string) {
	ldapAuthBasicURL = joinRequestPath(vaultAddress, ldapAuthBasicPath)
	ldapAuthPolicy = authPolicy
}

func DoLdapAuthentication(username, password string) bool {
	url := ldapAuthBasicURL + "/" + username
	payload := strings.NewReader(fmt.Sprintf("{\"password\": \"%v\"}", password))

	availablePolicies, err := sendVaultLdapRequest(url, payload)
	if err == ErrAuth {
		return false
	} else if err != nil {
		log.Println(err)
		return false
	}

	for _, policy := range availablePolicies {
		if policy.(string) == ldapAuthPolicy {
			return true
		}
	}
	return false
}

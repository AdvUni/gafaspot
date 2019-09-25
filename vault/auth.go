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
	"strings"

	"github.com/AdvUni/gafaspot/util"
)

const (
	createEphemeralTokenPath = "auth/approle/login"
	createOrphanTokenPath    = "auth/token/create-orphan"
	ldapAuthBasicPath        = "auth/ldap/login"
)

var (
	ldapAuthBasicURL  string
	ldapAuthPolicy    string
	getOrphanTokenURL string
	apprl             approle
)

type approle struct {
	roleID      string
	secretID    string
	getTokenURL string
}

func initAuth(c util.GafaspotConfig) {
	// init approle
	getTokenURL := joinRequestPath(c.VaultAddress, createEphemeralTokenPath)
	apprl = approle{
		c.ApproleID,
		c.ApproleSecret,
		getTokenURL,
	}

	// init orphan token
	getOrphanTokenURL = joinRequestPath(c.VaultAddress, createOrphanTokenPath)
	tuneLeaseDuration(joinRequestPath(c.VaultAddress, "sys", "auth", "token", "tune"), c.MaxBookingDays)

	// init LDAP
	ldapAuthBasicURL = joinRequestPath(c.VaultAddress, ldapAuthBasicPath)
	ldapAuthPolicy = c.UserPolicy
}

// createOrphanVaultToken is for generating a long-living vault token. The
// function first performs an approle login and then uses the received
// ephemeral token to create an orphan token (orphan tokens do not get revoked
// as soon as their parents expire). The orphan token can be created with an
// individual life span, so they can be used to generate secrets leases at the
// start of a reservation.
func createOrphanVaultToken(ttl string) string {
	payload := fmt.Sprintf("{\"ttl\": \"%s\"}", ttl)
	token, err := sendVaultTokenRequest(getOrphanTokenURL, createEphemeralVaultToken(), strings.NewReader(payload))
	if err != nil {
		logger.Error(err)
	}
	return token
}

// createEphemeralVaultToken performs an approle login to vault and returns the
// received token. The token is only valid for a short time; this depends on
// the approle role configuration in vault. Do not use this tokens for
// generating secrets leases, as those leases would expire with the tokens.
func createEphemeralVaultToken() string {
	payload := fmt.Sprintf("{\"role_id\": \"%v\", \"secret_id\": \"%v\"}", apprl.roleID, apprl.secretID)
	token, err := sendVaultTokenRequest(apprl.getTokenURL, "", strings.NewReader(payload))
	if err != nil {
		logger.Error(err)
	}
	return token
}

// DoLdapAuthentication performs an LDAP authentication against a Vault LDAP Auth Method.
// It checks, whether username and password are accepted by the configured ldap server at all.
// If so, it checks whether Vault assigns the ldap-group-policy given in gafaspot_config.yaml
// to the login data. This is the case if the user is member of the correct LDAP group (and the
// vault auth method is configured correctly).
func DoLdapAuthentication(username, password string) bool {
	url := ldapAuthBasicURL + "/" + username
	payload := strings.NewReader(fmt.Sprintf("{\"password\": \"%v\"}", password))

	availablePolicies, err := sendVaultLdapRequest(url, payload)
	if err == ErrAuth {
		return false
	} else if err != nil {
		logger.Error(err)
		return false
	}

	for _, policy := range availablePolicies {
		if policy.(string) == ldapAuthPolicy {
			return true
		}
	}
	return false
}

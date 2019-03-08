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

func CreateVaultToken() string {
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

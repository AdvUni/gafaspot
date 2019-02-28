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

type Approle struct {
	roleID      string
	secretID    string
	getTokenURL string
}

func NewApprole(approleID, approleSecret, vaultAddress string) *Approle {
	getTokenURL := joinRequestPath(vaultAddress, createTokenPath)
	approle := Approle{
		approleID,
		approleSecret,
		getTokenURL,
	}
	return &approle
}

func (approle Approle) CreateVaultToken() string {
	payload := fmt.Sprintf("{\"role_id\": \"%v\", \"secret_id\": \"%v\"}", approle.roleID, approle.secretID)
	token, err := sendVaultTokenRequest(approle.getTokenURL, strings.NewReader(payload))
	if err != nil {
		log.Println(err)
	}
	return token
}

func InitLDAP(authPolicy, vaultAddress string) {
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

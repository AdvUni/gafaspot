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

type AuthLDAP struct {
	authBasicURL string
	authPolicy   string
}

func NewAuthLDAP(authPolicy, vaultAddress string) AuthLDAP {
	authBasicURL := joinRequestPath(vaultAddress, ldapAuthBasicPath)
	authLDAP := AuthLDAP{
		authBasicURL,
		authPolicy,
	}
	return authLDAP
}

func (ldap AuthLDAP) DoLdapAuthentication(username, password string) bool {
	url := ldap.authBasicURL + "/" + username
	payload := strings.NewReader(fmt.Sprintf("{\"password\": \"%v\"}", password))

	availablePolicies, err := sendVaultLdapRequest(url, payload)
	if err == ErrAuth {
		return false
	} else if err != nil {
		log.Println(err)
		return false
	}

	for _, policy := range availablePolicies {
		if policy.(string) == ldap.authPolicy {
			return true
		}
	}
	return false
}

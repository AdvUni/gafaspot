package vault

import (
	"fmt"
	"log"
	"strings"
)

const (
	createTokenPath = "auth/approle/login"
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
	token, err := sendVaultTokenRequest("POST", approle.getTokenURL, strings.NewReader(payload))
	if err != nil {
		log.Println(err)
	}
	return token
}

package vault

import (
	"fmt"
	"log"
	"strings"
)

const (
	sshSignPath = "sign"
)

type SshSecretEngine struct {
	signURL      string
	storeDataURL string
}

func NewSshSecretEngine(vaultAddress, operateBasicPath, storeBasicPath, middlePath, role string) SshSecretEngine {
	signURL := joinRequestPath(vaultAddress, operateBasicPath, middlePath, sshSignPath, role)
	log.Println("sign path: ", signURL)
	storeDataURL := joinRequestPath(vaultAddress, storeBasicPath, middlePath, role, "signature")
	log.Println("kv path: ", storeDataURL)

	return SshSecretEngine{
		signURL,
		storeDataURL,
	}
}

// TODO: implements methods

func (ssh SshSecretEngine) StartBooking(vaultToken, sshKey string) {
	data := fmt.Sprintf("{\"signature\": \"%v\"}", ssh.signKey(vaultToken, sshKey))
	log.Println(data)
	WriteSecret(vaultToken, ssh.storeDataURL, data)
}

func (ssh SshSecretEngine) EndBooking(vaultToken, sshKey string) {
	DeleteSecret(vaultToken, ssh.storeDataURL)
}

func (ssh SshSecretEngine) signKey(vaultToken, sshKey string) interface{} {

	data, err := sendVaultRequest("POST", ssh.signURL, vaultToken, strings.NewReader(sshKey))
	if err != nil {
		log.Println(err)
	}
	return data
}

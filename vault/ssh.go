package vault

import (
	"fmt"
	"log"
)

const (
	sshCredsPath = "creds"
)

type SshSecretEngine struct {
	VaultAddress string
	VaultPath    string
	Role         string
}

// TODO: implements methods

func (ssh SshSecretEngine) StartBooking(vaultToken, sshKey string) {
	fmt.Println(ssh.signKey(vaultToken, sshKey))
}

func (ssh SshSecretEngine) EndBooking(vaultToken, sshKey string) {
	fmt.Println("empty method")
}

func (ssh SshSecretEngine) signKey(vaultToken, sshKey string) interface{} {

	requestPath := joinRequestPath(ssh.VaultAddress, ssh.VaultPath, sshCredsPath, ssh.Role)

	log.Println("repuestPath: ", requestPath)

	data, err := sendVaultRequest("GET", requestPath, vaultToken, nil)
	if err != nil {
		log.Println(err)
	}
	return data
}

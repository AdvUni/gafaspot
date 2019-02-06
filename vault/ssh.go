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
	VaultAddress string
	VaultPath    string
	Role         string
}

// TODO: implements methods

func (ssh SshSecretEngine) StartBooking(vaultToken, sshKey string) {
	fmt.Println(ssh.signKey(vaultToken, sshKey))
	// TODO: Write results into kv secret engine
}

func (ssh SshSecretEngine) EndBooking(vaultToken, sshKey string) {
	// TODO: Delete contents from kv secret engine
	fmt.Println("empty method")
}

func (ssh SshSecretEngine) signKey(vaultToken, sshKey string) interface{} {

	requestPath := joinRequestPath(ssh.VaultAddress, ssh.VaultPath, sshSignPath, ssh.Role)

	log.Println("repuestPath: ", requestPath)

	data, err := sendVaultRequest("POST", requestPath, vaultToken, strings.NewReader(sshKey))
	if err != nil {
		log.Println(err)
	}
	return data
}

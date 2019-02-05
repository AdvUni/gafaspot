package vault

import (
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

func (ssh SshSecretEngine) ChangeCreds(vaultToken, sshKey string) interface{} {

	requestPath := joinRequestPath(ssh.VaultAddress, ssh.VaultPath, sshCredsPath, ssh.Role)

	log.Println("repuestPath: ", requestPath)

	data, err := sendVaultRequest("GET", requestPath, vaultToken, nil)
	if err != nil {
		log.Println(err)
	}
	return data
}

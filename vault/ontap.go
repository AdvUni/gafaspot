package vault

import (
	"log"
)

const (
	credsPath = "creds"
)

type OntapSecretEngine struct {
	VaultAddress string
	VaultPath    string
	Role         string
}

func (ontap OntapSecretEngine) ChangeCreds(vaultToken string) string {

	requestPath := joinRequestPath(ontap.VaultAddress, ontap.VaultPath, credsPath, ontap.Role)

	log.Println("repuestPath: ", requestPath)

	response := sendVaultRequest("GET", requestPath, vaultToken, nil)

	return response
}

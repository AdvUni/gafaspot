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

func (ontap OntapSecretEngine) ChangeCreds(vaultToken string) interface{} {

	requestPath := joinRequestPath(ontap.VaultAddress, ontap.VaultPath, credsPath, ontap.Role)

	log.Println("repuestPath: ", requestPath)

	data, err := sendVaultRequest("GET", requestPath, vaultToken, nil)
	if err != nil {
		log.Println(err)
	}
	return data
}

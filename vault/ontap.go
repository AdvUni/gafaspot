package vault

import (
	"fmt"
	"log"
)

const (
	ontapCredsPath = "creds"
)

type OntapSecretEngine struct {
	VaultAddress string
	VaultPath    string
	Role         string
	VaultKVPath  string
}

func (ontap OntapSecretEngine) StartBooking(vaultToken, _ string) {
	fmt.Println(ontap.changeCreds(vaultToken))
	// TODO: Write results into kv secret engine
}

func (ontap OntapSecretEngine) EndBooking(vaultToken, _ string) {
	// TODO: Delete contents from kv secret engine
	fmt.Println(ontap.changeCreds(vaultToken))
}

func (ontap OntapSecretEngine) changeCreds(vaultToken string) interface{} {

	requestPath := joinRequestPath(ontap.VaultAddress, ontap.VaultPath, ontapCredsPath, ontap.Role)

	log.Println("repuestPath: ", requestPath)

	data, err := sendVaultRequest("GET", requestPath, vaultToken, nil)
	if err != nil {
		log.Println(err)
	}
	return data
}

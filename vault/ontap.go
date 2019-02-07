package vault

import (
	"log"
	"fmt"
)

const (
	ontapCredsPath = "creds"
)

type OntapSecretEngine struct {
	changeCredsURL string
	storeDataURL   string
}

func NewOntapSecretEngine(vaultAddress, operateBasicPath, storeBasicPath, middlePath, role string) OntapSecretEngine {
	changeCredsURL := joinRequestPath(vaultAddress, operateBasicPath, middlePath, ontapCredsPath, role)
	log.Println("creds path: ", changeCredsURL)
	storeDataURL := joinRequestPath(vaultAddress, storeBasicPath, middlePath, role, "data")
	log.Println("kv path: ", storeDataURL)

	return OntapSecretEngine{
		changeCredsURL,
		storeDataURL,
	}
}

func (ontap OntapSecretEngine) StartBooking(vaultToken, _ string) {
	data := fmt.Sprintf("%v", ontap.changeCreds(vaultToken))
	log.Println(data)
	WriteSecret(vaultToken, ontap.storeDataURL, data)
}

func (ontap OntapSecretEngine) EndBooking(vaultToken, _ string) {
	DeleteSecret(vaultToken, ontap.storeDataURL)
	log.Println(ontap.changeCreds(vaultToken))
}

func (ontap OntapSecretEngine) changeCreds(vaultToken string) interface{} {

	data, err := sendVaultRequest("GET", ontap.changeCredsURL, vaultToken, nil)
	if err != nil {
		log.Println(err)
	}
	return data
}

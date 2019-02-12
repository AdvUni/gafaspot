package vault

import (
	"fmt"
	"log"
)

const (
	ontapCredsPath = "creds"
)

type OntapSecretEngine struct {
	changeCredsURL string
	storeDataURL   string
}

func (ontap OntapSecretEngine) StartBooking(vaultToken, _ string) {
	data := fmt.Sprintf("{\"data\": \"%v\"}", ontap.changeCreds(vaultToken))
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

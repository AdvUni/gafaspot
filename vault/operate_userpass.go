package vault

import (
	"fmt"
	"log"
)

const (
	userpassCredsPath = "creds"
)

type userpassSecEng struct {
	changeCredsURL string
	storeDataURL   string
}

func (secEng userpassSecEng) startBooking(vaultToken, _ string, _ int) {
	data := fmt.Sprintf("{\"data\": \"%v\"}", secEng.changeCreds(vaultToken))
	log.Println(data)
	vaultStorageWrite(vaultToken, secEng.storeDataURL, data)
}

func (secEng userpassSecEng) endBooking(vaultToken string) {
	vaultStorageDelete(vaultToken, secEng.storeDataURL)
	log.Println(secEng.changeCreds(vaultToken))
}

func (secEng userpassSecEng) changeCreds(vaultToken string) interface{} {

	data, err := sendVaultDataRequest("GET", secEng.changeCredsURL, vaultToken, nil)
	if err != nil {
		log.Println(err)
	}
	return data
}

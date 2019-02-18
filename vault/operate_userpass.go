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

func (secEng userpassSecEng) startBooking(vaultToken, _ string) {
	data := fmt.Sprintf("{\"data\": \"%v\"}", secEng.changeCreds(vaultToken))
	log.Println(data)
	vaultStorageWrite(vaultToken, secEng.storeDataURL, data)
}

func (secEng userpassSecEng) endBooking(vaultToken, _ string) {
	vaultStorageDelete(vaultToken, secEng.storeDataURL)
	log.Println(secEng.changeCreds(vaultToken))
}

func (SecEng userpassSecEng) changeCreds(vaultToken string) interface{} {

	data, err := sendVaultRequest("GET", SecEng.changeCredsURL, vaultToken, nil)
	if err != nil {
		log.Println(err)
	}
	return data
}

package vault

import (
	"fmt"
	"log"
)

const (
	userpassCredsPath = "creds"
)

// userpassSecEng is a SecEng implementation which works for all Vault secret engines listening to
// URL .../creds/rolename for changing credentials. This apllies to most of the credential secret engines
// such as ad, ontap and database.
type userpassSecEng struct {
	changeCredsURL string
	storeDataURL   string
}

// startBooking for a userpassSecEng means to change the credentials and store it inside the respective
// kv secret engine inside Vault.
func (secEng userpassSecEng) startBooking(vaultToken, _ string, _ int) {
	data := fmt.Sprintf("{\"data\": \"%v\"}", secEng.changeCreds(vaultToken))
	log.Println(data)
	vaultStorageWrite(vaultToken, secEng.storeDataURL, data)
}

// endBooking for a userpassSecEng means to delete the stored credentials from kv storage and then
// change the credentials again for them to become unknown.
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

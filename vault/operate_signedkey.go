package vault

import (
	"fmt"
	"log"
	"strings"
)

const (
	signPath = "sign"
)

type signedkeySecEng struct {
	signURL      string
	storeDataURL string
}

func (secEng signedkeySecEng) startBooking(vaultToken, sshKey string) {
	data := fmt.Sprintf("{\"signature\": \"%v\"}", secEng.signKey(vaultToken, sshKey))
	log.Println(data)
	vaultStorageWrite(vaultToken, secEng.storeDataURL, data)
}

func (secEng signedkeySecEng) endBooking(vaultToken, sshKey string) {
	vaultStorageDelete(vaultToken, secEng.storeDataURL)
}

func (secEng signedkeySecEng) signKey(vaultToken, sshKey string) interface{} {

	data, err := sendVaultRequest("POST", secEng.signURL, vaultToken, strings.NewReader(sshKey))
	if err != nil {
		log.Println(err)
	}
	return data
}

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

func (secEng signedkeySecEng) startBooking(vaultToken, sshKey string, ttl int) {

	data := fmt.Sprintf("{\"signature\": \"%v\"}", secEng.signKey(vaultToken, sshKey, ttl))
	log.Println(data)
	vaultStorageWrite(vaultToken, secEng.storeDataURL, data)
}

func (secEng signedkeySecEng) endBooking(vaultToken string) {
	vaultStorageDelete(vaultToken, secEng.storeDataURL)
}

func (secEng signedkeySecEng) signKey(vaultToken, sshKey string, ttl int) interface{} {

	payload := fmt.Sprintf("{\"public_key\": \"%v\", \"ttl\": \"%vs\"}", sshKey, ttl)

	data, err := sendVaultDataRequest("POST", secEng.signURL, vaultToken, strings.NewReader(payload))
	if err != nil {
		log.Println(err)
	}
	return data
}

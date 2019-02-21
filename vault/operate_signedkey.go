package vault

import (
	"fmt"
	"log"
	"strings"
)

const (
	signPath = "sign"
)

// signedkeySecEng is a SecEng implementation which for Vault's ssh secret engine used with signed certificates.
// In contrast to most of the other credential secret engines, signing certificates workes under URL
// .../sign/rolename. This returnes not a username and a password, but an ssh signature which can be used to
// log in into machines which are configured for this.
type signedkeySecEng struct {
	signURL      string
	storeDataURL string
}

// startBooking means for an ssh secret engine used with signed certificates to create an ssh signature for a given
// public key. The signature is valid for a specified duration. As it should expire exactly with the booking's
// expiration, the ttl in seconds is needed already at the booking's begin.
func (secEng signedkeySecEng) startBooking(vaultToken, sshKey string, ttl int) {

	data := fmt.Sprintf("{\"signature\": \"%v\"}", secEng.signKey(vaultToken, sshKey, ttl))
	log.Println(data)
	vaultStorageWrite(vaultToken, secEng.storeDataURL, data)
}

// endBooking only needs to delete the data from Vault's kv storage, as the signature expires at its own.
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

// Copyright 2019, Advanced UniByte GmbH.
// Author Marie Lohbeck.
//
// This file is part of Gafaspot.
//
// Gafaspot is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Gafaspot is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Gafaspot.  If not, see <https://www.gnu.org/licenses/>.

package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// signedkeySecEng is a SecEng implementation which for Vault's ssh secret engine used with signed certificates.
// In contrast to most of the other credential secret engines, signing certificates workes with endpoint
// .../sign/rolename. This returnes not a username and a password, but an ssh signature which can be used to
// log in into machines which are configured for this.
type signedkeySecEng struct {
	name         string
	signURL      string
	storeDataURL string
}

// startBooking means for an ssh secret engine used with signed certificates to create an ssh signature for a given
// public key. The signature is valid for a specified duration. As it should expire exactly with the booking's
// expiration, the ttl in seconds is needed already at the booking's begin.
func (secEng signedkeySecEng) startBooking(vaultToken, sshKey string, ttl int) {
	data, err := json.Marshal(secEng.signKey(vaultToken, sshKey, ttl))
	if err != nil {
		log.Println(err)
	}
	// remove the line feed from data, which is returned by the ssh secrets engine, as it corrupts the json
	bytes.Replace(data, []byte("\n"), nil, -1)

	vaultStorageWrite(vaultToken, secEng.storeDataURL, data)
}

func (secEng signedkeySecEng) getName() string {
	return secEng.name
}

// endBooking only needs to delete the data from Vault's kv storage, as the signature expires at its own.
func (secEng signedkeySecEng) endBooking(vaultToken string) {
	vaultStorageDelete(vaultToken, secEng.storeDataURL)
}

func (secEng signedkeySecEng) readCreds(vaultToken string) (map[string]interface{}, error) {
	return vaultStorageRead(vaultToken, secEng.storeDataURL)
}

func (secEng signedkeySecEng) signKey(vaultToken, sshKey string, ttl int) map[string]interface{} {

	payload := fmt.Sprintf("{\"public_key\": \"%v\", \"ttl\": \"%vs\"}", sshKey, ttl)

	data, err := sendVaultDataRequest("POST", secEng.signURL, vaultToken, strings.NewReader(payload))
	if err != nil {
		log.Println(err)
	}
	return data
}

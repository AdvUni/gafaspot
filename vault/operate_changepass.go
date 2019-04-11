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
	"encoding/json"
)

// changepassSecEng is a SecEng implementation which works for Vault secrets engines listening to
// endpoint .../creds/rolename for performing a password change. Currently, secrets engines of type
// ad and ontap belong to this implementation. This secrets engines don't work with leases and
// there is nothing like a lease duration or lease revocation to worry about.
type changepassSecEng struct {
	name           string
	changeCredsURL string
	storeDataURL   string
}

func (secEng changepassSecEng) getName() string {
	return secEng.name
}

// startBooking for a changepassSecEng means to change the credentials and store it inside the respective
// kv secret engine inside Vault.
func (secEng changepassSecEng) startBooking(vaultToken, _ string, _ int) {
	data, err := json.Marshal(secEng.changeCreds(vaultToken))
	if err != nil {
		logger.Errorf("not able to marshal new creds: %v", err)
	}
	vaultStorageWrite(vaultToken, secEng.storeDataURL, data)
}

// endBooking for a changepassSecEng means to delete the stored credentials from kv storage and then
// change the credentials again for them to become unknown.
func (secEng changepassSecEng) endBooking(vaultToken string) {
	vaultStorageDelete(vaultToken, secEng.storeDataURL)
	secEng.changeCreds(vaultToken)
}

func (secEng changepassSecEng) readCreds(vaultToken string) (map[string]interface{}, error) {
	return vaultStorageRead(vaultToken, secEng.storeDataURL)
}

func (secEng changepassSecEng) changeCreds(vaultToken string) map[string]interface{} {
	data, err := sendVaultDataRequest("GET", secEng.changeCredsURL, vaultToken, nil)
	if err != nil {
		logger.Errorf("not able to change Creds: %v", err)
	}
	return data
}

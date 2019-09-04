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
	"time"

	"github.com/AdvUni/gafaspot/util"
	logging "github.com/alexcesaro/log"
)

var (
	logger       logging.Logger
	environments map[string][]SecEng
)

// InitVaultParams initializes the vault package from gafaspot. Besides setting the logger, it
// reads several values from config and readies gafaspot to communicate with the vault
// Auth Methods. Further, it creates several SecEng objects to communicate with vault secrets
// engines.
func InitVaultParams(l logging.Logger, config util.GafaspotConfig) {

	logger = l

	initApprole(config.ApproleID, config.ApproleSecret, config.VaultAddress)
	initLDAP(config.UserPolicy, config.VaultAddress)
	environments = initSecEngs(config.Environments, config.VaultAddress, config.MaxBookingDays)
}

// StartBooking starts a booking for a whole environment. As the environment may include ssh secret
// engines, this function needs an ssh key. It also needs the time 'until' to determine the ttl for
// a possible ssh signature. If there is no ssh secret engine inside
// the environment, the ssKey parameter will be ignored everywhere.
func StartBooking(envPlainName, sshKey string, until time.Time) {
	vaultToken := createVaultToken()
	ttl := int(until.Sub(time.Now()).Seconds())
	environment, ok := environments[envPlainName]
	if !ok {
		logger.Errorf("tried to start booking for environment '%v' but it does not exist", envPlainName)
		return
	}
	for _, secEng := range environment {
		secEng.startBooking(vaultToken, sshKey, ttl)
	}
}

// EndBooking ends a booking for a whole environment.
func EndBooking(envPlainName string) {
	vaultToken := createVaultToken()
	environment, ok := environments[envPlainName]
	if !ok {
		logger.Errorf("tried to end booking for environment '%v' but it does not exist", envPlainName)
		return
	}
	for _, secEng := range environment {
		secEng.endBooking(vaultToken)
	}
}

// ReadCredentials reads the credentials from all KV Secrets Engine related to the environment
// envPlainName and returns them as map. Map keys are the Secrets Engine's names. If it is not
// possible to retrieve any credentials because the environment does not exist, an error message
// gets logged and the result is nill. If retrieving of credentials fails for a specific
// Secrets Engine, a small error message gets written into the map instead of the credentials, so
// that it will be automatically displayed in the creds view.
func ReadCredentials(envPlainName string) map[string]map[string]interface{} {
	environment, ok := environments[envPlainName]
	if !ok {
		logger.Warningf("tried to read creds for environment '%v' which does not exist", envPlainName)
		return nil
	}

	vaultToken := createVaultToken()

	credentials := make(map[string]map[string]interface{})
	for _, secEng := range environment {
		c, err := secEng.readCreds(vaultToken)
		if err != nil {
			logger.Warningf("failed to read creds from Secrets Engine '%v' in environment '%v': %v", secEng.getName(), envPlainName, err)
			credentials[secEng.getName()] = map[string]interface{}{"error": "not possible to provide credentials"}
		} else {
			credentials[secEng.getName()] = c
		}
	}
	return credentials
}

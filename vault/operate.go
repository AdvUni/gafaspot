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
	"fmt"

	"github.com/AdvUni/gafaspot/util"
)

const (
	wordOperate = "operate"
	wordStore   = "store"

	wordCreds = "creds"
	wordSign  = "sign"
)

// SecEng is a handler for one credential secrets engine such as "ad" or "ssh" inside Vault.
// As the secrets retrieved from a secrets engine needs to be saved somewhere, each credential secrets
// engine has an equivalently named kv secrets engine as storage which is also obtained by this interface.
// A SecEng stores the URLs to which the secrets engines listen to and provides the functionality which
// is needed to start and end bookings, as changing credentials and storing or deleting them.
type SecEng interface {
	getName() string
	startBooking(vaultToken, sshKey string, ttl int)
	endBooking(vaultToken string)
	readCreds(vaultToken string) (map[string]interface{}, error)
}

// NewSecEng creates a new SecEng. From string engineType, it decides, which implementation of the interface
// must be instanciated. The path snippets vaultAddress, env, name and role get assembled to the the
// URLs, to which the vault secrets engines listen to.
func newSecEng(engineType, vaultAddress, env, name, role string, maxBookingDays int) SecEng {
	switch engineType {
	case "ad", "ontap":
		secEng := changepassSecEng{}
		secEng.name = name
		secEng.changeCredsURL = joinRequestPath(vaultAddress, wordOperate, env, name, wordCreds, role)
		secEng.storeDataURL = joinRequestPath(vaultAddress, wordStore, env, name, role, "data")
		return secEng

	case "database":
		secEng := leaseSecEng{}
		secEng.name = name
		secEng.createLeaseURL = joinRequestPath(vaultAddress, wordOperate, env, name, wordCreds, role)
		secEng.revokeLeaseURL = joinRequestPath(vaultAddress, "sys", "leases", "revoke-prefix", wordOperate, env, name, wordCreds, role)
		secEng.storeDataURL = joinRequestPath(vaultAddress, wordStore, env, name, role, "data")

		tuneLeaseDurationURL := joinRequestPath(vaultAddress, "sys", "mounts", wordOperate, env, name, "tune")
		tuneLeaseDuration(tuneLeaseDurationURL, maxBookingDays)
		return secEng

	case "ssh":
		secEng := signedkeySecEng{}
		secEng.name = name
		secEng.signURL = joinRequestPath(vaultAddress, wordOperate, env, name, wordSign, role)
		secEng.storeDataURL = joinRequestPath(vaultAddress, wordStore, env, name, role, "signature")
		return secEng

	default:
		logger.Warning(fmt.Errorf("Unsupported Secrets Engine type: %v", engineType))
		return nil
	}
}

func initSecEngs(environmentConfigs map[string]util.EnvironmentConfig, vaultAddress string, maxBookingDays int) map[string][]SecEng {
	environments := make(map[string][]SecEng)
	for envPlainName, envConf := range environmentConfigs {
		envPlainName = util.CreatePlainIdentifier(envPlainName)
		var secretEngines []SecEng
		for _, engine := range envConf.SecretsEngines {
			logger.Debugf("name: %v, type: %v, role: %v\n", engine.NiceName, engine.EngineType, engine.Role)
			secretEngine := newSecEng(engine.EngineType, vaultAddress, envPlainName, engine.NiceName, engine.Role, maxBookingDays)
			logger.Debug(secretEngine)
			if secretEngine != nil {
				secretEngines = append(secretEngines, secretEngine)
			}
		}
		environments[envPlainName] = secretEngines
	}
	return environments
}

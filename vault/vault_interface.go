package vault

import (
	"fmt"
	"log"
	"time"

	"gitlab-vs.informatik.uni-ulm.de/gafaspot/util"
)

var environments map[string][]SecEng

func InitVaultParams(config util.GafaspotConfig) {
	initApprole(config.ApproleID, config.ApproleSecret, config.VaultAddress)
	initLDAP(config.UserPolicy, config.VaultAddress)
	environments = initSecEngs(config.Environments, config.VaultAddress)
}

// StartBooking starts a booking for a whole environment. As the environment may include ssh secret
// engines, this function needs an ssh key. It also needs the time 'until' to determine the ttl for
// a possible ssh signature. If there is no ssh secret engine inside
// the environment, the ssKey parameter will be ignored everywhere.
func StartBooking(envName, sshKey string, until time.Time) {
	vaultToken := CreateVaultToken()
	ttl := int(until.Sub(time.Now()).Seconds())
	environment, ok := environments[envName]
	if !ok {
		log.Fatalf("tried to start booking for environment '%v' but it does not exist", envName)
	}
	for _, secEng := range environment {
		secEng.startBooking(vaultToken, sshKey, ttl)
	}
}

// EndBooking ends a booking for a whole environment.
func EndBooking(envName string) {
	vaultToken := CreateVaultToken()
	environment, ok := environments[envName]
	if !ok {
		log.Fatalf("tried to end booking for environment '%v' but it does not exist", envName)
	}
	for _, secEng := range environment {
		secEng.endBooking(vaultToken)
	}
}

func ReadCredentials(envName string, vaultToken string) (map[string]interface{}, error) {
	environment, ok := environments[envName]
	if !ok {
		return nil, fmt.Errorf("environment '%v' does not exist", envName)
	}

	credentials := make(map[string]interface{})
	for _, secEng := range environment {
		c, err := secEng.readCreds(vaultToken)
		if err != nil {
			return nil, err
		}
		credentials[secEng.getName()] = c
	}
	return credentials, nil
}

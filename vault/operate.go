package vault

import (
	"fmt"
	"log"
	"time"
)

const (
	operateBasicPath = "operate"
	storeBasicPath   = "store"
)

// SecEng is a handler for one credential secret engine such as "ad" or "ssh" inside Vault.
// As the secrets retrieved from a secret engine needs to be saved somewhere, each credential secret
// engine has an equivalently named kv secret engine as storage which is also obtained by this interface.
// A SecEng stores the URLs to which the secret engines listen to and provides the functionality which
// is needed to start and end bookings, as changing credentials and storing or deleting them.
type SecEng interface {
	startBooking(vaultToken, sshKey string, ttl int)
	endBooking(vaultToken string)
}

// NewSecEng creates a new SecEng. From string engineType, it decides, which implementation of the interface
// must be instanciated. The path snippets vaultAddress, env, name and role get assembled to the the
// URLs, to which the vault secret engines listen to.
func NewSecEng(engineType, vaultAddress, env, name, role string) SecEng {
	switch engineType {
	case "ad", "ontap":
		log.Println("adding a creds secret engine")

		changeCredsURL := joinRequestPath(vaultAddress, operateBasicPath, env, name, userpassCredsPath, role)
		log.Println("creds path: ", changeCredsURL)
		storeDataURL := joinRequestPath(vaultAddress, storeBasicPath, env, name, role, "data")
		log.Println("kv path: ", storeDataURL)

		return userpassSecEng{
			changeCredsURL,
			storeDataURL,
		}
	case "ssh":
		log.Println("adding ssh secret engine")

		signURL := joinRequestPath(vaultAddress, operateBasicPath, env, name, signPath, role)
		log.Println("sign path: ", signURL)
		storeDataURL := joinRequestPath(vaultAddress, storeBasicPath, env, name, role, "signature")
		log.Println("kv path: ", storeDataURL)

		return signedkeySecEng{
			signURL,
			storeDataURL,
		}

	default:
		log.Println(fmt.Errorf("Unsupported Secret Engine type: %v", engineType))
		return nil
	}
}

// StartBooking starts a booking for a whole environment. As the environment may include ssh secret
// engines, this function needs an ssh key. It also needs the time string of format constants.timeLayout
// until to determine the ttl for a possible ssh signature. If there is no ssh secret engine inside
// the environment, the ssKey parameter will be ignored everywhere.
func StartBooking(environment []SecEng, vaultToken, sshKey string, until time.Time) {
	ttl := int(until.Sub(time.Now()).Seconds())
	for _, secEng := range environment {
		secEng.startBooking(vaultToken, sshKey, ttl)
	}
}

// EndBooking ends a booking for a whole environment.
func EndBooking(environment []SecEng, vaultToken string) {
	for _, secEng := range environment {
		secEng.endBooking(vaultToken)
	}
}

package vault

import (
	"fmt"
	"log"

	"gitlab-vs.informatik.uni-ulm.de/gafaspot/util"
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
	getName() string
	startBooking(vaultToken, sshKey string, ttl int)
	endBooking(vaultToken string)
	readCreds(vaultToken string) (interface{}, error)
}

// NewSecEng creates a new SecEng. From string engineType, it decides, which implementation of the interface
// must be instanciated. The path snippets vaultAddress, env, name and role get assembled to the the
// URLs, to which the vault secret engines listen to.
func newSecEng(engineType, vaultAddress, env, name, role string) SecEng {
	switch engineType {
	case "ad", "ontap":
		log.Println("adding a creds secret engine")

		changeCredsURL := joinRequestPath(vaultAddress, operateBasicPath, env, name, userpassCredsPath, role)
		log.Println("creds path: ", changeCredsURL)
		storeDataURL := joinRequestPath(vaultAddress, storeBasicPath, env, name, role, "data")
		log.Println("kv path: ", storeDataURL)

		return userpassSecEng{
			name,
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
			name,
			signURL,
			storeDataURL,
		}

	default:
		log.Println(fmt.Errorf("Unsupported Secret Engine type: %v", engineType))
		return nil
	}
}

func initSecEngs(environmentConfigs map[string]util.EnvironmentConfig, vaultAddress string) map[string][]SecEng {
	environments := make(map[string][]SecEng)
	for envName, envConf := range environmentConfigs {
		var secretEngines []SecEng
		for _, engine := range envConf.SecretEngines {
			fmt.Printf("name: %v, type: %v, role: %v\n", engine.Name, engine.EngineType, engine.Role)
			secretEngine := newSecEng(engine.EngineType, vaultAddress, envName, engine.Name, engine.Role)
			fmt.Println(secretEngine)
			if secretEngine != nil {
				secretEngines = append(secretEngines, secretEngine)
			}
		}
		environments[envName] = secretEngines
	}
	return environments
}

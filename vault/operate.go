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

// SecEng is a handler for one credential secrets engine such as "ad" or "ssh" inside Vault.
// As the secrets retrieved from a secrets engine needs to be saved somewhere, each credential secrets
// engine has an equivalently named kv secrets engine as storage which is also obtained by this interface.
// A SecEng stores the URLs to which the secrets engines listen to and provides the functionality which
// is needed to start and end bookings, as changing credentials and storing or deleting them.
type SecEng interface {
	getName() string
	startBooking(vaultToken, sshKey string, ttl int)
	endBooking(vaultToken string)
	readCreds(vaultToken string) (interface{}, error)
}

// NewSecEng creates a new SecEng. From string engineType, it decides, which implementation of the interface
// must be instanciated. The path snippets vaultAddress, env, name and role get assembled to the the
// URLs, to which the vault secrets engines listen to.
func newSecEng(engineType, vaultAddress, env, name, role string) SecEng {
	switch engineType {
	case "ad", "ontap":
		log.Println("adding a creds secrets engine")

		changeCredsURL := joinRequestPath(vaultAddress, operateBasicPath, env, name, changepassCredsPath, role)
		log.Println("creds path: ", changeCredsURL)
		storeDataURL := joinRequestPath(vaultAddress, storeBasicPath, env, name, role, "data")
		log.Println("kv path: ", storeDataURL)

		return changepassSecEng{
			name,
			changeCredsURL,
			storeDataURL,
		}
	case "ssh":
		log.Println("adding ssh secrets engine")

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
		log.Println(fmt.Errorf("Unsupported Secrets Engine type: %v", engineType))
		return nil
	}
}

func initSecEngs(environmentConfigs map[string]util.EnvironmentConfig, vaultAddress string) map[string][]SecEng {
	environments := make(map[string][]SecEng)
	for envPlainName, envConf := range environmentConfigs {
		envPlainName = util.CreatePlainIdentifier(envPlainName)
		var secretEngines []SecEng
		for _, engine := range envConf.SecretEngines {
			fmt.Printf("name: %v, type: %v, role: %v\n", engine.NiceName, engine.EngineType, engine.Role)
			secretEngine := newSecEng(engine.EngineType, vaultAddress, envPlainName, engine.NiceName, engine.Role)
			fmt.Println(secretEngine)
			if secretEngine != nil {
				secretEngines = append(secretEngines, secretEngine)
			}
		}
		environments[envPlainName] = secretEngines
	}
	return environments
}

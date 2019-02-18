package vault

import (
	"fmt"
	"log"
)

const (
	operateBasicPath = "operate"
	storeBasicPath   = "store"
)

type SecEng interface {
	startBooking(vaultToken, sshKey string)
	endBooking(vaultToken, sshKey string)
}

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

func StartBooking(environment []SecEng, vaultToken, sshKey string) {
	for _, secEng := range environment {
		secEng.startBooking(vaultToken, sshKey)
	}
}

func EndBooking(environment []SecEng, vaultToken, sshKey string) {
	for _, secEng := range environment {
		secEng.endBooking(vaultToken, sshKey)
	}
}

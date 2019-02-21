package vault

import (
	"fmt"
	"log"
	"time"

	"gitlab-vs.informatik.uni-ulm.de/gafaspot/constants"
)

const (
	operateBasicPath = "operate"
	storeBasicPath   = "store"
)

type SecEng interface {
	startBooking(vaultToken, sshKey string, ttl int)
	endBooking(vaultToken string)
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

func StartBooking(environment []SecEng, vaultToken, sshKey string, until string) {
	untilTime, err := time.ParseInLocation(constants.TimeLayout, until, time.Local)
	if err != nil {
		log.Fatal(err)
	}
	ttl := int(untilTime.Sub(time.Now()).Seconds())
	for _, secEng := range environment {
		secEng.startBooking(vaultToken, sshKey, ttl)
	}
}

func EndBooking(environment []SecEng, vaultToken string) {
	for _, secEng := range environment {
		secEng.endBooking(vaultToken)
	}
}

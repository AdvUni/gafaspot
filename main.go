package main

import (
	"gitlab-vs.informatik.uni-ulm.de/gafaspot/vault"
)

const (
	vaultAddress     = "http://127.0.0.1:8200/v1"
	operateBasicPath = "operate"
	storeBasicPath   = "store"
	sshKey           = ""
	vaultToken       = ""
)

func main() {

	testOntap := vault.NewOntapSecretEngine(vaultAddress, operateBasicPath, storeBasicPath, "ontap", "gafaspot")

	testSSH := vault.NewSshSecretEngine(vaultAddress, operateBasicPath, storeBasicPath, "ssh", "gafaspot")

	testOntap.StartBooking(vaultToken, "")
	testSSH.StartBooking(vaultToken, sshKey)
}

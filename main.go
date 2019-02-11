package main

import (
	//"gitlab-vs.informatik.uni-ulm.de/gafaspot/vault"
	"fmt"
)

const (
	sshKey     = ""
	vaultToken = ""
)

func main() {

	environments := initSecretEngines()
	fmt.Println(environments)

	//testOntap := vault.NewOntapSecretEngine(vaultAddress, operateBasicPath, storeBasicPath, "ontap", "gafaspot")

	//testSSH := vault.NewSshSecretEngine(vaultAddress, operateBasicPath, storeBasicPath, "ssh", "gafaspot")

	//testOntap.StartBooking(vaultToken, "")
	//testSSH.StartBooking(vaultToken, sshKey)
}

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
	//fmt.Println(environments)

	demo0 = environments["demo0"]
	fmt.Println(demo0)

	for _, secEng := range demo0 {
		secEng.StartBooking(vaultToken, sshKey)
	}

	//testOntap := vault.NewOntapSecretEngine(vaultAddress, operateBasicPath, storeBasicPath, "ontap", "gafaspot")

	//testSSH := vault.NewSshSecretEngine(vaultAddress, operateBasicPath, storeBasicPath, "ssh", "gafaspot")

	//testOntap.StartBooking(vaultToken, "")
	//testSSH.StartBooking(vaultToken, sshKey)
}

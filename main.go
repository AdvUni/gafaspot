package main

import (
	"fmt"
	"gitlab-vs.informatik.uni-ulm.de/gafaspot/vault"
)

const (
	vaultAddress = "http://127.0.0.1:8200/v1"
	vaultToken   = ""
)

func main() {
	testOntap := vault.OntapSecretEngine{
		vaultAddress,
		"ontap",
		"gafaspot",
	}

	fmt.Println(testOntap.ChangeCreds(vaultToken, ""))
}

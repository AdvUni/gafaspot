package main

import (
	"gitlab-vs.informatik.uni-ulm.de/gafaspot/vault"
)

const (
	vaultAddress = "http://127.0.0.1:8200/v1"
	sshKey       = ""
	vaultToken   = ""
)

func main() {
	testOntap := vault.OntapSecretEngine{
		vaultAddress,
		"ontap",
		"gafaspot",
		"",
	}

	testSSH := vault.SshSecretEngine{
		vaultAddress,
		"ontap",
		"gafaspot",
	}

	testOntap.StartBooking(vaultToken, "")
	testSSH.StartBooking(vaultToken, sshKey)
}

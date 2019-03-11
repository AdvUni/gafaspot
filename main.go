package main

import (
	"gitlab-vs.informatik.uni-ulm.de/gafaspot/database"
	"gitlab-vs.informatik.uni-ulm.de/gafaspot/ui"
	"gitlab-vs.informatik.uni-ulm.de/gafaspot/vault"
)

const (
	sshKey   = ""
	username = ""
	password = ""
)

func main() {

	// get all config
	config := readConfig()

	// do initialization with config values
	vault.InitVaultParams(config)
	database.InitDB(config)

	// start webserver and routine for processing reservations
	go handleReservationScanning()
	ui.RunWebserver(config.WebserviceAddress)
}

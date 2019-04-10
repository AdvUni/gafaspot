// Copyright 2019, Advanced UniByte GmbH.
// Author Marie Lohbeck.
//
// This file is part of Gafaspot.
//
// Gafaspot is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Gafaspot is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Gafaspot.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"gitlab-vs.informatik.uni-ulm.de/gafaspot/database"
	"gitlab-vs.informatik.uni-ulm.de/gafaspot/ui"
	"gitlab-vs.informatik.uni-ulm.de/gafaspot/vault"
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

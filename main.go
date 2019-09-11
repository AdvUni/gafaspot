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
	"github.com/alexcesaro/log/stdlog"

	"github.com/AdvUni/gafaspot/database"
	"github.com/AdvUni/gafaspot/email"
	"github.com/AdvUni/gafaspot/ui"
	"github.com/AdvUni/gafaspot/vault"
)

func main() {
	// init logger
	logger := stdlog.GetFromFlags()
	logger.Info("Welcome to Gafaspot!")

	// get all config
	logger.Info("Reading config...")
	config := readConfig(logger)

	// do initialization with config values
	logger.Info("Initialization...")
	vault.InitVaultParams(logger, config)
	database.InitDB(logger, config)
	email.InitMailing(logger, config)

	// start webserver and routine for processing reservations
	logger.Info("Starting reservation scanning routine...")
	go handleReservationScanning(logger, config.ScanningInterval)
	logger.Info("Starting web server...")
	ui.RunWebserver(logger, config.WebserviceAddress)
}

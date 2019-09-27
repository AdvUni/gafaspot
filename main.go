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
	"os"

	"github.com/alexcesaro/log/stdlog"

	"github.com/AdvUni/gafaspot/database"
	"github.com/AdvUni/gafaspot/email"
	"github.com/AdvUni/gafaspot/ui"
	"github.com/AdvUni/gafaspot/vault"
	"github.com/hashicorp/vault/sdk/helper/mlock"
)

func main() {
	// init logger
	logger := stdlog.GetFromFlags()
	logger.Debug("Warning: Log level DEBUG will print sensible information!")
	logger.Info("Welcome to Gafaspot!")

	// get all config
	logger.Info("Reading config...")
	config := readConfig(logger)

	// mlock
	if config.DisableMlock {
		logger.Debug("mlock is disabled by Gafaspot config")
	} else {
		// try to prevent Gafaspot's memory pages from swapping with mlock the same way as Vault does
		if !mlock.Supported() {
			logger.Emergency("Gafaspot uses mlock to prevent memory from being swapped to disk, but the mlock syscall is not supported by your system. Please disable Gafaspot from using it by setting the 'disable_mlock' option in Gafaspot's configuration file")
			os.Exit(1)
		} else {
			err := mlock.LockMemory()
			if err != nil {
				logger.Emergency("Gafaspot uses mlock to prevent memory from being swapped to disk, but the mlock syscall fails. Please enable mlock on your system - maybe you must run Gafaspot as root - or disable Gafaspot from using it by setting the 'disable_mlock' option in Gafaspot's configuration file")
				logger.Debugf("error with mlock: %v", err)
				os.Exit(1)
			}
		}
	}

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

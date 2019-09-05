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

	"github.com/AdvUni/gafaspot/util"
	logging "github.com/alexcesaro/log"
	"github.com/spf13/viper"
)

var (
	configDefaults = map[string]interface{}{
		"webservice-address":            "0.0.0.0:80",
		"gafaspot-mailaddress":          "gafaspot@gafaspot.com",
		"max-reservation-duration-days": 30,
		"max-queuing-time-months":       2,
		"db-path":                       "./gafaspot.db",
		"database-ttl-months":           12,
		"vault-address":                 "http://127.0.0.1:8200/v1",
		"ldap-group-policy":             "gafaspot-user-ldap",
	}
)

// readConfig unmarshals the config file into a GafaspotConfig struct.
func readConfig(logger logging.Logger) util.GafaspotConfig {

	// set config defaults
	for key, value := range configDefaults {
		viper.SetDefault(key, value)
	}

	viper.SetConfigName("gafaspot_config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		logger.Emergencyf("failed to read config file: %s", err)
		os.Exit(1)
	}

	config := util.GafaspotConfig{}
	err = viper.Unmarshal(&config)
	if err != nil {
		logger.Emergencyf("unable to decode config into GafaspotConfig struct: %v", err)
		os.Exit(1)
	}

	// check completeness
	if config.ApproleID == "" || config.ApproleSecret == "" {
		logger.Emergency("parameters approle-roleID and approle-secretID must be specified in config")
		os.Exit(1)
	}

	return config
}

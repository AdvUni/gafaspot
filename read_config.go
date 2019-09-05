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

// readConfig unmarshals the config file into a GafaspotConfig struct.
func readConfig(logger logging.Logger) util.GafaspotConfig {
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
	return config
}

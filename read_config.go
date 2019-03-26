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
	"log"

	"github.com/spf13/viper"
	"gitlab-vs.informatik.uni-ulm.de/gafaspot/util"
)

// readConfig unmarshals the config file into a GafaspotConfig struct.
func readConfig() util.GafaspotConfig {
	viper.SetConfigName("gafaspot_config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("can't read config file: %s", err)
	}

	config := util.GafaspotConfig{}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct: %v", err)
	}
	return config
}

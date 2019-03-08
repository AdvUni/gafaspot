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

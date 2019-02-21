package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/spf13/viper"
)

// GafaspotConfig is a struct to load every information from config file.
type GafaspotConfig struct {
	VaultAddress  string                       `mapstructure:"vault-address"`
	ApproleID     string                       `mapstructure:"approle-roleID"`
	ApproleSecret string                       `mapstructure:"approle-secretID"`
	UserPolicy    string                       `mapstructure:"ldap-group-policy"`
	Database      string                       `mapstructure:"db-path"`
	Environments  map[string]environmentConfig //`yaml:"environments"`
}

// environmentConfig is a struct to load information about one environment from config file.
type environmentConfig struct {
	SecretEngines []SecretEngineConfig //`yaml:"secretEngines"`
	Description   string               //`yaml:"description"`
}

// SecretEngineConfig is a struct to load information about one Secret Engine from config file.
type SecretEngineConfig struct {
	Name       string //`yaml:"name"`
	EngineType string `mapstructure:"type"`
	Role       string //`yaml:"role"`
}

// readConfig unmarshals the config file into a GafaspotConfig struct.
func readConfig() GafaspotConfig {
	viper.SetConfigName("gafaspot_config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("can't read config file: %s", err)
	}

	config := GafaspotConfig{}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct: %v", err)
	}
	return config
}

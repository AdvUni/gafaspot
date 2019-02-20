package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/spf13/viper"
)

type GafaspotConfig struct {
	VaultAddress  string                       `mapstructure:"vault-address"`
	ApproleID     string                       `mapstructure:"approle-roleID"`
	ApproleSecret string                       `mapstructure:"approle-secretID"`
	UserPolicy    string                       `mapstructure:"ldap-group-policy"`
	Database      string                       `mapstructure:"db-path"`
	Environments  map[string]environmentConfig //`yaml:"environments"`
}

type environmentConfig struct {
	SecretEngines []SecretEngineConfig //`yaml:"secretEngines"`
	Description   string               //`yaml:"description"`
}

type SecretEngineConfig struct {
	Name       string //`yaml:"name"`
	EngineType string `mapstructure:"type"`
	Role       string //`yaml:"role"`
}

func readConfig() GafaspotConfig {
	viper.SetConfigName("gafaspot_config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s \n", err))
	}

	config := GafaspotConfig{}
	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("unable to decode into struct, %v", err)
	}
	return config
}

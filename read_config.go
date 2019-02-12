package main

import (
	"fmt"
	"github.com/spf13/viper"
	"gitlab-vs.informatik.uni-ulm.de/gafaspot/vault"
	"log"
)

type GafaspotConfig struct {
	VaultAddress string                       `mapstructure:"vault-address"`
	Environments map[string]environmentConfig //`yaml:"environments"`
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
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	config := GafaspotConfig{}
	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("unable to decode into struct, %v", err)
	}
	return config
}

func initSecretEngines() map[string][]vault.SecretEngine {
	config := readConfig()

	log.Println(config.VaultAddress)

	environments := make(map[string][]vault.SecretEngine)
	for envName, envConf := range config.Environments {
		var secretEngines []vault.SecretEngine
		for _, engine := range envConf.SecretEngines {
			fmt.Printf("name: %v, type: %v, role: %v\n", engine.Name, engine.EngineType, engine.Role)
			secretEngine := vault.NewSecretEngine(engine.EngineType, config.VaultAddress, envName, engine.Name, engine.Role)
			fmt.Println(secretEngine)
			if secretEngine != nil {
				secretEngines = append(secretEngines, secretEngine)
			}
		}
		environments[envName] = secretEngines
	}

	return environments

}

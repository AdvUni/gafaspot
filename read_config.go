package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"gitlab-vs.informatik.uni-ulm.de/gafaspot/vault"
	"log"
)

type GafaspotConfig struct {
	VaultAddress string                       `mapstructure:"vault-address"`
	Database     string                       `mapstructure:"db-path"`
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

func initSecretEngines(config GafaspotConfig) map[string][]vault.SecretEngine {

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

func initDB(config GafaspotConfig) *sql.DB {
	log.Println(config.Database)
	db, err := sql.Open("sqlite3", config.Database)
	if err != nil {
		log.Fatal("Not able to open database: ", err)
	}

	// Create table reservations. If it already exists, don't overwrite
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS reservations (id INTEGER PRIMARY KEY, status TEXT NOT NULL, username TEXT NOT NULL, env_name TEXT NOT NULL, start DATETIME NOT NULL, end DATETIME NOT NULL, subject TEXT, labels TEXT, delete_on DATE NOT NULL);")
	if err != nil {
		log.Fatal(err)
	}

	// Create table users. If it already exists, don't overwrite
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (username TEXT UNIQUE NOT NULL, ssh_pub_key BLOB, delete_on DATE NOT NULL);")
	if err != nil {
		log.Fatal(err)
	}

	// Create table environments. If it already exist, deltete it first. Someone might have updated the environment configurations before system restart. So we want to create this table from scratch.
	_, err = db.Exec("DROP TABLE IF EXISTS environments;")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE environments (env_name TEXT UNIQUE NOT NULL, has_ssh BOOLEAN NOT NULL, description TEXT);")
	if err != nil {
		log.Fatal(err)
	}

	// Fill empty table environments with information from configuration file
	for envName, envConf := range config.Environments {
		envDescription := envConf.Description
		envHasSSH := false
		for _, secEng := range envConf.SecretEngines {
			if secEng.EngineType == "ssh" {
				envHasSSH = true
			}
		}
		_, err = db.Exec("INSERT INTO environments VALUES (?, ?, ?);", envName, envHasSSH, envDescription)
		if err != nil {
			log.Fatal(err)
		}
	}

	return db
}

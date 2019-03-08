package main

import (
	"database/sql"
	"log"

	"gitlab-vs.informatik.uni-ulm.de/gafaspot/util"

	_ "github.com/mattn/go-sqlite3"

	"gitlab-vs.informatik.uni-ulm.de/gafaspot/ui"
	"gitlab-vs.informatik.uni-ulm.de/gafaspot/vault"
)

const (
	sshKey   = ""
	username = ""
	password = ""
)

func main() {

	// get all config
	config := readConfig()

	// do initialization with config values
	vault.InitVaultParams(config)
	db := initDB(config)
	log.Printf("db: %v\n", db)

	// start webserver and routine for processing reservations
	go handleReservationScanning(db)
	ui.RunWebserver(db, config.WebserviceAddress)
}

func initDB(config util.GafaspotConfig) *sql.DB {
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

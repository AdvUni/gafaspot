package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gitlab-vs.informatik.uni-ulm.de/gafaspot/ui"
	"gitlab-vs.informatik.uni-ulm.de/gafaspot/vault"
)

//"gitlab-vs.informatik.uni-ulm.de/gafaspot/vault"

const (
	sshKey   = ""
	username = ""
	password = ""
)

func main() {

	config := readConfig()

	environments := initSecEngs(config)
	log.Printf("environments: %v\n", environments)
	db := initDB(config)
	log.Printf("db: %v\n", db)
	approle := initApprole(config)
	ldap := initLdapAuth(config)

	log.Println("Auhtentication test:")
	log.Printf("Should return false: %v", ldap.DoLdapAuthentication("wrongUsername", "wrongPassword"))
	log.Printf("Should return true: %v", ldap.DoLdapAuthentication(username, password))

	ui.CreateReservation(db, "firstuser", "demo0", "testsubject", "", time.Date(2019, time.February, 25, 9, 0, 0, 0, time.Local), time.Date(2019, time.February, 25, 10, 0, 0, 0, time.Local))
	ui.CreateReservation(db, "seconduser", "demo0", "testsubject", "", time.Date(2019, time.February, 25, 10, 1, 0, 0, time.Local), time.Date(2019, time.February, 26, 10, 15, 0, 0, time.Local))
	ui.CreateReservation(db, "thirduser", "demo0", "testsubject", "", time.Date(2019, time.February, 27, 9, 0, 0, 0, time.Local), time.Date(2019, time.February, 28, 10, 0, 0, 0, time.Local))

	handleBookings(db, environments, approle)
}

func initSecEngs(config GafaspotConfig) map[string][]vault.SecEng {

	log.Println(config.VaultAddress)

	environments := make(map[string][]vault.SecEng)
	for envName, envConf := range config.Environments {
		var secretEngines []vault.SecEng
		for _, engine := range envConf.SecretEngines {
			fmt.Printf("name: %v, type: %v, role: %v\n", engine.Name, engine.EngineType, engine.Role)
			secretEngine := vault.NewSecEng(engine.EngineType, config.VaultAddress, envName, engine.Name, engine.Role)
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

func initApprole(config GafaspotConfig) *vault.Approle {
	return vault.NewApprole(config.ApproleID, config.ApproleSecret, config.VaultAddress)
}

func initLdapAuth(config GafaspotConfig) vault.AuthLDAP {
	return vault.NewAuthLDAP(config.UserPolicy, config.VaultAddress)
}

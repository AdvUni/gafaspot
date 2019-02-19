package main

import (
	"database/sql"
	"fmt"
	"log"

	"gitlab-vs.informatik.uni-ulm.de/gafaspot/vault"
)

//"gitlab-vs.informatik.uni-ulm.de/gafaspot/vault"

const (
	sshKey     = ""
	vaultToken = ""
)

func main() {

	config := readConfig()

	environments := initSecEngs(config)
	log.Printf("environments: %v\n", environments)
	db := initDB(config)
	log.Printf("db: %v\n", db)
	approle := initApprole(config)

	handleBookings(db, environments, approle)

	stmt, err := db.Prepare("INSERT INTO reservations (status, username, env_name, start, end, delete_on) VALUES(?,?,?,?,?,?);")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec("upcoming", "some_user", "demo0", "2019-02-14 22:00:00", "2019-02-22 00:00:00", "2020-02-15 00:00:00")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)

	/* 	res, err = stmt.Exec("upcoming", "some_user", "demo1", "2019-02-14 00:00:00", "2019-02-15 00:00:00", "2020-02-15 00:00:00")
	   	if err != nil {
	   		log.Fatal(err)
	   	}
	   	log.Println(res) */

	res, err = stmt.Exec("upcoming", "other_user", "demo0", "2019-02-12 00:00:00", "2019-02-14 10:00:00", "2020-02-15 00:00:00")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)

	res, err = stmt.Exec("upcoming", "third_user", "demo0", "2019-02-25 00:00:00", "2019-02-26 10:00:00", "2020-02-15 00:00:00")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)

	//demo0 := environments["demo0"]
	//fmt.Println(demo0)

	//for _, secEng := range demo0 {
	//	secEng.StartBooking(vaultToken, sshKey)
	//}

	//testOntap := vault.NewOntapSecretEngine(vaultAddress, operateBasicPath, storeBasicPath, "ontap", "gafaspot")

	//testSSH := vault.NewSshSecretEngine(vaultAddress, operateBasicPath, storeBasicPath, "ssh", "gafaspot")

	//testOntap.StartBooking(vaultToken, "")
	//testSSH.StartBooking(vaultToken, sshKey)
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

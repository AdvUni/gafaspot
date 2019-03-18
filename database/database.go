package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"gitlab-vs.informatik.uni-ulm.de/gafaspot/util"
)

const (
	// General TTL for old database entries in years. Applies to tables users and reservations.
	yearsTTL = 2
)

var db *sql.DB

// InitDB prepares the database for gafaspot. Opens the database at the path given in config file.
// As SQLITE is used, database doesn't even need to exist yet. Prepares all database tables and
// fills the environments table with the information from config file.
// Sets the package variable db to enable every function in the package to access the database.
func InitDB(config util.GafaspotConfig) {
	log.Println(config.Database)
	var err error

	// Open database. SQLITE databases are simple files, and if database doesn't exist yet, a new file will be createt at the specified path
	db, err = sql.Open("sqlite3", config.Database)
	if err != nil {
		log.Fatal("Not able to open database: ", err)
	}

	// Create table reservations. If it already exists, don't overwrite
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS reservations (id INTEGER PRIMARY KEY, status TEXT NOT NULL, username TEXT NOT NULL, env_plain_name TEXT NOT NULL, start DATETIME NOT NULL, end DATETIME NOT NULL, subject TEXT, labels TEXT, delete_on DATE NOT NULL);")
	if err != nil {
		log.Fatal(err)
	}

	// Create table users. If it already exists, don't overwrite
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (username TEXT UNIQUE NOT NULL, ssh_pub_key BLOB, delete_on DATE NOT NULL);")
	if err != nil {
		log.Fatal(err)
	}

	// Create table environments. If it already exist, deltete it first. Someone might have updated the environment configurations before system restart. So this table should be created from scratch.
	_, err = db.Exec("DROP TABLE IF EXISTS environments;")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE environments (env_plain_name TEXT UNIQUE NOT NULL, env_nice_name TEXT NOT NULL, has_ssh BOOLEAN NOT NULL, description TEXT);")
	if err != nil {
		log.Fatal(err)
	}

	// Fill empty table environments with information from configuration file
	for envPlainName, envConf := range config.Environments {
		envPlainName = util.CreatePlainIdentifier(envPlainName)
		envNiceName := envConf.NiceName
		if envNiceName == "" {
			envNiceName = envPlainName
		}
		envDescription := envConf.Description
		envHasSSH := false
		for _, secEng := range envConf.SecretEngines {
			if secEng.EngineType == "ssh" {
				envHasSSH = true
			}
		}
		_, err = db.Exec("INSERT INTO environments VALUES (?, ?, ?, ?);", envPlainName, envNiceName, envHasSSH, envDescription)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func beginTransaction() *sql.Tx {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	return tx
}

func commitTransaction(tx *sql.Tx) {
	err := tx.Commit()
	if err != nil {
		log.Println(err)
	}
}

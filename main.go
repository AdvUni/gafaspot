package main

import (
	//"gitlab-vs.informatik.uni-ulm.de/gafaspot/vault"
	"fmt"
)

const (
	sshKey     = ""
	vaultToken = ""
)

func main() {

	config := readConfig()
	environments := initSecretEngines(config)
	db := initDB(config)

	stmt, err := db.Prepare("INSERT INTO reservations (upcoming, username, env_name, start, end, deletion_date) VALUES(?,?,?,?,?);")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec("some_user", "demo0", "2019-02-14 00:00:00", "2019-02-15 00:00:00", "2020-02-15 00:00:00")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec("some_user", "demo1", "2019-02-14 00:00:00", "2019-02-15 00:00:00", "2020-02-15 00:00:00")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec("other_user", "demo0", "2019-02-12 00:00:00", "2019-02-14 10:00:00", "2020-02-15 00:00:00")
	if err != nil {
		log.Fatal(err)
	}

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

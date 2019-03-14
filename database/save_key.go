package database

import (
	"log"
	"time"
)

func SaveUserSSH(username string, ssh []byte) {
	deleteOn := time.Now().AddDate(yearsTTL, 0, 0)
	stmt, err := db.Prepare("REPLACE INTO users (username, ssh_pub_key, delete_on) VALUES(?,?,?);")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, ssh, deleteOn)
	if err != nil {
		log.Println(err)
	}
}

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

func RefreshDeletionDate(username string) {
	deleteOn := time.Now().AddDate(yearsTTL, 0, 0)
	stmt, err := db.Prepare("UPDATE users SET delete_on=? WHERE username=?;")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(deleteOn, username)
	if err != nil {
		log.Println(err)
	}
}

func DeleteUser(username string) {
	stmt, err := db.Prepare("DELETE FROM users WHERE username=?;")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(username)
	if err != nil {
		log.Println(err)
	}
}

func DeleteOldUserEntries(now time.Time) {
	stmt, err := db.Prepare("DELETE FROM users WHERE delete_on<=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(now)
	if err != nil {
		log.Println(err)
	}
}

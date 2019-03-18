package database

import (
	"log"
	"time"
)

// SaveUserSSH takes an ssh key and stores it with the username to database table users.
// Function does not perform any checks, so make sure you validate the key format earlier.
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

// RefreshDeletionDate updates the column "delete_on" in the users table for a specific user, if
// a user entry exists. The delete_on entry is for enable the program to delete old entries for
// users, which haven't logged in for a long time. So, each time a user actually logs in, the
// deletion date has to be refreshed.
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

// DeleteUser deltes a database entry in table users for a specific username.
// Note: There must not be any database entry for a user to use gafaspot. The users table is
// recently only nessecary for associating ssh keys with users. So this function is mainly for
// deleting a user's ssh key.
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

// DeleteOldUserEntries deletes all users from database table "users", who haven't logged in for a
// long time ("long time" is defined by constant "yearsTTL"). Old user entries are recognized by
// their delete_on column. So, this function deletes all user entries, whoseq delete_on dates are
// exeeded.
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

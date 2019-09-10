// Copyright 2019, Advanced UniByte GmbH.
// Author Marie Lohbeck.
//
// This file is part of Gafaspot.
//
// Gafaspot is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Gafaspot is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Gafaspot.  If not, see <https://www.gnu.org/licenses/>.

package database

import (
	"bytes"
	"fmt"
	"os"
	"time"
)

// SaveUserSSH takes an ssh key and stores it with the username to database table users.
// Function does not perform any checks, so make sure you validate the key format earlier.
func SaveUserSSH(username string, ssh []byte) {

	// remove line breaks from ssh key
	ssh = bytes.Replace(ssh, []byte("\n"), nil, -1)
	ssh = bytes.Replace(ssh, []byte("\r"), nil, -1)

	saveUserAttribute(username, "ssh_pub_key", ssh)
}

// SaveUserEmail takes an mail address and stores it with the username to database table users.
// Function does not perform any checks, so make sure you validate the address format earlier.
func SaveUserEmail(username, email string) {
	saveUserAttribute(username, "email", email)
}

func saveUserAttribute(username, attribute string, value interface{}) {
	deleteOn := addTTL(time.Now())

	tx := beginTransaction()
	defer commitTransaction(tx)

	// if user already exists, he gets updated, whereupon other attributes get not overwritten
	// if user does not exist, this statement does nothing
	stmtstring := fmt.Sprintf("UPDATE users SET %s=?, delete_on=? WHERE username=?;", attribute)
	stmt, err := tx.Prepare(stmtstring)
	if err != nil {
		logger.Emergency(err)
		os.Exit(1)
	}
	defer stmt.Close()
	_, err = stmt.Exec(value, deleteOn, username)
	if err != nil {
		logger.Error(err)
	}

	// if user didn't exist, this statement creates a new user
	// if user already exist, this statement does nothing
	stmtstring = fmt.Sprintf("INSERT OR IGNORE INTO users (username, %s, delete_on) VALUES(?,?,?);", attribute)
	stmt, err = tx.Prepare(stmtstring)
	if err != nil {
		logger.Emergency(err)
		os.Exit(1)
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, value, deleteOn)
	if err != nil {
		logger.Error(err)
	}
}

// RefreshDeletionDate updates the column "delete_on" in the users table for a specific user, if
// a user entry exists. The delete_on entry is for enable the program to delete old entries for
// users, which haven't logged in for a long time. So, each time a user actually logs in, the
// deletion date has to be refreshed.
func RefreshDeletionDate(username string) {
	deleteOn := addTTL(time.Now())
	stmt, err := db.Prepare("UPDATE users SET delete_on=? WHERE username=?;")
	if err != nil {
		logger.Emergency(err)
		os.Exit(1)
	}
	defer stmt.Close()

	_, err = stmt.Exec(deleteOn, username)
	if err != nil {
		logger.Error(err)
	}
}

// deleteUser deletes a database entry in table users for a specific username.
// Note: There must not be any database entry for a user to use gafaspot. The users table is
// only necessary for associating ssh keys and email addresses with users.
func deleteUser(username string) {
	stmt, err := db.Prepare("DELETE FROM users WHERE username=?;")
	if err != nil {
		logger.Emergency(err)
		os.Exit(1)
	}
	defer stmt.Close()

	_, err = stmt.Exec(username)
	if err != nil {
		logger.Error(err)
	}
}

// DeleteUserSSH deletes a user's ssh public key from database. The rest of the user's
// database entry remains.
func DeleteUserSSH(username string) {
	if UserHasEmail(username) {
		SaveUserSSH(username, []byte(""))
	} else {
		deleteUser(username)
	}
}

// DeleteUserEmail deletes a user's e-mail address from database. The rest of the user's
// database entry remains.
func DeleteUserEmail(username string) {
	if UserHasSSH(username) {
		SaveUserEmail(username, "")
	} else {
		deleteUser(username)
	}
}

// DeleteOldUserEntries deletes all users from database table "users", who haven't logged in for a
// long time ("long time" is defined by constant "yearsTTL"). Old user entries are recognized by
// their delete_on column. So, this function deletes all user entries, whose delete_on dates are
// exceeded.
func DeleteOldUserEntries(now time.Time) {
	stmt, err := db.Prepare("DELETE FROM users WHERE delete_on<=?")
	if err != nil {
		logger.Emergency(err)
		os.Exit(1)
	}
	defer stmt.Close()

	_, err = stmt.Exec(now)
	if err != nil {
		logger.Error(err)
	}
}

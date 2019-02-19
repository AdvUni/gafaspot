package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gitlab-vs.informatik.uni-ulm.de/gafaspot/vault"
)

func handleBookings(db *sql.DB, environments map[string][]vault.SecEng, approle *vault.Approle) {

	stmt, err := db.Prepare("SELECT ? FROM ? WHERE ?;")
	s, err := db.Prepare("SELECT (id, username, env_name) FROM reservations WHERE (status='?') AND (start<='" + time.Now().String() + "');")
	update, err := db.Prepare("UPDATE reservations SET status='?' WHERE id=?;")

	// TODO: endless loop, triggered each 5 minutes

	// have to start any upcoming bookings?
	resRows, err := s.Query("upcoming")
	if err != nil {
		log.Println(err)
	}
	defer resRows.Close()
	for resRows.Next() {
		var reservationID int
		var username, envName string
		err := resRows.Scan(&reservationID, &username, &envName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(reservationID, username, envName)

		// Does environment exist? Is ssh key needed?
		var hasSSH bool
		err = stmt.QueryRow("has_ssh", "environments", "(env_name='"+envName+"')").Scan(&hasSSH)
		if err == sql.ErrNoRows {
			log.Fatalf("Environment %v does not exist. Can't book it.", envName)
		} else if err != nil {
			log.Fatal(err)
		}

		sshKey := ""
		if hasSSH {
			// retrieve ssh key from user table
			err := stmt.QueryRow("ssh_pub_key", "users", "(username='"+username+"')").Scan(&sshKey)
			if err != nil {
				log.Fatal(err)
			}
		}

		// trigger the start of the booking
		vaultToken := approle.CreateVaultToken()
		vault.StartBooking(environments[envName], vaultToken, sshKey)

		// change booking status in database
		_, err = update.Exec("active", reservationID)
	}

	// TODO: have to end any active bookings?
	resRows, err = s.Query("active")
	if err != nil {
		log.Println(err)
	}
	defer resRows.Close()
	for resRows.Next() {
	}

	// TODO: have to delete any expired bookings?
}

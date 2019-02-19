package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gitlab-vs.informatik.uni-ulm.de/gafaspot/vault"
)

func readValues(db *sql.DB, resRows *sql.Rows, lookupSSH bool) (int, string, string) {
	// TODO: error handling
	var reservationID int
	var username, envName string
	err := resRows.Scan(&reservationID, &username, &envName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Values from matching query: id - %v, username - %v, envname %v", reservationID, username, envName)

	// Does environment exist? Is ssh key needed?
	var hasSSH bool
	err = db.QueryRow("SELECT has_ssh FROM environments WHERE (env_name='" + envName + "');").Scan(&hasSSH)
	if err == sql.ErrNoRows {
		log.Fatalf("Environment %v does not exist. Can't book it.", envName)
	} else if err != nil {
		log.Fatal(err)
	}

	sshKey := ""
	// only search for ssh key, if lookupSSH is true. Otherwise, return parameter sshKey is always empty.
	if lookupSSH && hasSSH {
		// retrieve ssh key from user table
		err := db.QueryRow("SELECT ssh_pub_key FROM users WHERE (username='" + username + "');").Scan(&sshKey)
		if err != nil {
			log.Fatal(err)
		}
	}

	return reservationID, envName, sshKey
}

func handleBookings(db *sql.DB, environments map[string][]vault.SecEng, approle *vault.Approle) {

	log.Println("startet booking handle procedure")

	selectCurrentEvents, err := db.Prepare("SELECT (id, username, env_name) FROM reservations WHERE (status='?') AND (?<='" + time.Now().String() + "');")
	updateState, err := db.Prepare("UPDATE reservations SET status='?' WHERE id=?;")

	// TODO: endless loop, triggered each 5 minutes

	// have to start any upcoming bookings?
	resRows, err := selectCurrentEvents.Query("upcoming", "start")
	if err != nil {
		log.Println(err)
	}
	defer resRows.Close()
	for resRows.Next() {
		reservationID, envName, sshKey := readValues(db, resRows, true)

		// trigger the start of the booking
		vaultToken := approle.CreateVaultToken()
		vault.StartBooking(environments[envName], vaultToken, sshKey)

		// change booking status in database
		_, err = updateState.Exec("active", reservationID)
	}

	// have to end any active bookings?
	resRows, err = selectCurrentEvents.Query("active", "end")
	if err != nil {
		log.Println(err)
	}
	defer resRows.Close()
	for resRows.Next() {
		reservationID, envName, _ := readValues(db, resRows, false)

		// trigger the end of the booking
		vaultToken := approle.CreateVaultToken()
		vault.EndBooking(environments[envName], vaultToken)

		// change booking status in database
		_, err = updateState.Exec("expired", reservationID)
	}

	// have to delete any expired bookings?
	resRows, err = selectCurrentEvents.Query("expired", "delete_on")
	if err != nil {
		log.Println(err)
	}
	defer resRows.Close()
	for resRows.Next() {
		reservationID, _, _ := readValues(db, resRows, false)

		// delete booking from database
		_, err = db.Exec("DELETE FROM reservations WHERE id=?;", reservationID)
	}
	log.Println("end booking handle procedure")
}

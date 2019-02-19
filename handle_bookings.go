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
	fmt.Printf("Values from matching query: id - %v, username - %v, envname %v\n", reservationID, username, envName)

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

	// use sprintf formatting here instead of prepared statement, because prepared statements seems not to cope with coulumn name insertion
	// this should be save because non of the parameters is user input
	selectCurrentEvents := "SELECT id, username, env_name FROM reservations WHERE (status='%v') AND (%v<='" + time.Now().String() + "');"
	updateState, err := db.Prepare("UPDATE reservations SET status=? WHERE id=?;")
	defer updateState.Close()
	if err != nil {
		log.Println(err)
	}

	// TODO: endless loop, triggered each 5 minutes

	// have to start any upcoming bookings?
	resRows, err := db.Query(fmt.Sprintf(selectCurrentEvents, "upcoming", "start"))
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
		log.Println("executed start of booking")
		_, err = updateState.Exec("active", reservationID)
		if err != nil {
			log.Printf("did not change status from upcoming to active due to following error: %v\n", err)
		}
	}

	// have to end any active bookings?
	resRows, err = db.Query(fmt.Sprintf(selectCurrentEvents, "active", "end"))
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
		log.Println("executed end of booking")
		_, err = updateState.Exec("expired", reservationID)
		if err != nil {
			log.Printf("did not change status from active to expired due to following error: %v\n", err)
		}
	}

	// have to delete any expired bookings?
	resRows, err = db.Query(fmt.Sprintf(selectCurrentEvents, "expired", "delete_on"))
	if err != nil {
		log.Println(err)
	}
	defer resRows.Close()
	for resRows.Next() {
		reservationID, _, _ := readValues(db, resRows, false)

		// delete booking from database
		_, err = db.Exec("DELETE FROM reservations WHERE id=?;", reservationID)
		if err != nil {
			log.Printf("did not delete database entry due to following error: %v\n", err)
		}
	}
	log.Println("end booking handle procedure")
}

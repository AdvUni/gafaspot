package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gitlab-vs.informatik.uni-ulm.de/gafaspot/vault"
)

type reservationRow struct {
	id       int
	username string
	envName  string
	hasSSH   bool
}

func (row reservationRow) check(db *sql.DB) {
	// Does environment exist? Does it has components which require an ssh key for login?
	err := db.QueryRow("SELECT has_ssh FROM environments WHERE (env_name='" + row.envName + "');").Scan(&row.hasSSH)
	if err == sql.ErrNoRows {
		log.Fatalf("environment %v does not exist", row.envName)
	} else if err != nil {
		log.Fatal(err)
	}
}

func getRows(db *sql.DB, status, timeCol string) []reservationRow {
	// use sprintf formatting here instead of prepared statement, because prepared statements seems not to cope with coulumn name insertion
	// this should be save because non of the parameters is user input
	selectCurrentEvents := "SELECT id, username, env_name FROM reservations WHERE (status='%v') AND (%v<='" + time.Now().String() + "');"
	resRows, err := db.Query(fmt.Sprintf(selectCurrentEvents, status, timeCol))
	if err != nil {
		log.Println(err)
	}
	defer resRows.Close()

	var rows []reservationRow

	for resRows.Next() {
		var reservationID int
		var username, envName string
		err := resRows.Scan(&reservationID, &username, &envName)
		if err != nil {
			log.Fatal(err)
		}
		row := reservationRow{
			reservationID,
			username,
			envName,
			false,
		}
		rows = append(rows, row)
		fmt.Printf("Values from matching query: id - %v, username - %v, envname %v\n", reservationID, username, envName)
	}
	return rows
}

func handleBookings(db *sql.DB, environments map[string][]vault.SecEng, approle *vault.Approle) {
	updateState, err := db.Prepare("UPDATE reservations SET status=? WHERE id=?;")
	defer updateState.Close()
	if err != nil {
		log.Println(err)
	}
	// TODO: endless loop, triggered each 5 minutes

	log.Println("startet booking handle procedure")

	// have to start any upcoming bookings?
	rows := getRows(db, "upcoming", "start")
	for _, row := range rows {
		// check, if enironment in reservation exists and fill in the information has_ssh
		row.check(db)

		sshKey := ""
		if row.hasSSH {
			// retrieve ssh key from user table
			err := db.QueryRow("SELECT ssh_pub_key FROM users WHERE (username='" + row.username + "');").Scan(&sshKey)
			if err != nil {
				log.Fatal(err)
			}
		}

		// trigger the start of the booking
		vaultToken := approle.CreateVaultToken()
		vault.StartBooking(environments[row.envName], vaultToken, sshKey)

		// change booking status in database
		log.Println("executed start of booking")
		_, err = updateState.Exec("active", row.id)
		if err != nil {
			log.Printf("did not change status from upcoming to active due to following error: %v\n", err)
		}
	}

	// have to end any active bookings?
	rows = getRows(db, "active", "end")
	for _, row := range rows {
		// check, if enironment in reservation exists (and fill in the information has_ssh, which is not needed)
		row.check(db)

		// trigger the end of the booking
		vaultToken := approle.CreateVaultToken()
		vault.EndBooking(environments[row.envName], vaultToken)

		// change booking status in database
		log.Println("executed end of booking")
		_, err = updateState.Exec("expired", row.id)
		if err != nil {
			log.Printf("did not change status from active to expired due to following error: %v\n", err)
		}
	}

	// have to delete any expired bookings?
	rows = getRows(db, "expired", "delete_on")
	for _, row := range rows {

		// delete booking from database
		_, err = db.Exec("DELETE FROM reservations WHERE id=?;", row.id)
		if err != nil {
			log.Printf("did not delete database entry due to following error: %v\n", err)
		}
	}

	log.Println("end booking handle procedure")
}

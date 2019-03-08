package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gitlab-vs.informatik.uni-ulm.de/gafaspot/vault"
)

const (
	scanningInterval = 5 * time.Minute
)

func handleReservationScanning(db *sql.DB) {
	// endless loop, triggered each 5 minutes
	tick := time.NewTicker(scanningInterval).C
	for {
		<-tick
		reservationScan(db)
	}
}

type reservationRow struct {
	id       int
	username string
	envName  string
	end      time.Time
	hasSSH   bool
}

func (row reservationRow) check(tx *sql.Tx) bool {
	// Does environment exist? Does it has components which require an ssh key for login?
	err := tx.QueryRow("SELECT has_ssh FROM environments WHERE (env_name='" + row.envName + "');").Scan(&row.hasSSH)
	if err == sql.ErrNoRows {
		log.Printf("environment %v does not exis; delete reservation with id=%v for user=%v from database", row.envName, row.id, row.username)
		// delete booking from database
		_, err = tx.Exec("DELETE FROM reservations WHERE id=?;", row.id)
		if err != nil {
			log.Printf("did not delete database entry due to following error: %v\n", err)
		}
		return false
	} else if err != nil {
		log.Fatal(err)
	}
	return true
}

func getRows(tx *sql.Tx, now time.Time, status, timeCol string) []reservationRow {
	stmt, err := tx.Prepare("SELECT id, username, env_name, end FROM reservations WHERE (status=?) AND (" + timeCol + "<=?);")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	resRows, err := stmt.Query(status, now)
	if err != nil {
		log.Println(err)
	}
	defer resRows.Close()

	var rows []reservationRow

	for resRows.Next() {
		var reservationID int
		var username, envName string
		var end time.Time
		err := resRows.Scan(&reservationID, &username, &envName, &end)
		if err != nil {
			log.Fatal(err)
		}
		row := reservationRow{
			reservationID,
			username,
			envName,
			end,
			false,
		}
		rows = append(rows, row)
		fmt.Printf("Values from matching query: id - %v, username - %v, envname %v\n", reservationID, username, envName)
	}
	return rows
}

func reservationScan(db *sql.DB) {

	log.Println("startet booking handle procedure")

	now := time.Now()

	// any active bookings which should end?
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	rows := getRows(tx, now, "active", "end")
	for _, row := range rows {
		// check, if enironment in reservation exists (and fill in the information has_ssh, which is not needed)
		ok := row.check(tx)
		if ok {
			// trigger the end of the booking
			vaultToken := vault.CreateVaultToken()
			vault.EndBooking(row.envName, vaultToken)

			// change booking status in database
			changeStatus(tx, row.id, "expired")
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Println(err)
	}

	// any upcoming bookings which should start?
	tx, err = db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	rows = getRows(tx, now, "upcoming", "start")
	for _, row := range rows {

		if row.end.Before(now) {
			// in case the end time of the upcoming booking which never was active is already reached for some reason, don't start the booking, just expire it in database
			changeStatus(tx, row.id, "expired")
		} else {

			// check, if enironment in reservation exists and fill in the information has_ssh
			ok := row.check(tx)

			if ok {
				var sshKey sql.NullString
				if row.hasSSH {
					// retrieve ssh key from user table
					err := tx.QueryRow("SELECT ssh_pub_key FROM users WHERE (username='" + row.username + "');").Scan(&sshKey)
					if err == sql.ErrNoRows || !sshKey.Valid {
						log.Fatalf("there is no ssh public key stored for user %v, but it is required for booking environment %v", row.username, row.envName)
					} else if err != nil {
						log.Println(err)
					}
				}

				// trigger the start of the booking
				vaultToken := vault.CreateVaultToken()
				vault.StartBooking(row.envName, vaultToken, sshKey.String, row.end)

				// change booking status in database
				changeStatus(tx, row.id, "active")
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Println(err)
	}

	// any expired bookings which should get deleted?
	tx, err = db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	rows = getRows(tx, now, "expired", "delete_on")
	for _, row := range rows {

		// delete booking from database
		_, err = tx.Exec("DELETE FROM reservations WHERE id=?;", row.id)
		if err != nil {
			log.Printf("did not delete database entry due to following error: %v\n", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Println(err)
	}

	log.Println("end booking handle procedure")
}

func changeStatus(tx *sql.Tx, rowID int, status string) {
	_, err := tx.Exec("UPDATE reservations SET status=? WHERE id=?;", status, rowID)
	if err != nil {
		log.Printf("did not change status due to following error: %v\n", err)
	}
}

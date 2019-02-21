package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gitlab-vs.informatik.uni-ulm.de/gafaspot/constants"
	"gitlab-vs.informatik.uni-ulm.de/gafaspot/vault"
)

type reservationRow struct {
	id       int
	username string
	envName  string
	end      string
	hasSSH   bool
}

func (row reservationRow) check(tx *sql.Tx) {
	// Does environment exist? Does it has components which require an ssh key for login?
	err := tx.QueryRow("SELECT has_ssh FROM environments WHERE (env_name='" + row.envName + "');").Scan(&row.hasSSH)
	if err == sql.ErrNoRows {
		log.Fatalf("environment %v does not exist", row.envName)
	} else if err != nil {
		log.Fatal(err)
	}
}

func getRows(tx *sql.Tx, now, status, timeCol string) []reservationRow {
	// use sprintf formatting here instead of prepared statement, because prepared statements seems not to cope with coulumn name insertion
	// this should be save because non of the parameters is user input
	resRows, err := tx.Query(fmt.Sprintf("SELECT id, username, env_name, end FROM reservations WHERE (status='%v') AND (%v<='%v');", status, timeCol, now))
	if err != nil {
		log.Println(err)
	}
	defer resRows.Close()

	var rows []reservationRow

	for resRows.Next() {
		var reservationID int
		var username, envName, end string
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

func handleBookings(db *sql.DB, environments map[string][]vault.SecEng, approle *vault.Approle) {
	updateState, err := db.Prepare("UPDATE reservations SET status=? WHERE id=?;")
	defer updateState.Close()
	if err != nil {
		log.Println(err)
	}
	// TODO: endless loop, triggered each 5 minutes

	log.Println("startet booking handle procedure")

	now := time.Now().Format(constants.TimeLayout)

	// any active bookings which should end?
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	rows := getRows(tx, now, "active", "end")
	for _, row := range rows {
		// check, if enironment in reservation exists (and fill in the information has_ssh, which is not needed)
		row.check(tx)

		// trigger the end of the booking
		vaultToken := approle.CreateVaultToken()
		vault.EndBooking(environments[row.envName], vaultToken)

		// change booking status in database
		log.Println("executed end of booking")
		_, err = tx.Stmt(updateState).Exec("expired", row.id)
		if err != nil {
			log.Printf("did not change status from active to expired due to following error: %v\n", err)
		}
	}
	tx.Commit()

	// any upcoming bookings which should start?
	tx, err = db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	rows = getRows(tx, now, "upcoming", "start")
	for _, row := range rows {

		if row.end <= now {
			// in case the end time of the upcoming booking which never was active is already reached for some reason, don't start the booking, just expire it in database
			_, err = tx.Stmt(updateState).Exec("expired", row.id)
			if err != nil {
				log.Printf("did not change status from upcoming to expired due to following error: %v\n", err)
			}
		} else {

			// check, if enironment in reservation exists and fill in the information has_ssh
			row.check(tx)

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
			vaultToken := approle.CreateVaultToken()
			vault.StartBooking(environments[row.envName], vaultToken, sshKey.String, row.end)

			// change booking status in database
			log.Println("executed start of booking")
			_, err = tx.Stmt(updateState).Exec("active", row.id)
			if err != nil {
				log.Printf("did not change status from upcoming to active due to following error: %v\n", err)
			}
		}
	}
	tx.Commit()

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
	tx.Commit()

	log.Println("end booking handle procedure")
}

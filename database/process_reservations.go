package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gitlab-vs.informatik.uni-ulm.de/gafaspot/util"
)

func changeStatus(tx *sql.Tx, rowID int, status string) {
	_, err := tx.Exec("UPDATE reservations SET status=? WHERE id=?;", status, rowID)
	if err != nil {
		log.Printf("did not change status due to following error: %v\n", err)
	}
}

func deleteReservation(tx *sql.Tx, reservationID int) {
	_, err := tx.Exec("DELETE FROM reservations WHERE id=?;", reservationID)
	if err != nil {
		log.Printf("did not delete database entry due to following error: %v\n", err)
	}
}

func getApplicableReservations(tx *sql.Tx, now time.Time, status, timeCol string) []util.Reservation {
	stmt, err := tx.Prepare("SELECT id, username, env_name, end FROM reservations WHERE (status=?) AND (" + timeCol + "<=?);")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(status, now)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var reservations []util.Reservation

	for rows.Next() {
		r := util.Reservation{}

		err := rows.Scan(&r.ID, &r.User, &r.EnvName, &r.End)
		if err != nil {
			log.Fatal(err)
		}
		reservations = append(reservations, r)
		fmt.Printf("Values from matching query: id - %v, username - %v, envname %v\n", r.ID, r.User, r.EnvName)
	}
	return reservations
}

func check(tx *sql.Tx, r util.Reservation, hasSSH *bool) bool {
	// Does environment exist? Does it has components which require an ssh key for login?
	err := tx.QueryRow("SELECT has_ssh FROM environments WHERE (env_name='" + r.EnvName + "');").Scan(hasSSH)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		log.Fatal(err)
	}
	return true
}

type startBookingFunc func(envName, sshKey string, until time.Time)

func StartUpcomingReservations(now time.Time, startBooking startBookingFunc) {
	tx := beginTransaction()
	defer commitTransaction(tx)

	reservations := getApplicableReservations(tx, now, "upcoming", "start")
	for _, r := range reservations {

		if r.End.Before(now) {
			// in case the end time of the upcoming booking which never was active is already reached for some reason, don't start the booking, just expire it in database
			changeStatus(tx, r.ID, "expired")
			return
		}

		// check, if enironment in reservation exists and fill in the information has_ssh
		var hasSSH bool
		ok := check(tx, r, &hasSSH)
		if !ok {
			log.Printf("environment %v does not exis; mark reservation with id=%v for user=%v as error", r.EnvName, r.ID, r.User)
			changeStatus(tx, r.ID, "error")
			return
		}

		sshKey := ""
		if hasSSH {
			// retrieve ssh key from user table
			sshKey, ok = GetUserSSH(r.User)
			if !ok {
				log.Printf("there is no ssh public key stored for user %v, but it is required for booking environment %v", r.User, r.EnvName)
				changeStatus(tx, r.ID, "error")
				return
			}
		}

		// trigger the start of the booking
		startBooking(r.EnvName, sshKey, r.End)

		// change booking status in database
		changeStatus(tx, r.ID, "active")

	}
}

type endBookingFunc func(envName string)

func ExpireActiveReservations(now time.Time, endBooking endBookingFunc) {
	tx := beginTransaction()
	defer commitTransaction(tx)

	reservations := getApplicableReservations(tx, now, "active", "end")
	for _, r := range reservations {
		// check, if enironment in reservation exists (and fill in the information has_ssh, which is not needed)
		ok := check(tx, r, new(bool))
		if ok {
			// trigger the end of the booking
			endBooking(r.EnvName)

			// change booking status in database
			changeStatus(tx, r.ID, "expired")
		}
	}
}

func DeleteOldReservations(now time.Time) {
	tx := beginTransaction()
	defer commitTransaction(tx)

	reservations := getApplicableReservations(tx, now, "expired", "delete_on")
	reservations = append(reservations, getApplicableReservations(tx, now, "error", "delete_on")...)
	for _, r := range reservations {

		// delete booking from database
		deleteReservation(tx, r.ID)
	}
}

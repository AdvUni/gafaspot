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
	"database/sql"
	"fmt"
	"log"
	"time"

	"gitlab-vs.informatik.uni-ulm.de/gafaspot/util"
)

// changeStatus sets the status column of the reservation with the given id to the given status string.
func changeStatus(tx *sql.Tx, id int, status string) {
	_, err := tx.Exec("UPDATE reservations SET status=? WHERE id=?;", status, id)
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

// getApplicableReservations fetches specific reservations from database. Which reservations are
// concerned is specified by the reservation status and by the exceeding of a specific time column.
// For example, if you give the parameters status="upcoming" and timeCol="start", the function will
// return all upcoming reservation, which have a start date lieing in the past. Herefore, the
// reference time "now" has to be explicitely passed. tx is the transaction, in which the database
// request should be executed.
func getApplicableReservations(tx *sql.Tx, now time.Time, status, timeCol string) []util.Reservation {
	stmt, err := tx.Prepare("SELECT id, username, env_plain_name, end FROM reservations WHERE (status=?) AND (" + timeCol + "<=?);")
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

		err := rows.Scan(&r.ID, &r.User, &r.EnvPlainName, &r.End)
		if err != nil {
			log.Fatal(err)
		}
		reservations = append(reservations, r)
		fmt.Printf("Values from matching query: id - %v, username - %v, env_plain_name %v\n", r.ID, r.User, r.EnvPlainName)
	}
	return reservations
}

// When triggering a booking start or end, the environments table should be checked to be sure the
// environment in the reservation actually exists. At the same time, it can be
// retrieved the information, whether an ssh key is needed for booking the concerned environment as
// it might be needed later.
// tx is the transaction inside which the database query shall be executed.
// r is the reservation, for which the check should be performed.
// hasSSH is the pointer to a boolean. The function will fill the information, whether booking
// needs an ssh key in here.
// Function returns a boolean, which contains the check's outcome.
func check(tx *sql.Tx, r util.Reservation, hasSSH *bool) bool {
	// Does environment exist? Does it has components which require an ssh key for login?
	err := tx.QueryRow("SELECT has_ssh FROM environments WHERE (env_plain_name='" + r.EnvPlainName + "');").Scan(hasSSH)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		log.Fatal(err)
	}
	return true
}

type startBookingFunc func(envPlainName, sshKey string, until time.Time)

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
			log.Printf("environment %v does not exis; mark reservation with id=%v for user=%v as error", r.EnvPlainName, r.ID, r.User)
			changeStatus(tx, r.ID, "error")
			return
		}

		sshKey := ""
		if hasSSH {
			// retrieve ssh key from user table
			sshKey, ok = GetUserSSH(r.User)
			if !ok {
				log.Printf("there is no ssh public key stored for user %v, but it is required for booking environment %v", r.User, r.EnvPlainName)
				changeStatus(tx, r.ID, "error")
				return
			}
		}

		// trigger the start of the booking
		startBooking(r.EnvPlainName, sshKey, r.End)

		// change booking status in database
		changeStatus(tx, r.ID, "active")

	}
}

type endBookingFunc func(envPlainName string)

func ExpireActiveReservations(now time.Time, endBooking endBookingFunc) {
	tx := beginTransaction()
	defer commitTransaction(tx)

	reservations := getApplicableReservations(tx, now, "active", "end")
	for _, r := range reservations {
		// check, if enironment in reservation exists (and fill in the information has_ssh, which is not needed)
		ok := check(tx, r, new(bool))
		if ok {
			// trigger the end of the booking
			endBooking(r.EnvPlainName)

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

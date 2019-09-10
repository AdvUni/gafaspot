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
	"os"
	"time"

	"github.com/AdvUni/gafaspot/email"
	"github.com/AdvUni/gafaspot/util"
)

// changeStatus sets the status column of the reservation with the given id to the given status string.
func changeStatus(tx *sql.Tx, id int, status string) {
	_, err := tx.Exec("UPDATE reservations SET status=? WHERE id=?;", status, id)
	if err != nil {
		logger.Error("did not change status due to following error: %v\n", err)
	}
}

func deleteReservation(tx *sql.Tx, reservationID int) {
	_, err := tx.Exec("DELETE FROM reservations WHERE id=?;", reservationID)
	if err != nil {
		logger.Error("did not delete database entry due to following error: %v\n", err)
	}
}

// getApplicableReservations fetches specific reservations from database. Which reservations are
// concerned is specified by the reservation status and by the exceeding of a specific time column.
// For example, if you give the parameters status="upcoming" and timeCol="start", the function will
// return all upcoming reservation, which have a start date lieing in the past. Herefore, the
// reference time "now" has to be explicitly passed. tx is the transaction, in which the database
// request should be executed.
func getApplicableReservations(tx *sql.Tx, now time.Time, status, timeCol string) []util.Reservation {
	stmt, err := tx.Prepare("SELECT id, status, username, env_plain_name, start, end, subject, labels, start_mail, end_mail FROM reservations WHERE (status=?) AND (" + timeCol + "<=?);")
	if err != nil {
		logger.Emergency(err)
		os.Exit(1)
	}
	defer stmt.Close()

	rows, err := stmt.Query(status, now)
	if err != nil {
		logger.Error(err)
	}
	defer rows.Close()

	return assembleReservations(rows)
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
		logger.Error(err)
	}
	return true
}

type startBookingFunc func(envPlainName, sshKey string, until time.Time)

// StartUpcomingReservations selects all upcoming reservations from database, wich have a start
// time smaller than now. It applies the startBooking function to all environments which are
// affected by those reservations. After, it changes the reservation's status in database.
// The reason, why the startBooking function is passed as parameter
// here is the ambition to preserve the separation of database and vault package. The time now is
// passed because an unchanging reference is needed over several function calls to avoid
// inconsistencies.
func StartUpcomingReservations(now time.Time, startBooking startBookingFunc) {
	tx := beginTransaction()
	defer commitTransaction(tx)

	reservations := getApplicableReservations(tx, now, "upcoming", "start")
	for _, r := range reservations {

		if r.End.Before(now) {
			// in case the end time of the upcoming booking which never was active is already reached for some reason, don't start the booking, just expire it in database
			changeStatus(tx, r.ID, "expired")
			// TODO: Possibly write an email?
			return
		}

		// check, if environment in reservation exists and fill in the information has_ssh
		var hasSSH bool
		ok := check(tx, r, &hasSSH)
		if !ok {
			logger.Warningf("environment %v does not exist; mark reservation with id=%v for user=%v as error", r.EnvPlainName, r.ID, r.User)
			changeStatus(tx, r.ID, "error")
			// TODO: Possibly write an email?
			return
		}

		sshKey := ""
		if hasSSH {
			// retrieve ssh key from user table
			sshKey, ok = GetUserSSH(r.User)
			if !ok {
				logger.Warningf("there is no ssh public key stored for user %v, but it is required for booking environment %v. Mark reservation with error",
					r.User, r.EnvPlainName)
				changeStatus(tx, r.ID, "error")
				return
			}
		}

		// trigger the start of the booking
		logger.Infof("Starting reservation... %+v", r)
		startBooking(r.EnvPlainName, sshKey, r.End)

		// change booking status in database
		changeStatus(tx, r.ID, "active")

		// send email to user, if wished and if mailing is enabled in gafaspot config
		if r.SendStartMail && email.MailingEnabled {
			mailAddress, ok := GetUserEmail(r.User)
			if ok {
				email.SendBeginReservationMail(mailAddress, r)
			} else {
				logger.Warningf("tried to send an e-mail to user '%s', but there is not mail address stored for him in database (anymore)", r.User)
			}
		}

	}
}

type endBookingFunc func(envPlainName string)

// ExpireActiveReservations selects all active reservations from database, wich have an end
// time smaller than now. It applies the endBooking function to all environments which are
// affected by those reservations. After, it changes the reservation's status in database.
// The reason, why the endBooking function is passed as parameter
// here is the ambition to preserve the separation of database and vault package. The time now is
// passed because an unchanging reference is needed over several function calls to avoid
// inconsistencies.
func ExpireActiveReservations(now time.Time, endBooking endBookingFunc) {
	tx := beginTransaction()
	defer commitTransaction(tx)

	reservations := getApplicableReservations(tx, now, "active", "end")
	for _, r := range reservations {
		// check, if environment in reservation exists (and fill in the information has_ssh, which is not needed)
		ok := check(tx, r, new(bool))
		if ok {
			// trigger the end of the booking
			logger.Infof("Ending reservation... %+v", r)
			endBooking(r.EnvPlainName)
		} else {
			logger.Infof("Ended reservation for an environment, which does not seam to exist (anymore): %+v", r)
		}
		// change booking status in database
		changeStatus(tx, r.ID, "expired")

		// send email to user, if wished and if mailing is enabled in gafaspot config
		if r.SendEndMail && email.MailingEnabled {
			mailAddress, ok := GetUserEmail(r.User)
			if ok {
				email.SendEndReservationMail(mailAddress, r)
			} else {
				logger.Warningf("tried to send an e-mail to user '%s', but there is not mail address stored for him in database (anymore)", r.User)
			}
		}
	}
}

// DeleteOldReservations selects all expired reservations from database, which have a delete_on time
// smaller than now. It deletes all those reservations from database.
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

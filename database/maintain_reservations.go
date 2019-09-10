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
	"os"
	"time"

	"github.com/AdvUni/gafaspot/email"

	"github.com/AdvUni/gafaspot/util"
)

// ReservationError is thrown by the CreateReservation func, if reservation parameters are not
// valid. In this case, reservation will not be created.
type ReservationError string

func (err ReservationError) Error() string {
	return fmt.Sprintf("reservation is invalid: %v", string(err))
}

// CreateReservation puts a new reservation entry to the database. Bevor writing to database,
// several checks are performed. Function checks time parameters for plausibility, tests, if
// user has an ssh key uploaded if necessary, and checks for possible conflicts with existing
// reservations. If everything is fine, reservation will be created. Otherwise, function returns
// a reservation error.
func CreateReservation(r util.Reservation) error {

	// check, whether reservation is in future
	if !r.Start.After(time.Now()) {
		return ReservationError("cannot do reservation for the past")
	}

	// check whether start < end
	if !r.Start.Before(r.End) {
		return ReservationError("end of reservation must be after start of reservation")
	}

	// check whether reservation duration is too long
	if r.Start.AddDate(0, 0, maxBookingDays).Before(r.End) {
		return ReservationError(fmt.Sprintf("you are only allowed to do reservations with a duration up to %v days", maxBookingDays))
	}

	// check whether reservation is too far in the future
	if time.Now().AddDate(0, maxQueuingMonths, 0).Before(r.Start) {
		return ReservationError(fmt.Sprintf("you are not allowed to do reservations which start more than %v months int the future", maxQueuingMonths))
	}

	// start a transaction
	tx := beginTransaction()
	defer commitTransaction(tx)

	// check, whether environment exists and determine, whether the reservation needs an ssh key
	stmt, err := tx.Prepare("SELECT has_ssh FROM environments WHERE (env_plain_name=?);")
	if err != nil {
		logger.Emergency(err)
		os.Exit(1)
	}
	defer stmt.Close()

	var hasSSH bool
	err = stmt.QueryRow(r.EnvPlainName).Scan(&hasSSH)
	if err == sql.ErrNoRows {
		return ReservationError(fmt.Sprintf("environment %v does not exist", r.EnvPlainName))
	} else if err != nil {
		logger.Emergency(err)
		os.Exit(1)
	}

	// check, whether there is stored an ssh key for the user, if it is needed for the reservation
	if hasSSH {
		if !UserHasSSH(r.User) {
			return ReservationError(fmt.Sprintf("there is no ssh public key stored for user %v, but it is required for booking environment %v", r.User, r.EnvPlainName))
		}
	}

	// check possibility of sending e-mails
	if r.SendEndMail || r.SendStartMail {
		if !email.MailingEnabled {
			return ReservationError("gafaspot is not configured to send e-mails")
		}
		if !UserHasEmail(r.User) {
			return ReservationError(fmt.Sprintf("there is no e-mail address stored for user %v, so Gafaspot can't mail him", r.User))
		}
	}

	// check the environment's availability within the requested time range:
	// a conflict occurs iff ((start1 <= end2) && (end1 >= start2))
	stmt, err = tx.Prepare("SELECT start, end FROM reservations WHERE (env_plain_name=?) AND (start<=?) AND (end>=?);")
	if err != nil {
		logger.Emergency(err)
		os.Exit(1)
	}
	defer stmt.Close()

	var conflictStart, conflictEnd time.Time
	err = stmt.QueryRow(r.EnvPlainName, r.End, r.Start).Scan(&conflictStart, &conflictEnd)
	// there is a conflict, if answer is NOT empty; means, if there is NO sql.ErrNoRows
	if err == nil {
		return ReservationError(fmt.Sprintf("reservation conflicts with an existing reservation from %v to %v", conflictStart.Format(util.TimeLayout), conflictEnd.Format(util.TimeLayout)))
	}
	if err != sql.ErrNoRows {
		logger.Error(err)
	}

	// generate the deletion date of reservation entry in database
	reservationDeleteDate := addTTL(r.End)

	// finally write reservation into database
	stmt, err = tx.Prepare("INSERT INTO reservations (status, username, env_plain_name, start, end, subject, labels, start_mail, end_mail, delete_on) VALUES(?,?,?,?,?,?,?,?,?,?);")
	if err != nil {
		logger.Emergency(err)
		os.Exit(1)
	}
	defer stmt.Close()
	_, err = stmt.Exec("upcoming", r.User, r.EnvPlainName, r.Start, r.End, r.Subject, r.Labels, r.SendStartMail, r.SendEndMail, reservationDeleteDate)
	if err != nil {
		logger.Error(err)
	} else {
		logger.Infof("new reservation created: %+v", r)
	}

	return nil
}

// AbortReservation deletes a reservation entry from database. This is only possible, if the
// reservation is still upcoming and not active yet. This is because an active reservation
// has to be ended, whereas an upcoming reservation just can be deleted. Further, a reservation
// is only abortable by the user who created it.
// Function parameter id is the reservation's database id.
func AbortReservation(username string, id int) error {
	// start a transaction
	tx := beginTransaction()
	defer commitTransaction(tx)

	// fetch reservation from database
	stmt, err := tx.Prepare("SELECT status FROM reservations WHERE (username=?) AND (id=?);")
	if err != nil {
		logger.Emergency(err)
		os.Exit(1)
	}
	defer stmt.Close()
	status := ""
	err = stmt.QueryRow(username, id).Scan(&status)
	if err == sql.ErrNoRows {
		logger.Warning(fmt.Errorf("tried to abort reservation which does not exist or not belongs to specified user; id '%v', user '%v'", id, username))
		return nil
	}
	if err != nil {
		logger.Error(err)
	}

	// check reservation status (can only abort upcoming reservations)
	if status != "upcoming" {
		return fmt.Errorf("reservation is already active or expired, though it is not possible anymore to abort it")
	}

	// delete reservation from database
	deleteReservation(tx, id)

	return nil
}

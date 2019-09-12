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
	"html/template"
	"os"
	"sort"

	"github.com/AdvUni/gafaspot/util"
)

// UserHasSSH determines, whether an SSH public key is stored in database for a given username
func UserHasSSH(username string) bool {
	_, ok := GetUserSSH(username)
	return ok
}

// GetUserSSH returns the SSH public key for a given user from database if present. If not, an
// empty string will be returned, together with the second return value saying false.
func GetUserSSH(username string) (string, bool) {
	return getUserAttribute(username, "ssh_pub_key")
}

// UserHasEmail determines, whether an e-mail address is stored in database for a given username
func UserHasEmail(username string) bool {
	_, ok := GetUserEmail(username)
	return ok
}

// GetUserEmail returns the mail address for a given user from database if present. If not, an
// empty string will be returned, together with the second return value saying false.
func GetUserEmail(username string) (string, bool) {
	return getUserAttribute(username, "email")
}

func getUserAttribute(username, attribute string) (string, bool) {
	var value sql.NullString
	stmtstring := fmt.Sprintf("SELECT %s FROM users WHERE (username=?);", attribute)
	stmt, err := db.Prepare(stmtstring)
	if err != nil {
		logger.Emergency(err)
		os.Exit(1)
	}
	defer stmt.Close()
	err = stmt.QueryRow(username).Scan(&value)
	if err == sql.ErrNoRows || !value.Valid {
		return "", false
	} else if err != nil {
		logger.Error(err)
		return "", false
	}
	if value.String == "" {
		return "", false
	}
	return value.String, true
}

// GetEnvironments reads all environments from database and returns them as a map with the PlainNames as keys.
func GetEnvironments() map[string]util.Environment {
	rows, err := db.Query("SELECT env_plain_name, env_nice_name, has_ssh, description FROM environments ORDER BY env_nice_name;")
	if err != nil {
		logger.Error(err)
		return nil
	}
	defer rows.Close()

	envMap := make(map[string]util.Environment)
	for rows.Next() {
		e := util.Environment{}
		description := sql.NullString{}
		err := rows.Scan(&e.PlainName, &e.NiceName, &e.HasSSH, &description)
		if err != nil {
			logger.Emergency(err)
			os.Exit(1)
		}
		if description.Valid {
			e.Description = template.HTML(description.String)
		}

		envMap[e.PlainName] = e
	}
	return envMap
}

// GetEnvReservations returns all reservations stored in database for a specific environment.
func GetEnvReservations(envPlainName string) []util.Reservation {
	return getReservations("env_plain_name", envPlainName)
}

// GetUserReservations returns all reservations stored in database for a specific username.
func GetUserReservations(username string) []util.Reservation {
	return getReservations("username", username)
}

// getReservations allows to select all reservations from database by one specific condition. The
// condition is: 'WHERE conditionKey=conditionVal', where conditionKey and conditionVal are
// function parameters.
func getReservations(conditionKey, conditionVal string) []util.Reservation {
	stmtstring := fmt.Sprintf("SELECT id, status, username, env_plain_name, start, end, subject, labels, start_mail, end_mail FROM reservations WHERE %v=?", conditionKey)
	stmt, err := db.Prepare(stmtstring)
	if err != nil {
		logger.Emergency(err)
		os.Exit(1)
	}
	defer stmt.Close()

	rows, err := stmt.Query(conditionVal)
	if err != nil {
		logger.Error(err)
		return nil
	}
	defer rows.Close()
	return assembleReservations(rows)
}

// CollectUserCreds bundles all valid credentials for a user. It searches for the user's
// reservations with status 'active', adds the Environment information and looks up the
// corresponding credentials.
// As reading credentials from vault is a matter of the vault package, and it is tried to
// keep the packages database and vault separately, the readCreds function is passed as
// parameter.
// If a reservation is found for which no environment exists in database, the function
// creates kind of a dummy Environment struct using the EnvPlainName given in the Reservation.
// No error or similar will arise.
func CollectUserCreds(username string, readCreds readCredsFunc) []util.ReservationCreds {
	// get all active reservations of user
	stmt, err := db.Prepare("SELECT id, status, username, env_plain_name, start, end, subject, labels, start_mail, end_mail FROM reservations WHERE (status='active') AND (username=?);")
	if err != nil {
		logger.Emergency(err)
		os.Exit(1)
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		logger.Error(err)
	}
	defer rows.Close()
	reservations := assembleReservations(rows)

	// sort reservations
	sort.Slice(reservations, func(i, j int) bool {
		return reservations[i].EnvPlainName < reservations[j].EnvPlainName
	})
	// add environment info
	resEnvCreds := collateReservationEnvironment(reservations)

	// add creds info
	for i := range resEnvCreds {
		resEnvCreds[i].Creds = readCreds(resEnvCreds[i].Env.PlainName)
	}

	return resEnvCreds
}

// collectReservationCreds bundles the credentials for an active reservation.
// The function adds the Environment information to the reservation and looks up the
// corresponding credentials.
// As reading credentials from vault is a matter of the vault package, and it is tried to
// keep the packages database and vault separately, the readCreds function is passed as
// parameter.
// If the reservation's environment does not exist in database, the function creates
// kind of a dummy Environment struct using the EnvPlainName given in the Reservation.
// No error or similar will arise.
// Pass only active reservations, otherwise the readCreds function will not be able to
// return proper values.
func collectReservationCreds(reservation util.Reservation, readCreds readCredsFunc) util.ReservationCreds {
	reservationCreds := collateReservationEnvironment([]util.Reservation{reservation})[0]
	reservationCreds.Creds = readCreds(reservation.EnvPlainName)
	return reservationCreds
}

// collateReservationEnvironment takes a list of Reservations and looks up the Environment
// for each Reservation. The function returns both packed together in a list of
// util.ReservationCreds structs, in which the Creds attribute is not set.
// If a reservation is found for which no environment exists in database, the function
// creates kind of a dummy Environment struct using the EnvPlainName given in the Reservation.
// No error or similar will arise.
// The size of the result list is the same es the size of reservations
func collateReservationEnvironment(reservations []util.Reservation) []util.ReservationCreds {
	// get all environments
	environments := GetEnvironments()

	// associate environments and reservations to each other
	// if environment does not exist, create a dummy Environment object which contains the PlainName
	activeReservationsInfo := []util.ReservationCreds{}
	for _, r := range reservations {
		env, ok := environments[r.EnvPlainName]
		if !ok {
			logger.Debugf("environment from reservation does not exist. Reservation: %v", r)
			env = util.Environment{NiceName: r.EnvPlainName, PlainName: r.EnvPlainName}
		}
		//creds := vault.ReadCredentials(r.EnvPlainName)
		activeReservationsInfo = append(activeReservationsInfo, util.ReservationCreds{Res: r, Env: env})
	}
	return activeReservationsInfo
}

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

func userHasEmail(username string) bool {
	_, ok := GetUserEmail(username)
	return ok
}

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
	return value.String, true
}

// GetEnvironments reads all environments from database and returns them in a slice, ordered by
// NiceName. Further, it returns a map to access the environments by PlainName.
func GetEnvironments() ([]util.Environment, map[string]util.Environment) {
	rows, err := db.Query("SELECT env_plain_name, env_nice_name, has_ssh, description FROM environments ORDER BY env_nice_name;")
	if err != nil {
		logger.Error(err)
		return nil, nil
	}
	defer rows.Close()

	envs := []util.Environment{}
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

		envs = append(envs, e)
		envMap[e.PlainName] = e
	}
	return envs, envMap
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

// GetUserActiveReservationEnv returns all environment's PlainNames, which are stored for a
// specific username and have status 'active'.
func GetUserActiveReservationEnv(username string) []util.Reservation {
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

	return assembleReservations(rows)
}

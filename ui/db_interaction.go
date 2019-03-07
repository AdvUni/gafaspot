package ui

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gitlab-vs.informatik.uni-ulm.de/gafaspot/constants"
)

const (
	yearsTTL = 2
)

type ReservationError string

func (err ReservationError) Error() string {
	return fmt.Sprintf("reservation is invalid: %v", string(err))
}

func CreateReservation(db *sql.DB, username, envName, subject, labels string, start, end time.Time) error {

	// check, whether reservation is in future
	if !start.After(time.Now()) {
		return ReservationError("cannot do reservation for the past")
	}

	// check whether start < end
	if !start.Before(end) {
		return ReservationError("end of reservation must be after start of reservation")
	}

	// start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	// defer the transaction's commit as this function may return an error
	defer func() {
		err = tx.Commit()
		if err != nil {
			log.Println(err)
		}
	}()

	// check, whether environment exists and determine, whether the reservation needs an ssh key
	stmt, err := tx.Prepare("SELECT has_ssh FROM environments WHERE (env_name=?);")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var hasSSH bool
	err = stmt.QueryRow(envName).Scan(&hasSSH)
	if err == sql.ErrNoRows {
		return ReservationError(fmt.Sprintf("environment %v does not exist", envName))
	} else if err != nil {
		log.Fatal(err)
	}

	// check, whether there is stored an ssh key for the user, if it is needed for the reservation
	if hasSSH {
		var sshKey sql.NullString
		stmt, err = tx.Prepare("SELECT ssh_pub_key FROM users WHERE (username=?);")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		err = stmt.QueryRow(username).Scan(&sshKey)
		if err == sql.ErrNoRows || !sshKey.Valid {
			return ReservationError(fmt.Sprintf("there is no ssh public key stored for user %v, but it is required for booking environment %v", username, envName))
		} else if err != nil {
			log.Println(err)
		}
	}

	// check the environment's availability within the requested time range:
	// a conflict occurs iff ((start1 <= end2) && (end1 >= start2))
	stmt, err = tx.Prepare("SELECT start, end FROM reservations WHERE (env_name=?) AND (start<=?) AND (end>=?);")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var conflictStart, conflictEnd time.Time
	err = stmt.QueryRow(envName, end, start).Scan(&conflictStart, &conflictEnd)
	// there is a conflict, if answer is NOT empty; means, if there is NO sql.ErrNoRows
	if err == nil {
		return ReservationError(fmt.Sprintf("reservation conflicts with an existing reservation from %v to %v", conflictStart.Format(constants.TimeLayout), conflictEnd.Format(constants.TimeLayout)))
	}
	if err != sql.ErrNoRows {
		log.Fatal(err)
	}

	// generate the deletion date of reservation entry in database
	reservationDeleteDate := end.AddDate(yearsTTL, 0, 0)

	// finally write reservation into database
	stmt, err = tx.Prepare("INSERT INTO reservations (status, username, env_name, start, end, subject, labels, delete_on) VALUES(?,?,?,?,?,?,?,?);")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec("upcoming", username, envName, start, end, subject, labels, reservationDeleteDate)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func AbortReservation(db *sql.DB, username, string, id int) error {
	// start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	// defer the transaction's commit as this function may return an error
	defer func() {
		err = tx.Commit()
		if err != nil {
			log.Println(err)
		}
	}()

	// fetch reservation from database
	stmt, err := tx.Prepare("SELECT status FROM reservations WHERE (username=?) AND (id=?);")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	status := ""
	err = stmt.QueryRow(username, id).Scan(&status)
	if err == sql.ErrNoRows {
		log.Println(fmt.Errorf("tried to abort reservation which does not exist or not belongs to specified user; id '%v', user '%v'", id, username))
		return nil
	}
	if err != nil {
		log.Fatal(err)
	}

	// check reservation status (can only abort upcoming reservations)
	if status != "upcoming" {
		return fmt.Errorf("reservation is already active or expired, though it is not possible anymore to abort it")
	}

	// delete reservation from database
	stmt, err = tx.Prepare("DELETE FROM reservations WHERE id=?;")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func getEnvironments(db *sql.DB) ([]env, map[string]env) {
	rows, err := db.Query("SELECT env_name, has_ssh, description FROM environments;")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	envList := []env{}
	envMap := make(map[string]env)
	for rows.Next() {
		e := env{}
		description := sql.NullString{}
		err := rows.Scan(&e.Name, &e.HasSSH, &description)
		if err != nil {
			log.Fatal(err)
		}
		if description.Valid {
			e.Description = description.String
		}
		e.NamePlain = createPlainIdentifier(e.Name)

		envList = append(envList, e)
		envMap[e.NamePlain] = e
	}
	return envList, envMap
}

func getEnvReservations(db *sql.DB, envName string) []reservation {
	return getReservations(db, "env_name", envName)
}

func getUserReservations(db *sql.DB, username string) []reservation {
	return getReservations(db, "username", username)
}

func getReservations(db *sql.DB, conditionKey, conditionVal string) []reservation {
	stmtstring := fmt.Sprintf("SELECT id, status, username, env_name, start, end, subject, labels FROM reservations WHERE %v=?", conditionKey)
	stmt, err := db.Prepare(stmtstring)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(conditionVal)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	reservations := []reservation{}
	for rows.Next() {
		r := reservation{}
		var starttime, endtime time.Time
		var subject, labels sql.NullString
		err := rows.Scan(&r.ID, &r.Status, &r.User, &r.EnvName, &starttime, &endtime, &subject, &labels)
		if err != nil {
			log.Fatal(err)
		}
		r.Start = starttime.Format(constants.TimeLayout)
		r.End = endtime.Format(constants.TimeLayout)
		if subject.Valid {
			r.Subject = subject.String
		}
		if labels.Valid {
			r.Labels = labels.String
		}

		reservations = append(reservations, r)
	}
	return reservations
}

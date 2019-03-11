package database

import (
	"database/sql"
	"fmt"
	"log"

	"gitlab-vs.informatik.uni-ulm.de/gafaspot/util"
)

func UserHasSSH(username string) bool {
	_, ok := GetUserSSH(username)
	return ok
}

func GetUserSSH(username string) (string, bool) {
	var sshKey sql.NullString
	stmt, err := db.Prepare("SELECT ssh_pub_key FROM users WHERE (username=?);")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(username).Scan(&sshKey)
	if err == sql.ErrNoRows || !sshKey.Valid {
		return "", false
	} else if err != nil {
		log.Println(err)
		return "", false
	}
	return sshKey.String, true
}

func GetEnvironments() ([]util.Environment, map[string]bool) {
	rows, err := db.Query("SELECT env_plain_name, env_nice_name, has_ssh, description FROM environments;")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	envs := []util.Environment{}
	envHasSSHMap := make(map[string]bool)
	for rows.Next() {
		e := util.Environment{}
		description := sql.NullString{}
		err := rows.Scan(&e.PlainName, &e.NiceName, &e.HasSSH, &description)
		if err != nil {
			log.Fatal(err)
		}
		if description.Valid {
			e.Description = description.String
		}

		envs = append(envs, e)
		envHasSSHMap[e.PlainName] = e.HasSSH
	}
	return envs, envHasSSHMap
}

func GetEnvReservations(envPlainName string) []util.Reservation {
	return getReservations("env_plain_name", envPlainName)
}

func GetUserReservations(username string) []util.Reservation {
	return getReservations("username", username)
}

func getReservations(conditionKey, conditionVal string) []util.Reservation {
	stmtstring := fmt.Sprintf("SELECT id, status, username, env_plain_name, start, end, subject, labels FROM reservations WHERE %v=?", conditionKey)
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

	reservations := []util.Reservation{}
	for rows.Next() {
		r := util.Reservation{}
		var subject, labels sql.NullString
		err := rows.Scan(&r.ID, &r.Status, &r.User, &r.EnvPlainName, &r.Start, &r.End, &subject, &labels)
		if err != nil {
			log.Fatal(err)
		}
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

func GetUserActiveReservationEnv(username string) []string {
	stmt, err := db.Prepare("SELECT env_plain_name FROM reservations WHERE (status='active') AND (username=?);")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var envNames []string
	for rows.Next() {
		env := ""
		err := rows.Scan(&env)
		if err != nil {
			log.Fatal(err)
		}
		envNames = append(envNames, env)
	}
	return envNames
}

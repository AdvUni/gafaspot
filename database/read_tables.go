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

func GetEnvironments() ([]util.Environment, map[string]util.Environment) {
	rows, err := db.Query("SELECT env_name, has_ssh, description FROM environments;")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	envList := []util.Environment{}
	envMap := make(map[string]util.Environment)
	for rows.Next() {
		e := util.Environment{}
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

func GetEnvReservations(envName string) []util.Reservation {
	return getReservations("env_name", envName)
}

func GetUserReservations(username string) []util.Reservation {
	return getReservations("username", username)
}

func getReservations(conditionKey, conditionVal string) []util.Reservation {
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

	reservations := []util.Reservation{}
	for rows.Next() {
		r := util.Reservation{}
		var subject, labels sql.NullString
		err := rows.Scan(&r.ID, &r.Status, &r.User, &r.EnvName, &r.Start, &r.End, &subject, &labels)
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
	stmt, err := db.Prepare("SELECT env_name FROM reservations WHERE (status='active') AND (username=?);")
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

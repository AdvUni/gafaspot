package vault

import (
	"database/sql"
	"log"
	"time"
)

func handle_bookings(db *sql.DB) {

	stmt, err := db.Prepare("SELECT ? FROM ? WHERE ?;")

	// TODO: endless loop, triggered each 5 minutes

	// have to start any upcoming bookings?
	resRows, err := stmt.Query("(id, env_name)", "reservations", "(status='upcoming') AND (start<='"+string(time.Now)+"')")
	if err != nil {
		log.Println(err)
	}
	defer resRows.Close()
	for resRows.Next() {
		var id int
		var username, envName string
		err := resRows.Scan(&id, &username, &envName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, username, envName)

		// Does environment exist? Is ssh key needed?
		var hasSSH bool
		err := stmt.QueryRow("has_ssh", "environments", "(env_name='"+envName+"')").Scan(&hasSSH)
		if err == sql.ErrNoRows {
			log.Fatal("Environment %v does not exist. Can't book it.", envName)
		} else if err != nil {
			log.Fatal(err)
		}

		sshKey := ""
		if hasSSH {
			err := stmt.QueryRow("ssh_pub_key", "users", "(username='"+username+"')").Scan(&sshKey)
			if err != nil {
				log.Fatal(err)
			}
		}
		
		// TODO: finally submit booking. From where is the secret engine list retrieved?

	}

	// TODO: have to end any active bookings?

	// TODO: have to delete any expired bookings?
}

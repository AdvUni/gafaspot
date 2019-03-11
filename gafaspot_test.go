package main

import (
	"testing"

	"gitlab-vs.informatik.uni-ulm.de/gafaspot/util"
)

func mockConfig() util.GafaspotConfig {

	envconfig := make(map[string]util.EnvironmentConfig)
	envconfig["demo0"] = util.EnvironmentConfig{SecretEngines: []util.SecretEngineConfig{}}
	envconfig["demo1"] = util.EnvironmentConfig{SecretEngines: []util.SecretEngineConfig{}}
	return util.GafaspotConfig{Database: "./gafaspot-test.db", Environments: envconfig}

}

func TestCreateInvalidReservations(t *testing.T) {
	/*dummyconfig := mockConfig()
	// delete database file if it already exists
	os.Remove(dummyconfig.Database)

	database.InitDB(dummyconfig)

	now := time.Now()
	past := now.Add(-10 * time.Minute)
	future1 := now.Add(10 * time.Minute)
	future2 := now.Add(20 * time.Minute)

	err := ui.CreateReservation(db, "testuser", "demo0", "", "", past, future1)
	t.Log(err)
	_, ok := err.(ui.ReservationError)
	if !ok {
		t.Fail()
	}
	err = ui.CreateReservation(db, "testuser", "demo0", "", "", future2, future1)
	t.Log(err)
	_, ok = err.(ui.ReservationError)
	if !ok {
		t.Fail()
	}
	err = ui.CreateReservation(db, "testuser", "not a valid env", "", "", future1, future2)
	t.Log(err)
	_, ok = err.(ui.ReservationError)
	if !ok {
		t.Fail()
	}
	*/
}

func TestReservationStateRotation(t *testing.T) {
	/*

		dummyconfig := mockConfig()
		// delete database file if it already exists
		os.Remove(dummyconfig.Database)

		db := initDB(dummyconfig)
		defer db.Close()

		now := time.Now()
		time1 := now.Add(1 * time.Millisecond)
		time2 := time1.Add(100 * time.Millisecond)

		err := ui.CreateReservation(db, "testuser", "demo0", "", "", time1, time2)
		if err != nil {
			t.Fatal(err)
		}

		var id int
		var status string
		err = db.QueryRow("select id, status from reservations;").Scan(&id, &status)
		t.Log(id)
		if err != nil {
			t.Fatal(err)
		}
		if status != "upcoming" {
			t.Fatalf("expected status 'upcoming', instead got '%v'", status)
		}

		time.Sleep(2 * time.Millisecond)
		reservationScan(db)

		err = db.QueryRow(fmt.Sprintf("select status from reservations where id='%v';", id)).Scan(&status)
		if err != nil {
			t.Fatal(err)
		}
		if status != "active" {
			t.Fatalf("expected status 'active', instead got '%v'", status)
		}

		time.Sleep(100 * time.Millisecond)
		reservationScan(db)

		err = db.QueryRow(fmt.Sprintf("select status from reservations where id='%v';", id)).Scan(&status)
		if err != nil {
			t.Fatal(err)
		}
		if status != "expired" {
			t.Fatalf("expected status 'expired', instead got '%v'", status)
		}

		just := time.Now().Add(-2 * time.Millisecond)
		_, err = db.Exec(fmt.Sprintf("UPDATE reservations SET delete_on='%v' WHERE id=%v;", just, id))
		reservationScan(db)

		err = db.QueryRow(fmt.Sprintf("select status from reservations where id='%v';", id)).Scan(&status)
		if err != sql.ErrNoRows {
			t.Fatalf("expected no database entry for id %v, but database returned entry, including status %v", id, status)
		}*/

}

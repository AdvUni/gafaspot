package main

import (
	"log"
	"time"

	"gitlab-vs.informatik.uni-ulm.de/gafaspot/database"
	"gitlab-vs.informatik.uni-ulm.de/gafaspot/vault"
)

const (
	scanningInterval = 5 * time.Minute
)

func handleReservationScanning() {
	// endless loop, triggered each 5 minutes
	tick := time.NewTicker(scanningInterval).C
	for {
		<-tick
		reservationScan()
	}
}

func reservationScan() {

	now := time.Now()

	// any active bookings which should end?
	database.ExpireActiveReservations(now, vault.EndBooking)

	// any upcoming bookings which should start?
	database.StartUpcomingReservations(now, vault.StartBooking)

	// any expired bookings which should get deleted?
	database.DeleteOldReservations(now)

	// finally, check if some of the entries in users table reached deletion_date
	database.DeleteOldUserEntries(now)

	log.Println("finished reservation scan")
}

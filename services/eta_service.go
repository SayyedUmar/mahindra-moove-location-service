package services

import (
	"time"

	"github.com/MOOVE-Network/location_service/db"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func handleCheckinTrip(trip *db.Trip, currentLocation db.Location) error {
	ds := GetDurationService()
	ns := GetNotificationService()
	// TODO: Verify if driver_arrived should be one of the all checked in statuses
	if trip.AllCheckedIn() {
		endLocation := trip.TripRoutes[len(trip.TripRoutes)-1].ScheduledEndLocation
		dm, err := ds.GetDuration(currentLocation, endLocation, time.Now())
		if err != nil {
			return err
		}
		for _, tr := range trip.TripRoutes {
			go NotifyTripRoute(&tr, &dm, ns)
		}
		return nil
	}
	var offset time.Duration
	var trsToBeNotified []db.TripRoute
	lastDurationMetric := DurationMetrics{}
	for _, tr := range trip.TripRoutes {
		if tr.Status == "on_board" {
			// Notify them the last
			trsToBeNotified = append(trsToBeNotified, tr)
		}
		if tr.Status == "not_started" && offset == 0 {
			startLoc := currentLocation
			endLoc := tr.ScheduledStartLocation
			dm, err := ds.GetDuration(startLoc, endLoc, time.Now())
			if err != nil {
				return err
			}
			offset += dm.Duration
			lastDurationMetric = dm
			go NotifyTripRoute(&tr, &dm, ns)
		}
		if tr.Status == "not_started" && offset > 0 {
			startLoc := tr.ScheduledStartLocation
			endLoc := tr.ScheduledEndLocation
			dm, err := ds.GetDuration(startLoc, endLoc, time.Now().Add(offset))
			if err != nil {
				return err
			}
			offset += dm.Duration
			lastDurationMetric = dm
			go NotifyTripRoute(&tr, &dm, ns)
		}
	}
	for _, tr := range trsToBeNotified {
		go NotifyTripRoute(&tr, &lastDurationMetric, ns)
	}
	return nil
}

func handleCheckoutTrip(trip *db.Trip, currentLocation db.Location) error {
	ds := GetDurationService()
	ns := GetNotificationService()
	tripNotStarted := true
	for _, tr := range trip.TripRoutes {
		if tr.Status != "not_started" {
			tripNotStarted = false
			break
		}
	}
	if tripNotStarted {
		startLoc := currentLocation
		// This needs to be only for the first employee
		endLoc := trip.TripRoutes[0].ScheduledStartLocation
		dm, err := ds.GetDuration(startLoc, endLoc, time.Now())
		if err != nil {
			return err
		}
		for _, tr := range trip.TripRoutes {
			go NotifyTripRoute(&tr, &dm, ns)
		}
		return nil
	}

	for _, tr := range trip.TripRoutes {
		var offset time.Duration
		if tr.Status == "on_board" && offset == 0 {
			startLoc := currentLocation
			endLoc := tr.ScheduledEndLocation
			dm, err := ds.GetDuration(startLoc, endLoc, time.Now())
			if err != nil {
				return err
			}
			offset += dm.Duration
			go NotifyTripRoute(&tr, &dm, ns)
		}

		if tr.Status == "on_board" && offset > 0 {
			startLoc := tr.ScheduledStartLocation
			endLoc := tr.ScheduledEndLocation
			dm, err := ds.GetDuration(startLoc, endLoc, time.Now().Add(offset))
			if err != nil {
				return err
			}
			offset += dm.Duration
			go NotifyTripRoute(&tr, &dm, ns)
		}
	}
	return nil
}

func GetETAForTrip(q sqlx.Queryer, trip *db.Trip) error {
	if err := trip.LoadTripRoutes(q, false); err != nil {
		return err
	}
	tl, err := db.LatestTripLocation(q, trip.ID)
	if err != nil {
		return err
	}
	currentLocation := tl.Location
	if trip.TripType == db.TripTypeCheckIn {
		return handleCheckinTrip(trip, currentLocation)
	}
	return handleCheckoutTrip(trip, currentLocation)
}

func StartETAServiceTimer(cancelChan chan bool) {
	ticker := time.NewTicker(2 * time.Minute)
	for {
		select {
		case _ = <-ticker.C:
			activeTrips, err := db.GetTripsByStatus(db.CurrentDB(), "active")
			if err != nil {
				log.Errorf("unable to get active trips - %s", err)
			}
			for _, t := range activeTrips {
				go func(t *db.Trip) {
					err := GetETAForTrip(db.CurrentDB(), t)
					if err != nil {
						log.Errorf("Error processing ETA for Trip %d", t.ID)
						log.Error(err)
					}
				}(t)
			}
		case _ = <-cancelChan:
			break
		}
	}
}

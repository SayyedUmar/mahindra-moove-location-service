package services

import (
	"time"

	"github.com/MOOVE-Network/location_service/db"
	log "github.com/sirupsen/logrus"
)

func handleCheckinTrip(trip *db.Trip, currentLocation db.Location, clock Clock) error {
	ds := GetDurationService()
	ns := GetNotificationService()
	// TODO: Verify if driver_arrived should be one of the all checked in statuses
	log.Info(trip)
	if trip.AllCheckedIn() {
		endLocation := trip.TripRoutes[len(trip.TripRoutes)-1].ScheduledEndLocation
		log.Infof("Requesting eta of trip %d from %s to %s with offset of %d mins\n", trip.ID, currentLocation.ToString(), endLocation.ToString(), 0)
		dm, err := ds.GetDuration(currentLocation, endLocation, clock.Now())
		if err != nil {
			return err
		}
		for _, tr := range trip.TripRoutes {
			NotifyTripRoute(&tr, &dm, 0, ns)
		}
		return nil
	}
	var offset time.Duration
	var trsToBeNotified []db.TripRoute
	var previousEndLocation db.Location
	for _, tr := range trip.TripRoutes {
		if tr.Status == "on_board" {
			// Notify them the last
			trsToBeNotified = append(trsToBeNotified, tr)
		}
		if tr.Status == "not_started" && offset == 0 {
			startLoc := currentLocation
			endLoc := tr.ScheduledStartLocation
			previousEndLocation = tr.ScheduledStartLocation
			log.Infof("Requesting eta of trip %d from %s to %s with offset of %d mins\n", trip.ID, startLoc.ToString(), endLoc.ToString(), 0)
			dm, err := ds.GetDuration(startLoc, endLoc, clock.Now())
			if err != nil {
				return err
			}
			NotifyTripRoute(&tr, &dm, offset, ns)
			offset += dm.Duration
		} else if tr.Status == "not_started" && offset > 0 {
			startLoc := previousEndLocation
			endLoc := tr.ScheduledStartLocation
			log.Infof("Requesting eta of trip %d from %s to %s with offset of %d mins\n", trip.ID, startLoc.ToString(), endLoc.ToString(), int64(offset.Minutes()))
			dm, err := ds.GetDuration(startLoc, endLoc, clock.Now().Add(offset))
			if err != nil {
				return err
			}
			NotifyTripRoute(&tr, &dm, offset, ns)
			offset += dm.Duration
		}
	}
	if len(trsToBeNotified) > 0 {
		lastTr := trip.TripRoutes[len(trip.TripRoutes)-1]
		startLoc := lastTr.ScheduledStartLocation
		endLoc := lastTr.ScheduledEndLocation
		dm, err := ds.GetDuration(startLoc, endLoc, clock.Now().Add(offset))
		if err != nil {
			return err
		}
		for _, tr := range trsToBeNotified {
			NotifyTripRoute(&tr, &dm, offset, ns)
		}
	}
	return nil
}

func handleCheckoutTrip(trip *db.Trip, currentLocation db.Location, clock Clock) error {
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
		log.Infof("Requesting eta of trip %d from %s to %s with offset of %d mins\n", trip.ID, startLoc.ToString(), endLoc.ToString(), 0)
		dm, err := ds.GetDuration(startLoc, endLoc, clock.Now())
		if err != nil {
			return err
		}
		for _, tr := range trip.TripRoutes {
			go NotifyTripRoute(&tr, &dm, 0, ns)
		}
		return nil
	}

	for _, tr := range trip.TripRoutes {
		var offset time.Duration
		if tr.Status == "on_board" && offset == 0 {
			startLoc := currentLocation
			endLoc := tr.ScheduledEndLocation
			log.Infof("Requesting eta of trip %d from %s to %s with offset of %d mins\n", trip.ID, startLoc.ToString(), endLoc.ToString(), 0)
			dm, err := ds.GetDuration(startLoc, endLoc, clock.Now())
			if err != nil {
				return err
			}
			offset += dm.Duration
			go NotifyTripRoute(&tr, &dm, offset, ns)
		}

		if tr.Status == "on_board" && offset > 0 {
			startLoc := tr.ScheduledStartLocation
			endLoc := tr.ScheduledEndLocation
			log.Infof("Requesting eta of trip %d from %s to %s with offset of %d mins\n", trip.ID, startLoc.ToString(), endLoc.ToString(), int(offset.Minutes()))
			dm, err := ds.GetDuration(startLoc, endLoc, clock.Now().Add(offset))
			if err != nil {
				return err
			}
			offset += dm.Duration
			go NotifyTripRoute(&tr, &dm, offset, ns)
		}
	}
	return nil
}

type Clock interface {
	Now() time.Time
}
type realClock struct {
}

func (rc realClock) Now() time.Time {
	return time.Now()
}

var clock = realClock{}

func GetETAForTrip(trip *db.Trip, currentLocation db.Location, clock Clock) error {
	if trip.TripType == db.TripTypeCheckIn {
		return handleCheckinTrip(trip, currentLocation, clock)
	}
	return handleCheckoutTrip(trip, currentLocation, clock)
}

func GetETAForActiveTrips() {
	activeTrips, err := db.GetTripsByStatus(db.CurrentDB(), "active")
	if err != nil {
		log.Errorf("unable to get active trips - %s", err)
	}
	for _, t := range activeTrips {
		go func(t *db.Trip) {
			tl, err := db.LatestTripLocation(db.CurrentDB(), t.ID)
			if err != nil {
				log.Errorf("Error getting current location for Trip %d", t.ID)
				log.Error(err)
			}
			err = GetETAForTrip(t, tl.Location, clock)
			if err != nil {
				log.Errorf("Error processing ETA for Trip %d", t.ID)
				log.Error(err)
			}
		}(t)
	}
}

func StartETAServiceTimer(cancelChan chan bool) {
	GetETAForActiveTrips()
	ticker := time.NewTicker(5 * time.Minute)
	for {
		select {
		case _ = <-ticker.C:
			GetETAForActiveTrips()
		case _ = <-cancelChan:
			break
		}
	}
}

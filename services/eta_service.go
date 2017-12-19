package services

import (
	"time"

	"github.com/MOOVE-Network/location_service/db"
	"github.com/jmoiron/sqlx"
)

func handleCheckinTrip(trip *db.Trip, currentLocation db.Location) error {
	ds := GetDurationService()
	if trip.AllCheckedIn() {
		endLocation := trip.TripRoutes[len(trip.TripRoutes)-1].ScheduledEndLocation
		_, err := ds.GetDuration(currentLocation, endLocation, time.Now())
		if err != nil {
			return err
		}
		// notify all trip routes
		return nil
	}
	var toSiteDuration time.Duration
	var offset time.Duration
	for _, tr := range trip.TripRoutes {
		if tr.IsOnBoard() {
			if toSiteDuration == 0 {
				endLocation := trip.TripRoutes[len(trip.TripRoutes)-1].ScheduledEndLocation
				dm, err := ds.GetDuration(currentLocation, endLocation, time.Now())
				if err != nil {
					return err
				}
				toSiteDuration = dm.Duration
			}
			// notify tr with toSiteDuration
		}
		if tr.Status == "not_started" && offset == 0 {
			startLoc := currentLocation
			endLoc := tr.ScheduledStartLocation
			dm, err := ds.GetDuration(startLoc, endLoc, time.Now())
			if err != nil {
				return err
			}
			offset = dm.Duration
			// notify tr with dur
		}
		if tr.Status == "not_started" && offset > 0 {
			startLoc := tr.ScheduledStartLocation
			endLoc := tr.ScheduledEndLocation
			dm, err := ds.GetDuration(startLoc, endLoc, time.Now().Add(offset))
			if err != nil {
				return err
			}
			offset += dm.Duration
			// notify
		}
	}
	return nil
}

func handleCheckoutTrip(trip *db.Trip, currentLocation db.Location) error {
	ds := GetDurationService()
	var dmToSite DurationMetrics
	for _, tr := range trip.TripRoutes {
		if tr.Status == "not_started" || tr.Status == "driver_arrived" {
			if dmToSite.Duration == 0 {
				startLoc := currentLocation
				endLoc := tr.ScheduledStartLocation
				dm, err := ds.GetDuration(startLoc, endLoc, time.Now())
				if err != nil {
					return err
				}
				dmToSite = dm
			}
			// notify trip_route with dm
		}
		var offset time.Duration
		if tr.Status == "on_board" && offset == 0 {
			startLoc := currentLocation
			endLoc := tr.ScheduledEndLocation
			dm, err := ds.GetDuration(startLoc, endLoc, time.Now())
			if err != nil {
				return err
			}
			offset += dm.Duration
			// notify trip_route with dm
		}
		if tr.Status == "on_board" && offset > 0 {
			startLoc := tr.ScheduledStartLocation
			endLoc := tr.ScheduledEndLocation
			dm, err := ds.GetDuration(startLoc, endLoc, time.Now().Add(offset))
			if err != nil {
				return err
			}
			offset += dm.Duration
			// notify trip_route with dm
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

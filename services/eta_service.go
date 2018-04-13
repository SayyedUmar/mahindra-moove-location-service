package services

import (
	"fmt"
	"strconv"
	"time"

	"github.com/MOOVE-Network/location_service/db"
	log "github.com/sirupsen/logrus"
	null "gopkg.in/guregu/null.v3"
)

func NotNullTime(t time.Time) null.Time {
	return null.NewTime(t, true)
}

type ETATripRoute struct {
	ID             int       `json:"id"`
	PickupTime     null.Time `json:"pickup_time"`
	DropoffTime    null.Time `json:"dropoff_time"`
	ETAInMinutes   float64   `json:"eta_in_minutes"`
	EmployeeUserID int       `json:"employee_user_id"`
	Status         string    `json:"status"`
}
type ETAResponse struct {
	ID         int            `json:"id"`
	UpdatedAt  time.Time      `json:"updated_at"`
	TripRoutes []ETATripRoute `json:"trip_routes"`
}

func handleCheckinTrip(trip *db.Trip, currentLocation db.Location, clock Clock) (*ETAResponse, error) {
	etaBusMap := make(map[int]DurationMetrics)
	etaResp := ETAResponse{ID: trip.ID, UpdatedAt: clock.Now()}
	ds := GetDurationService()
	ns := GetNotificationService()
	// TODO: Verify if driver_arrived should be one of the all checked in statuses
	log.Info(trip)
	if trip.AllCheckedIn() {
		endLocation := trip.TripRoutes[len(trip.TripRoutes)-1].ScheduledEndLocation
		log.Infof("Requesting eta of trip %d from %s to %s with offset of %d mins\n", trip.ID, currentLocation.ToString(), endLocation.ToString(), 0)
		dm, err := ds.GetDuration(currentLocation, endLocation, clock.Now())
		if err != nil {
			return nil, err
		}
		for _, tr := range trip.TripRoutes {
			etaResp.TripRoutes = append(etaResp.TripRoutes, ETATripRoute{
				ID:             tr.ID,
				PickupTime:     NotNullTime(clock.Now().Add(dm.Duration)),
				ETAInMinutes:   dm.Duration.Minutes(),
				EmployeeUserID: tr.EmployeeUserID,
				Status:         tr.Status,
			})
			NotifyTripRouteToEmployee(&tr, &dm, 0, ns)
		}
		NotifyTripRouteToDriver(&trip.TripRoutes[0], &dm, 0, ns)
		return &etaResp, nil
	}
	var offset time.Duration
	var trsToBeNotified []db.TripRoute
	var previousEndLocation db.Location
	for _, tr := range trip.TripRoutes {
		if tr.Status != "on_board" && tr.Status != "not_started" {
			etaResp.TripRoutes = append(etaResp.TripRoutes, ETATripRoute{
				ID:             tr.ID,
				EmployeeUserID: tr.EmployeeUserID,
				Status:         tr.Status,
			})
		}
		if tr.Status == "on_board" {
			// Notify them the last
			trsToBeNotified = append(trsToBeNotified, tr)
		}
		if tr.Status == "not_started" && offset == 0 {
			startLoc := currentLocation
			endLoc := tr.ScheduledStartLocation
			previousEndLocation = tr.ScheduledStartLocation
			dm, ok := etaBusMap[tr.ScheduledRouteOrder]
			if !ok {
				log.Infof("Requesting eta of trip %d from %s to %s with offset of %d mins\n", trip.ID, startLoc.ToString(), endLoc.ToString(), 0)
				dmLocal, err := ds.GetDuration(startLoc, endLoc, clock.Now())
				if err != nil {
					return nil, err
				}
				etaBusMap[tr.ScheduledRouteOrder] = dmLocal
				dm = dmLocal
			}
			etaResp.TripRoutes = append(etaResp.TripRoutes, ETATripRoute{
				ID:             tr.ID,
				PickupTime:     NotNullTime(clock.Now().Add(dm.Duration)),
				ETAInMinutes:   dm.Duration.Minutes(),
				EmployeeUserID: tr.EmployeeUserID,
				Status:         tr.Status,
			})
			NotifyTripRouteToEmployee(&tr, &dm, offset, ns)
			NotifyTripRouteToDriver(&tr, &dm, offset, ns)
			offset += dm.Duration
		} else if tr.Status == "not_started" && offset > 0 {
			startLoc := previousEndLocation
			endLoc := tr.ScheduledStartLocation
			dm, ok := etaBusMap[tr.ScheduledRouteOrder]
			if !ok {
				log.Infof("Requesting eta of trip %d from %s to %s with offset of %d mins\n", trip.ID, startLoc.ToString(), endLoc.ToString(), int64(offset.Minutes()))
				dmLocal, err := ds.GetDuration(startLoc, endLoc, clock.Now().Add(offset))
				if err != nil {
					return nil, err
				}
				etaBusMap[tr.ScheduledRouteOrder] = dmLocal
				dm = dmLocal
			}
			etaResp.TripRoutes = append(etaResp.TripRoutes, ETATripRoute{
				ID:             tr.ID,
				PickupTime:     NotNullTime(clock.Now().Add(dm.Duration).Add(offset)),
				ETAInMinutes:   (dm.Duration + offset).Minutes(),
				EmployeeUserID: tr.EmployeeUserID,
				Status:         tr.Status,
			})
			NotifyTripRouteToEmployee(&tr, &dm, offset, ns)
			offset += dm.Duration
		}
	}
	if len(trsToBeNotified) > 0 {
		lastTr := trip.TripRoutes[len(trip.TripRoutes)-1]
		startLoc := lastTr.ScheduledStartLocation
		endLoc := lastTr.ScheduledEndLocation
		dm, err := ds.GetDuration(startLoc, endLoc, clock.Now().Add(offset))
		if err != nil {
			return nil, err
		}
		for _, tr := range trsToBeNotified {
			etaResp.TripRoutes = append(etaResp.TripRoutes, ETATripRoute{
				ID:             tr.ID,
				DropoffTime:    NotNullTime(clock.Now().Add(dm.Duration).Add(offset)),
				ETAInMinutes:   (dm.Duration + offset).Minutes(),
				EmployeeUserID: tr.EmployeeUserID,
				Status:         tr.Status,
			})
			NotifyTripRouteToEmployee(&tr, &dm, offset, ns)
		}
	}
	return &etaResp, nil
}

func handleCheckoutTrip(trip *db.Trip, currentLocation db.Location, clock Clock) (*ETAResponse, error) {
	etaBusMap := make(map[int]DurationMetrics)
	etaResp := ETAResponse{ID: trip.ID, UpdatedAt: clock.Now()}
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
			return nil, err
		}
		for _, tr := range trip.TripRoutes {
			etaResp.TripRoutes = append(etaResp.TripRoutes, ETATripRoute{
				ID:             tr.ID,
				PickupTime:     NotNullTime(clock.Now().Add(dm.Duration)),
				ETAInMinutes:   dm.Duration.Minutes(),
				EmployeeUserID: tr.EmployeeUserID,
				Status:         tr.Status,
			})
			NotifyTripRouteToEmployee(&tr, &dm, 0, ns)
		}
		return &etaResp, nil
	}

	for _, tr := range trip.TripRoutes {
		var offset time.Duration
		if tr.Status != "on_board" && tr.Status != "not_started" {
			etaResp.TripRoutes = append(etaResp.TripRoutes, ETATripRoute{
				ID:             tr.ID,
				EmployeeUserID: tr.EmployeeUserID,
				Status:         tr.Status,
			})
		}
		if tr.Status == "on_board" && offset == 0 {
			startLoc := currentLocation
			endLoc := tr.ScheduledEndLocation
			dm, ok := etaBusMap[tr.ScheduledRouteOrder]
			if !ok {
				log.Infof("Requesting eta of trip %d from %s to %s with offset of %d mins\n", trip.ID, startLoc.ToString(), endLoc.ToString(), 0)
				dmLocal, err := ds.GetDuration(startLoc, endLoc, clock.Now())
				if err != nil {
					return nil, err
				}
				etaBusMap[tr.ScheduledRouteOrder] = dmLocal
				dm = dmLocal
			}
			etaResp.TripRoutes = append(etaResp.TripRoutes, ETATripRoute{
				ID:             tr.ID,
				DropoffTime:    NotNullTime(clock.Now().Add(dm.Duration).Add(offset)),
				ETAInMinutes:   (dm.Duration + offset).Minutes(),
				EmployeeUserID: tr.EmployeeUserID,
				Status:         tr.Status,
			})
			NotifyTripRouteToEmployee(&tr, &dm, offset, ns)
			offset += dm.Duration
		} else if tr.Status == "on_board" && offset > 0 {
			startLoc := tr.ScheduledStartLocation
			endLoc := tr.ScheduledEndLocation
			dm, ok := etaBusMap[tr.ScheduledRouteOrder]
			if !ok {
				log.Infof("Requesting eta of trip %d from %s to %s with offset of %d mins\n", trip.ID, startLoc.ToString(), endLoc.ToString(), int(offset.Minutes()))
				dmLocal, err := ds.GetDuration(startLoc, endLoc, clock.Now().Add(offset))
				if err != nil {
					return nil, err
				}
				etaBusMap[tr.ScheduledRouteOrder] = dmLocal
				dm = dmLocal
			}
			etaResp.TripRoutes = append(etaResp.TripRoutes, ETATripRoute{
				ID:             tr.ID,
				DropoffTime:    NotNullTime(clock.Now().Add(dm.Duration).Add(offset)),
				ETAInMinutes:   (dm.Duration + offset).Minutes(),
				EmployeeUserID: tr.EmployeeUserID,
				Status:         tr.Status,
			})
			NotifyTripRouteToEmployee(&tr, &dm, offset, ns)
			offset += dm.Duration
		}
	}
	return &etaResp, nil
}

type Clock interface {
	Now() time.Time
}
type RealClock struct {
}

func (rc RealClock) Now() time.Time {
	return time.Now()
}

var clock = RealClock{}

func GetETAForTrip(trip *db.Trip, currentLocation db.Location, clock Clock) (*ETAResponse, error) {
	if len(trip.TripRoutes) < 1 {
		return nil, fmt.Errorf("The trip %d has no trip routes", trip.ID)
	}
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
			etas, err := GetETAForTrip(t, tl.Location, clock)
			if err != nil {
				log.Errorf("Error processing ETA for Trip %d", t.ID)
				log.Error(err)
			}
			for _, tr := range etas.TripRoutes {
				db.SaveEta(db.CurrentDB(), tr.ID, tr.PickupTime, tr.DropoffTime)
			}
		}(t)
	}
}

func GetETAForAssignedTrip() {
	assignedTrips, err := db.GetTripsByStatus(db.CurrentDB(), "assigned")
	if err != nil {
		log.Errorf("unable to get assigned trips - %s", err)
	}
	for _, t := range assignedTrips {
		go func(t *db.Trip) {
			if len(t.TripRoutes) == 0 {
				log.Infof("zero trip routes in trip - %d", t.ID)
				return
			}
			if !t.ScheduledStartDate.Valid {
				log.Errorf("No scheduled start time for trip - %d, can't notify driver about when to start trip", t.DriverID)
				log.Error(err)
				return
			}

			lastLocation, err := db.DriverLocation(db.CurrentDB(), t.DriverUserID)
			if err != nil {
				log.Errorf("Error getting last locaction for driver - %d", t.DriverID)
				log.Error(err)
				return
			}
			ds := GetDurationService()
			dm, err := ds.GetDuration(*lastLocation, t.TripRoutes[0].ScheduledStartLocation, clock.Now())
			if err != nil {
				log.Errorf("Error getting duration for trip - %d with start location as %v and stop location as %v", t.ID, lastLocation, t.TripRoutes[0].ScheduledStartLocation)
				log.Error(err)
				return
			}

			if clock.Now().Add(dm.Duration).Add(time.Duration(time.Minute * 20)).After(t.ScheduledStartDate.Time) {
				ns := GetNotificationService()
				data := make(map[string]interface{})
				data["push_type"] = "driver_should_start_trip"
				data["trip_id"] = t.ID
				err := ns.SendNotification(strconv.Itoa(t.DriverUserID), data, "driver")
				if err != nil {
					log.Errorf("Error while sending notification to start trip - %d", t.ID)
					log.Error(err)
					return
				}
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
			GetETAForAssignedTrip()
		case _ = <-cancelChan:
			break
		}
	}
}

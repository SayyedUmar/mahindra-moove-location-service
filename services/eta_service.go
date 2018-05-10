package services

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/MOOVE-Network/location_service/utils"

	"github.com/MOOVE-Network/location_service/db"
	log "github.com/sirupsen/logrus"
	null "gopkg.in/guregu/null.v3"
)

var tripShouldStartNotifiers map[int]*time.Timer

func init() {
	tripShouldStartNotifiers = make(map[int]*time.Timer)
}

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
	activeTrips, err := db.GetTripsByStatuses(db.CurrentDB(), "active")
	if err != nil {
		log.Errorf("unable to get active trips - %s", err)
		return
	}
	for _, t := range activeTrips {
		go func(t *db.Trip) {
			tl, err := db.LatestTripLocation(db.CurrentDB(), t.ID)
			if err != nil {
				log.Errorf("Error getting current location for Trip %d", t.ID)
				log.Error(err)
				return
			}
			etas, err := GetETAForTrip(t, tl.Location, clock)
			if err != nil {
				log.Errorf("Error processing ETA for Trip %d", t.ID)
				log.Error(err)
				return
			}
			for _, tr := range etas.TripRoutes {
				db.SaveEta(db.CurrentDB(), tr.ID, tr.PickupTime, tr.DropoffTime)
			}
		}(t)
	}
}

func getETAForAssignedTrip() {
	assignedTrips, err := db.GetTripsByStatuses(db.CurrentDB(), "assigned", "assign_requested", "assign_request_expired")
	if err != nil {
		log.Errorf("unable to get assigned trips - %s", err)
	}
	for _, t := range assignedTrips {
		go func(t *db.Trip) {
			log.Debugf("Checking for trip %d for eta", t.ID)
			if !t.ScheduledStartDate.Valid {
				log.Debugf("ScheduledStartDate for trip: %d is invalid", t.ID)
				return
			}

			if clock.Now().After(t.ScheduledStartDate.Time) {
				log.Debugf("Stoping eta calculation as first pickup time is already passed for trip: %d", t.ID)
				return
			}

			maxTimeToCalculateEta, err := db.GetMaxTimeToCalculateStartTripEta(db.CurrentDB())
			if err != nil {
				maxTimeToCalculateEta = 60 * 3 //Assigning default value of 3 hours
			}

			if clock.Now().Before(t.ScheduledStartDate.Time.Add(-(time.Duration(maxTimeToCalculateEta) * time.Minute))) {
				log.Debugf("ScheduledStartDate for trip %d is more than 3 hours ahead of current time: %s", t.ID, clock.Now().String())
				return
			}

			lastDriverLocation, err := db.DriverLocation(db.CurrentDB(), t.DriverUserID)
			if err != nil {
				log.Errorf("Error getting last location for driver - %d", t.DriverID)
				log.Error(err)
				return
			}

			lastLocation, err := t.GetDriverShouldStartTripLocation(db.CurrentDB())
			if err == nil {
				minDistanceToCalculateEta, err := db.GetMinDistanceToCalculateStartTripEta(db.CurrentDB())
				if err != nil {
					minDistanceToCalculateEta = 500 //Assigning default value of 500 meters
				}
				distance := utils.Distance(lastLocation.Lat, lastLocation.Lng, lastDriverLocation.Lat, lastDriverLocation.Lng)
				if distance < float64(minDistanceToCalculateEta) {
					log.Debugf("distance between driver location and driver_should_start_trip_location is less than %d for trip: %d", minDistanceToCalculateEta, t.ID)
					return
				}
			}

			newStartTime, err := FindWhenShouldDriverStartTrip(t, lastDriverLocation, clock)
			if err != nil {
				log.Errorf("Error: [%s] while find eta to start trip : %d", err.Error(), t.ID)
				log.Error(err)
				return
			}

			calculationTime := clock.Now()
			err = t.UpdateDriverShouldStartTripTimeAndLocation(db.CurrentDB(), *newStartTime, *lastDriverLocation, calculationTime)
			if err != nil {
				log.Errorf("Error updating driver should start info in trips table for trip - %d with time as %v and location as %v", t.ID, newStartTime, lastDriverLocation)
				log.Error(err)
				return
			}

			_, err = NotifyDriverShouldStartTrip(t, newStartTime, &calculationTime)
			if err != nil {
				log.Errorf("Error while sending notification to start trip: %d\n", t.ID)
				log.Error(err)
			}

			SetStartTripDelayTimer(t.ID, newStartTime)
		}(t)
	}
}

func SetStartTripDelayTimer(tripID int, startTime *time.Time) {
	durationForDelayedTripNotification, err := db.GetBufferDurationForDelayTripNotification(db.CurrentDB())
	if err != nil {
		durationForDelayedTripNotification = 10 //defaulting of 10 minutes
	}

	timerDuration := startTime.Add(time.Duration(durationForDelayedTripNotification) * time.Minute).Sub(time.Now())

	newTimer := time.AfterFunc(timerDuration, func() {
		trip, err := db.GetTripByID(db.CurrentDB(), tripID)
		if err != nil {
			log.Errorf("could not get trip for id %d to create trip should start notification.")
		}
		if trip.HasStarted() {
			log.Infof("Trip %d has already started", trip.ID)
			return
		}
		tx, err := db.CurrentDB().Beginx()
		if err != nil {
			log.Errorf("Could not create TripShouldStart notification for trip id %d because %v", trip.ID, err)
			return
		}
		defer tx.Commit()
		tss, err := db.CreateTripShouldStartNotification(tx, trip.ID, trip.DriverID)
		if err != nil {
			log.Errorf("Could not create TripShouldStart notification for trip id %d because %v", trip.ID, err)
			return
		}
		log.Infof("Created a trip should start notification for trip %d", tss.TripID)
	})

	oldTimer, ok := tripShouldStartNotifiers[tripID]
	if ok {
		oldTimer.Stop()
	}
	tripShouldStartNotifiers[tripID] = newTimer
}

func FindWhenShouldDriverStartTrip(trip *db.Trip, driverLocation *db.Location, clock Clock) (*time.Time, error) {
	if len(trip.TripRoutes) == 0 {
		return nil, errors.New("can not find start time for trip with zero trip routes")
	}
	if !trip.ScheduledStartDate.Valid {
		return nil, errors.New("ScheduledStartDate for trip is invalid")
	}

	ds := GetDurationService()
	dm, err := ds.GetDuration(*driverLocation, trip.TripRoutes[0].ScheduledStartLocation, clock.Now())
	if err != nil {
		log.Errorf("Error getting duration for trip - %d with start location as %v and stop location as %v", trip.ID, driverLocation, trip.TripRoutes[0].ScheduledStartLocation)
		return nil, err
	}

	newStartTime := trip.ScheduledStartDate.Time.Add(-dm.Duration)
	return &newStartTime, nil
}

func NotifyDriverShouldStartTrip(trip *db.Trip, newStartTime *time.Time, calculationTime *time.Time) (bool, error) {
	ns := GetNotificationService()
	data := make(map[string]interface{})
	data["push_type"] = "driver_should_start_trip"
	data["trip_id"] = trip.ID
	data["driver_should_start_trip_time"] = newStartTime.Unix()
	data["driver_should_start_trip_timestamp"] = calculationTime.Unix()
	data["driver_should_start_trip_calc_time"] = calculationTime.Unix() //For backword compatibility

	log.Debugf("Sending start trip notification to driver : %d, with notification payload: %v", trip.DriverUserID, data)

	err := ns.SendNotification(strconv.Itoa(trip.DriverUserID), data, "user")
	if err != nil {
		return false, err
	}
	return true, nil
}

func StartETAServiceTimer(cancelChan chan bool) {
	// GetETAForActiveTrips()
	getETAForAssignedTrip()
	ticker := time.NewTicker(5 * time.Minute)
	for {
		select {
		case _ = <-ticker.C:
			// GetETAForActiveTrips()
			getETAForAssignedTrip()
		case _ = <-cancelChan:
			break
		}
	}
}

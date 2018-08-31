package web

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/MOOVE-Network/location_service/db"
	"github.com/MOOVE-Network/location_service/identity"
	"github.com/MOOVE-Network/location_service/models"
	"github.com/MOOVE-Network/location_service/services"
	"github.com/MOOVE-Network/location_service/socketstore"
	log "github.com/sirupsen/logrus"
)

var hbMutex = &sync.Mutex{}
var heartBeats = make(map[int]*db.HeartBeat)
var tlMutex = &sync.Mutex{}
var tripLocations []db.TripLocation
var gfMutex = &sync.Mutex{}
var gfEvents []socketstore.GeofenceEvent

var dlMutex = &sync.Mutex{}
var driverLocations []models.DriverLocation

var driverLocationsForSpeedCheck []models.DriverLocation
var overSpeedDriversTime = make(map[int64]*time.Time)
var overSpeedNotificationsSent = make(map[int64]bool)
var overSpeedMutex = &sync.Mutex{}

func WriteHeartBeat(w http.ResponseWriter, req *http.Request) {
	ident := req.Context().Value("identity").(*identity.Identity)
	hb := &db.HeartBeat{UserID: ident.Id, UpdatedAt: time.Now()}
	latStr := req.URL.Query().Get("lat")
	lngStr := req.URL.Query().Get("lng")
	if latStr == "" || lngStr == "" {
		writeOk(w)
		return
	}
	lat, errLat := strconv.ParseFloat(latStr, 64)
	lng, errLng := strconv.ParseFloat(lngStr, 64)
	if errLat != nil || errLng != nil {
		writeOk(w)
		return
	}
	hb.Lat = lat
	hb.Lng = lng
	heartBeats[ident.Id] = hb
	err := hb.Save(db.CurrentDB())
	if err != nil {
		log.Error(err)
	}
	writeOk(w)
}

func setupHeartBeatTimer() {
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for _ = range ticker.C {
			var hbs []*db.HeartBeat
			hbMutex.Lock()
			for _, hb := range heartBeats {
				hbs = append(hbs, hb)
			}
			hbMutex.Unlock()
			tx := db.CurrentDB().MustBegin()
			for _, heartBeat := range hbs {
				err := heartBeat.Save(tx)
				if err != nil {
					log.Warn("Unable to save heartbeat", err)
				}
			}
			err := tx.Commit()
			if err != nil {
				log.Panic(err)
			}
		}
	}()
}

func setupTripLocationTimer() {
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		stmt := db.InsertTripLocationStatement(db.CurrentDB())
		for _ = range ticker.C {
			var tls []db.TripLocation
			tlMutex.Lock()
			for _, tl := range tripLocations {
				tls = append(tls, tl)
			}
			tripLocations = nil
			tlMutex.Unlock()
			tx := db.CurrentDB().MustBegin()
			for _, tripLocation := range tls {
				err := tripLocation.Save(stmt)
				if err != nil {
					log.Error("Unable to save trip location")
					log.Error(err)
				}
			}
			err := tx.Commit()
			if err != nil {
				log.Error("Unable to save trip location")
				log.Error(err)
			}
		}
	}()
}

func setupDriverLocationTimer() {
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		stmt, err := models.DriverLocationPrepareInsertStmt(models.CurrentDB())
		if err != nil {
			log.Errorf("Unable to prepare insert statement for driver location %s ", err)
			panic(err)
		}
		for _ = range ticker.C {
			var dls []models.DriverLocation
			dlMutex.Lock()
			for _, dl := range driverLocations {
				dls = append(dls, dl)
			}
			driverLocations = nil
			dlMutex.Unlock()
			for _, driverLocation := range dls {
				err := driverLocation.Insert(stmt)
				if err != nil {
					log.Errorf("Unable to save DriverLocation location - %s", err)
				}
			}
		}
	}()
}

func setupGeofenceTimer() {
	log.Infoln("started geofence timer")
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for _ = range ticker.C {
			var tempGfEvents []socketstore.GeofenceEvent
			tlMutex.Lock()
			tempGfEvents = append(tempGfEvents, gfEvents...)
			gfEvents = nil
			tlMutex.Unlock()
			processGeofenceEvents(tempGfEvents)
		}
	}()
}

func processGeofenceEvents(gfEvents []socketstore.GeofenceEvent) {
	//Map which can hold list of geofence events per tripId.
	//This is for db query optimization.
	//With this we will fetch trip only one for list of geofence events.
	gfEventMap := make(map[int][]socketstore.GeofenceEvent)

	//grouping geofence events based on tripId
	for _, geofenceEvent := range gfEvents {
		tripID := geofenceEvent.TripID
		if tripID > 0 {
			val, ok := gfEventMap[tripID]
			if !ok {
				//giving capacity as 2 because we get two events per geofence.
				val = make([]socketstore.GeofenceEvent, 0, 2)
			}
			gfEventMap[tripID] = append(val, geofenceEvent)
		}
	}

	for tripID, events := range gfEventMap {
		trip, err := db.GetTripByID(db.CurrentDB(), tripID)
		if err != nil {
			log.Errorf("Unable to get trip for tripId: %d", tripID)
			log.Error(err)
			continue
		}
		log.Infof("for trip Id %d, found trip %+v", tripID, trip)
		sendDriverArrivingNotification(trip, events)
		updateGeofenceInfoInTripRoutes(trip, events)
	}
}

func sendDriverArrivingNotification(trip *db.Trip, events []socketstore.GeofenceEvent) {
	for _, gfEvent := range events {
		//We send Notification only for Dwell event and for Wider Geofence.
		if gfEvent.IsDwellEvent() && gfEvent.IsForWiderGeofence() {
			driver, err := db.GetDriverByID(db.CurrentDB(), trip.DriverID)
			if err != nil {
				log.Errorf("Unable to get driver for tripId: %d", trip.ID)
				log.Error(err)
				continue
			}
			log.Infof("for driver Id %d, found driver %+v", trip.DriverID, driver)
			if gfEvent.IsForSite() {
				log.Info("Checking tripRoute for site")
				for _, tripRoute := range trip.TripRoutes {
					if tripRoute.IsNotStarted() {
						services.SendDriverArrivingNotification(trip.ID, tripRoute.EmployeeUserID, driver)
					}
				}
			} else if gfEvent.IsForNodalPoint() {
				log.Info("Checking tripRoute for nodal point")
				for _, tripRoute := range trip.TripRoutes {
					for _, tripRouteID := range gfEvent.TripRouteIDs {
						if tripRoute.ID == tripRouteID && tripRoute.IsNotStarted() {
							services.SendDriverArrivingNotification(trip.ID, tripRoute.EmployeeUserID, driver)
						}
					}
				}
			}
		}
	}
}

func updateGeofenceInfoInTripRoutes(trip *db.Trip, events []socketstore.GeofenceEvent) {
	tx := db.CurrentDB().MustBegin()
	defer func() {
		err := tx.Commit()
		if err != nil {
			log.Panic(err)
		}
	}()
	for _, gfEvent := range events {
		if gfEvent.IsDwellEvent() && gfEvent.IsForNarrowGeofence() {
			if gfEvent.IsForSite() {
				for _, tripRoute := range trip.TripRoutes {
					if trip.IsCheckIn() {
						tripRoute.UpdateCompletedGeofenceInfo(tx, gfEvent.GetLocation(), gfEvent.Timestamp)
					} else {
						tripRoute.UpdateDriverArrivedGeofenceInfo(tx, gfEvent.GetLocation(), gfEvent.Timestamp)
					}
				}
			} else if gfEvent.IsForNodalPoint() {
				for _, tripRoute := range trip.TripRoutes {
					for _, tripRouteID := range gfEvent.TripRouteIDs {
						if tripRoute.ID == tripRouteID {
							if trip.IsCheckIn() {
								tripRoute.UpdateDriverArrivedGeofenceInfo(tx, gfEvent.GetLocation(), gfEvent.Timestamp)
							} else {
								tripRoute.UpdateCompletedGeofenceInfo(tx, gfEvent.GetLocation(), gfEvent.Timestamp)
							}
						}
					}
				}
			}
		}
	}
}

func setupOverSpeedingCheckTimer() {
	log.Infoln("started over speeding check timer")
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for _ = range ticker.C {
			speedLimit, err := db.GetSpeedLimit(db.CurrentDB())
			if err != nil {
				speedLimit = 22.2222 //80Kmph.
			}
			log.Infoln("Speed Limit:", speedLimit)
			overSpeedDuration, err := db.GetSpeedLimitViolationDuration(db.CurrentDB())
			if err != nil {
				overSpeedDuration = 60 //1 minute.
			}

			tempDriverLocations := make(map[int64][]models.DriverLocation)
			overSpeedMutex.Lock()
			for _, dl := range driverLocationsForSpeedCheck {
				if dl.TripID.Valid {
					val, ok := tempDriverLocations[dl.TripID.Int64]
					if !ok {
						val = make([]models.DriverLocation, 0, 8) //Just for the sake of giving some capacity giving 8.
					}
					tempDriverLocations[dl.TripID.Int64] = append(val, dl)
				}
			}
			driverLocationsForSpeedCheck = nil
			overSpeedMutex.Unlock()
			for tripID, value := range tempDriverLocations {
				sort.Slice(value, func(i, j int) bool {
					return value[i].RecordedAt.Before(value[j].RecordedAt)
				})

				for _, dl := range value {
					overSpeedStartTime, OK := overSpeedDriversTime[tripID]
					if !OK {
						if dl.Speed <= speedLimit {
							//have not over speeded and not over speeding
							continue
						} else {
							//This if the first occurrence of over speeding
							overSpeedDriversTime[tripID] = &dl.RecordedAt
							continue
						}
					} else {
						if dl.Speed <= speedLimit {
							delete(overSpeedDriversTime, tripID)
							delete(overSpeedNotificationsSent, tripID)
							continue
						}

						if overSpeedStartTime.Add(time.Duration(overSpeedDuration) * time.Second).Before(dl.RecordedAt) {
							notificationSent, OK := overSpeedNotificationsSent[tripID]
							if OK && notificationSent {
								continue //Notification already sent, so no need to do anything
							}

							userID, err := strconv.Atoi(dl.UserID.String)
							if err != nil {
								log.Errorf("Unable to convert string userID - %s to int", dl.UserID.String)
								continue
							}
							dos, err := createOverSpeedingNotification(tripID, userID)
							if err != nil {
								log.Errorf("Could not create DriverOverSpeeding notification for trip id %d because %s", tripID, err)
								continue
							}
							overSpeedNotificationsSent[tripID] = true
							log.Infof("Created a driver over speeding notification for trip %d", dos.TripID)
						}
					}
				}
			}
		}
	}()
}

func createOverSpeedingNotification(tripID int64, userID int) (*db.Notification, error) {
	driver, err := db.GetDriverByUserID(db.CurrentDB(), userID)
	if err != nil {
		return nil, fmt.Errorf("Unable to get driver for user id - %d, error - %s", userID, err)
	}
	tx, err := db.CurrentDB().Beginx()
	if err != nil {
		return nil, fmt.Errorf("Could not create transaction in createOverSpeedingNotification func because %s", err)
	}
	defer tx.Commit()
	dos, err := db.CreateDriverOverSpeedingNotification(tx, int(tripID), driver.ID)
	if err != nil {
		return nil, err
	}
	return dos, nil
}

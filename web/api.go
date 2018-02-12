package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/MOOVE-Network/location_service/services"

	"github.com/MOOVE-Network/location_service/db"
	"github.com/MOOVE-Network/location_service/identity"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var hbMutex = &sync.Mutex{}
var heartBeats = make(map[int]*db.HeartBeat)
var tlMutex = &sync.Mutex{}
var tripLocations []db.TripLocation

func GetTripETA(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, found := vars["id"]
	if !found {
		ErrorWithMessage("Unable to find param id").Respond(w, 404)
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		ErrorWithMessage(fmt.Sprintf("Driver id is not an integer. got %s as driverID", id)).Respond(w, 422)
	}
	trip, err := db.GetTripByID(db.CurrentDB(), idInt)
	if err != nil {
		ErrorWithMessage(fmt.Sprintf("Unable to find trip. %s", err.Error())).Respond(w, 404)
	}
	tl, err := db.LatestTripLocation(db.CurrentDB(), trip.ID)
	if err != nil {
		log.Errorf("Error getting current location for Trip %d", trip.ID)
		ErrorWithMessage(fmt.Sprintf("Unable to find current location for trip %d. %s", trip.ID, err.Error())).Respond(w, 404)
	}
	resp, err := services.GetETAForTrip(trip, tl.Location, services.RealClock{})
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(resp)
	if err != nil {
		log.Error("Unable to encode eta response")
		log.Error(resp)
		ErrorWithMessage(fmt.Sprintf("Unable to encode json ETA response. %s", err.Error())).Respond(w, 500)
	}
}

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

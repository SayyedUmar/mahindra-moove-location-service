package web

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/MOOVE-Network/location_service/db"
	"github.com/MOOVE-Network/location_service/identity"
	"github.com/MOOVE-Network/location_service/models"
	log "github.com/sirupsen/logrus"
)

var hbMutex = &sync.Mutex{}
var heartBeats = make(map[int]*db.HeartBeat)
var tlMutex = &sync.Mutex{}
var tripLocations []db.TripLocation

var dlMutex = &sync.Mutex{}
var driverLocations []models.DriverLocation

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

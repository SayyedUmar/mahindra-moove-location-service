package web

import (
	"github.com/MOOVE-Network/location_service/db"
	"github.com/MOOVE-Network/location_service/identity"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

var heartBeats = make(map[int]*db.HeartBeat)

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
	writeOk(w)
}

func setupHeartBeatTimer() {
	ticker := time.NewTicker(time.Second * 1)
	go func() {
		for _ = range ticker.C {
			tx := db.CurrentDB().MustBegin()
			for _, heartBeat := range heartBeats {
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

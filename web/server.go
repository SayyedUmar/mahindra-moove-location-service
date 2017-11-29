package web

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MOOVE-Network/location_service/db"
	"github.com/MOOVE-Network/location_service/identity"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
	"time"
)

var heartBeats = make(map[int]*db.HeartBeat)

type Error struct {
	Message string `json:"error"`
}

func ErrorWithMessage(msg string) *Error {
	return &Error{Message: msg}
}

func (e *Error) Respond(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(e)
	if err != nil {
		log.Warn(err)
	}
}

func TokenAuth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		uid := req.Header.Get("uid")
		accessToken := req.Header.Get("access-token")
		client := req.Header.Get("access-token")
		if uid == "" || accessToken == "" || client == "" {
			ErrorWithMessage("Invalid Credentials").Respond(w, 401)
			return
		}
		ident := identity.FetchIdentityByUID(db.CurrentDB(), uid)
		if !ident.IsValid(client, accessToken) {
			ErrorWithMessage("Invalid Credentials").Respond(w, 401)
			return
		}
		ctx := context.WithValue(req.Context(), "identity", ident)
		newReq := req.WithContext(ctx)
		fn(w, newReq)
	}
}

func writeOk(w http.ResponseWriter) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(map[string]string{"message": "success"})
	if err != nil {
		ErrorWithMessage(err.Error()).Respond(w, 200)
	}
}

func WriteHeartBeat(w http.ResponseWriter, req *http.Request) {
	ident := req.Context().Value("identity").(identity.Identity)
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

func SetupServer() {
	port := os.Getenv("LOCATION_PORT")
	if port == "" {
		port = "4343"
	}
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/drivers/{id}/heart_beat", TokenAuth(WriteHeartBeat))
	log.Info("Starting ... ")
	log.Infof("Listening on port %s ... ", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		log.Panic(err)
	}
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
			log.Panic(err)
		}
	}()
}

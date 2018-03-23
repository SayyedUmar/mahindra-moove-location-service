package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/MOOVE-Network/location_service/redis"
	"github.com/MOOVE-Network/location_service/services"
	"github.com/MOOVE-Network/location_service/utils"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/MOOVE-Network/location_service/db"
)

func getTripID(req *http.Request) (int, error) {
	vars := mux.Vars(req)
	id, found := vars["id"]
	if !found {
		return 0, errors.New("Unable to find param id")
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Errorf("Driver id is not an integer. got %s as driverID", id)
	}
	return idInt, err
}

func getTrip(req *http.Request) (*db.Trip, error) {
	idInt, err := getTripID(req)
	if err != nil {
		return nil, fmt.Errorf("Invalid driver id in paramters %v - %s", mux.Vars(req), err)
	}
	return db.GetTripByID(db.CurrentDB(), idInt)
}
func GetTripETA(w http.ResponseWriter, req *http.Request) {
	trip, err := getTrip(req)
	if err != nil {
		log.Error("Error getting Trip %v", err)
		ErrorWithMessage(fmt.Sprintf("Unable to find trip %s", err.Error())).Respond(w, 404)
		return
	}
	tl, err := db.LatestTripLocation(db.CurrentDB(), trip.ID)
	if err != nil {
		log.Errorf("Error getting current location for Trip %d", trip.ID)
		ErrorWithMessage(fmt.Sprintf("Unable to find current location for trip %d. %s", trip.ID, err.Error())).Respond(w, 404)
		return
	}
	resp, err := services.GetETAForTrip(trip, tl.Location, services.RealClock{})
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(resp)
	if err != nil {
		log.Error("Unable to encode eta response")
		log.Error(resp)
		ErrorWithMessage(fmt.Sprintf("Unable to encode json ETA response. %s", err.Error())).Respond(w, 500)
		return
	}
}

func TripSummary(w http.ResponseWriter, r *http.Request) {
	useGoogle := false
	usePassThrough := false
	if r.URL.Query().Get("use_google") != "" {
		useGoogle = true
	}
	if r.URL.Query().Get("use_passthrough") != "" {
		usePassThrough = true
	}
	tripID, err := getTripID(r)
	if err != nil {
		ErrorWithMessage("Unable to get trip id").Respond(w, 422)
	}

	method := "osrm"
	if useGoogle {
		method = "google"
	}
	if usePassThrough {
		method = "passthrough"
	}
	redisClient := redis.GetClient()
	cacheKey := fmt.Sprintf("tripsummary-%d-%s", tripID, method)

	w.Header().Add("Content-Type", "application/json")
	cachedValue := redisClient.Get(cacheKey).Val()

	if cachedValue != "" {
		w.Write([]byte(cachedValue))
		return
	}
	trip, err := getTrip(r)
	if err != nil {
		ErrorWithMessage(fmt.Sprintf("Unable to find trip %s", err.Error())).Respond(w, 404)
		return
	}
	tripLocations, err := db.GetTripLocationsByTrip(db.CurrentDB(), trip.ID)
	if err != nil {
		ErrorWithMessage(fmt.Sprintf("Unable to get locations for trip %s", err.Error())).Respond(w, 404)
		return
	}

	var locs []utils.Location
	var timestamps []time.Time
	for _, tl := range tripLocations {
		locs = append(locs, tl.Location.Location)
		timestamps = append(timestamps, tl.Time)
	}
	var client services.RoadsService
	if useGoogle {
		client = services.GetGoogleRoadsService()
	} else if usePassThrough {
		client = &services.PassThroughRoadsService{}
	} else {
		client = services.NewOSRMClient(getOSRMURL())
	}

	resp, err := client.Match(locs, timestamps)
	if err != nil {
		ErrorWithMessage(fmt.Sprintf("Unable to get matching for trip %d, %s", trip.ID, err.Error())).Respond(w, 500)
		return
	}
	encodedJson, err := json.Marshal(resp)
	if err != nil {
		ErrorWithMessage(fmt.Sprintf("Unable to encode match response for trip %d, %s", trip.ID, err.Error())).Respond(w, 500)
		return
	}
	if redisClient.Set(cacheKey, encodedJson, 0).Err() != nil {
		log.Error("Could not cache response")
	}

	w.Write(encodedJson)
}

func getOSRMURL() string {
	osrmURL := os.Getenv("OSRM_URL")
	if osrmURL == "" {
		osrmURL = "http://ec2-13-127-26-106.ap-south-1.compute.amazonaws.com:5000"
	}
	return osrmURL
}

package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MOOVE-Network/location_service/db"
	"github.com/MOOVE-Network/location_service/utils"
	log "github.com/sirupsen/logrus"
)

type UCLRequest struct {
	Values []*Dict `json:"values"`
}

type Dict struct {
	NameValuePairs *TripLocation `json:"nameValuePairs"`
}
type TripLocation struct {
	LatString  string `json:"lat"`
	LngString  string `json:"lng"`
	Lat        float64
	Lng        float64
	Distance   int     `json:"distance"`
	Speed      float64 `json:"speed"`
	TripID     string  `json:"tripId"`
	TimeString string  `json:"time"`
	Time       time.Time
}

func (tl *TripLocation) ToDB() (*db.TripLocation, error) {

	location := db.Location{utils.Location{Lat: tl.Lat, Lng: tl.Lng}}

	tripId, err := strconv.Atoi(tl.TripID)

	if err != nil {
		return nil, err
	}

	dbTL := &db.TripLocation{
		TripID:    tripId,
		Location:  location,
		Time:      tl.Time,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Speed:     fmt.Sprintf("%f", tl.Speed),
		Distance:  tl.Distance,
	}
	return dbTL, nil
}

func (tl *TripLocation) parseTime() {
	dtString := tl.TimeString
	idx := strings.LastIndex(dtString, ":")
	dtString = fmt.Sprintf("%s%s", dtString[:idx], dtString[idx+1:])
	t, err := time.Parse("Mon Jan 2 15:04:05 MST-0700 2006", dtString)
	if err != nil {
		log.Warn("unable to parse time", tl.TimeString, err)
	}
	tl.Time = t
}

func (tl *TripLocation) parseLatLng() {
	lat, errLat := strconv.ParseFloat(tl.LatString, 64)
	lng, errLng := strconv.ParseFloat(tl.LngString, 64)
	if errLat != nil {
		log.Warn("Unable to parse lat ", tl.LatString, errLat)
	}
	tl.Lat = lat
	if errLng != nil {
		log.Warn("Unable to parse lng ", tl.LngString, errLng)
	}
	tl.Lng = lng
}
func (tl *TripLocation) parse() {
	tl.parseTime()
	tl.parseLatLng()
}

func UpdateCurrentLocation(w http.ResponseWriter, r *http.Request) {
	var uclRequest *UCLRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warn("Error reading location request body", err)
		ErrorWithMessage("Unable to update current location body").Respond(w, 422)
		return
	}
	log.Info(string(body))
	jsonBody := bytes.NewBuffer(body)
	decoder := json.NewDecoder(jsonBody)
	err = decoder.Decode(&uclRequest)
	if err != nil {
		log.Warn("unable to decode json request", err)
		ErrorWithMessage("unable decode json").Respond(w, 422)
		return
	}
	var tls []*db.TripLocation
	for _, v := range uclRequest.Values {
		v.NameValuePairs.parseTime()
		tl, err := v.NameValuePairs.ToDB()
		if err != nil {
			log.Warn(err)
			continue
		}
		tls = append(tls, tl)
	}
	err = saveTripLocations(tls)
	if err != nil {
		log.Warn(err)
		ErrorWithMessage("something went wrong while saving trips").Respond(w, 422)
	}
	writeOk(w)
}

func saveTripLocations(tls []*db.TripLocation) error {
	tx, err := db.CurrentDB().Beginx()
	if err != nil {
		return err
	}
	stmt := db.InsertTripLocationStatement(tx)
	if err != nil {
		return err
	}
	for _, tl := range tls {
		err := tl.Save(stmt)
		if err != nil {
			log.Warn(err)
			return tx.Rollback()
		}
	}
	return tx.Commit()
}

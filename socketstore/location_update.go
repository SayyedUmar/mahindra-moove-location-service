package socketstore

import (
	"bytes"
	"encoding/json"
	"strconv"
	"time"

	"gopkg.in/guregu/null.v3"

	"github.com/MOOVE-Network/location_service/db"
	"github.com/MOOVE-Network/location_service/models"
	"github.com/MOOVE-Network/location_service/utils"
)

// LocationUpdate maps location updates to Redis
type LocationUpdate struct {
	TripID      int `json:"tripId"`
	UserID      string
	Lat         float64
	Lng         float64
	Distance    int
	Speed       float64
	Accuracy    float64
	Timestamp   time.Time
	UnixSeconds int64
}

func LocationUpdateFromJSON(msg []byte) (*LocationUpdate, error) {
	var locUpdate LocationUpdate
	buffer := bytes.NewBuffer(msg)
	decoder := json.NewDecoder(buffer)
	err := decoder.Decode(&locUpdate)
	if err != nil {
		return nil, err
	}
	locUpdate.Timestamp = time.Unix(locUpdate.UnixSeconds/1000, 0)
	return &locUpdate, nil
}
func (lu *LocationUpdate) ToTripLocation() db.TripLocation {
	return db.TripLocation{
		TripID:    lu.TripID,
		Location:  db.Location{utils.Location{Lat: lu.Lat, Lng: lu.Lng}},
		Time:      lu.Timestamp,
		Speed:     strconv.FormatFloat(lu.Speed, 'f', -1, 64),
		Distance:  lu.Distance,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

}

func (lu *LocationUpdate) ToDriverLocation() models.DriverLocation {
	return models.DriverLocation{
		RecordedAt: lu.Timestamp,
		TripID:     null.IntFrom(int64(lu.TripID)),
		UserID:     null.StringFrom(lu.UserID),
		Location:   models.Point{Lat: lu.Lat, Lng: lu.Lng},
		Distance:   lu.Distance,
		Speed:      lu.Speed,
		Accuracy:   lu.Accuracy,
		CreatedAt:  time.Now(),
	}
}

package socketstore

import (
	"bytes"
	"encoding/json"
	"time"
)

// LocationUpdate maps location updates to Redis
type LocationUpdate struct {
	TripID      int `json:"tripId"`
	DriverID    int
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
	return &locUpdate, nil
}

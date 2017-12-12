package socketstore

import (
	"time"
)

type GeofenceEvent struct {
	GeofenceType string
	TripID       int
	TripRouteID  int
	DriverID     int
	Lat          float64
	Lng          float64
	Distance     int
	Speed        float64
	Timstamp     time.Time
}

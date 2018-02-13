package socketstore

import (
	"bytes"
	"encoding/json"
	"time"
)

type GeofenceEvent struct {
	TransitionType string
	TripID         int
	SiteID         int
	TripRouteID    int
	EmployeeID     int
	Lat            float64
	Lng            float64
	Speed          float32
	Bearing        float32
	Accuracy       float32
	Timestamp      time.Time
	UnixSeconds    int64
}

//GeofenceEventFromJSON cosntructs GeofenceEvent from json event
func GeofenceEventFromJSON(msg []byte) (*GeofenceEvent, error) {
	var geofenceEvent GeofenceEvent
	buffer := bytes.NewBuffer(msg)
	decoder := json.NewDecoder(buffer)
	err := decoder.Decode(&geofenceEvent)
	if err != nil {
		return nil, err
	}
	geofenceEvent.Timestamp = time.Unix(geofenceEvent.UnixSeconds/1000, 0)
	return &geofenceEvent, nil
}

//IsDwellEvent returns true if the event is Dwell type otherwise returns false.
func (ge *GeofenceEvent) IsDwellEvent() bool {
	return ge.TransitionType == "GEOFENCE_DWELL"
}

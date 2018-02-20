package socketstore

import (
	"bytes"
	"encoding/json"
	"time"
)

type GeofenceEvent struct {
	TransitionType string
	GeofenceType   string
	LocationType   string
	TripID         int
	SiteID         int
	TripRouteID    int
	BusStopName    string
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

//IsForEmployeeLocation returns true if the event is for Employee home location otherwise returns false.
func (ge *GeofenceEvent) IsForEmployeeLocation() bool {
	return ge.LocationType == "EMPLOYEE_HOME"
}

//IsForNodalPoint returns true if the event is for Nodal Bus stop location otherwise returns false.
func (ge *GeofenceEvent) IsForNodalPoint() bool {
	return ge.LocationType == "NODAL_POINT"
}

//IsForSite returns true if the event is for Site location otherwise returns false.
func (ge *GeofenceEvent) IsForSite() bool {
	return ge.LocationType == "SITE"
}

//IsForNarrowGeofence returns true if the event is triggered for narrow geofence area otherwise returns false.
//On Mobile side radius of geofence is 200 meters
func (ge *GeofenceEvent) IsForNarrowGeofence() bool {
	return ge.GeofenceType == "Narrow"
}

//IsForWiderGeofence returns true if the event is triggered for wider geofence area otherwise returns false.
//On Mobile side radius of geofence is 1500 meters
func (ge *GeofenceEvent) IsForWiderGeofence() bool {
	return ge.GeofenceType == "Wider"
}

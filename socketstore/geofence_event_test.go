package socketstore

import (
	"fmt"
	"testing"
	"time"

	tst "github.com/MOOVE-Network/location_service/testutils"
	"github.com/stretchr/testify/assert"
)

func createGeofenceEvent(transitionType string, locationType string, geofenceType string) GeofenceEvent {
	return GeofenceEvent{
		TransitionType: transitionType,
		LocationType:   locationType,
		GeofenceType:   geofenceType,
	}
}

func TestGeofenceEventFromJSON(t *testing.T) {
	currentTime := time.Now()
	geofenceJSONEvent := []byte(fmt.Sprintf(`
		{
			"transitionType":"GEOFENCE_ENTER",
			"locationType" : "SITE",
			"geofenceType" : "Wider",
			"lat":13.036566,
			"lng":80.216382,
			"bearing":0,
			"speed":30,
			"accuracy":20.5,
			"unixSeconds":%d,
			"tripId":36521,
			"siteId":2,
			"tripRouteId":442671,
			"busStopName":"bus stop 1"
		}
	`, currentTime.Unix()))
	geofenceEvent, err := GeofenceEventFromJSON(geofenceJSONEvent)
	tst.FailNowOnErr(t, err)
	assert.Equal(t, "GEOFENCE_ENTER", geofenceEvent.TransitionType)
	assert.Equal(t, "SITE", geofenceEvent.LocationType)
	assert.Equal(t, "Wider", geofenceEvent.GeofenceType)
	assert.Equal(t, 13.036566, geofenceEvent.Lat)
	assert.Equal(t, 80.216382, geofenceEvent.Lng)
	assert.Equal(t, float32(0.0), geofenceEvent.Bearing)
	assert.Equal(t, float32(30.0), geofenceEvent.Speed)
	assert.Equal(t, float32(20.5), geofenceEvent.Accuracy)
	assert.Equal(t, currentTime.Unix(), geofenceEvent.UnixSeconds)
	assert.Equal(t, 36521, geofenceEvent.TripID)
	assert.Equal(t, 2, geofenceEvent.SiteID)
	assert.Equal(t, 442671, geofenceEvent.TripRouteID)
	//Don't know how to test time.
	// assert.Equal(t, currentTime.Second(), geofenceEvent.Timestamp.Second())

	geofenceJSONEvent = []byte(fmt.Sprintf(`
		{
			"transitionType":"GEOFENCE_ENTER",
			"locationType" : "SITE",
			"geofenceType" : "Wider",
			"lat":13.036566,
			"lng":80.216382,
			"bearing":0,
			"speed":30,
			"accuracy":20.5,
			"unixSeconds":%d,
			"tripId":36521,
			"siteId":2,
			"tripRouteId":442671,
			"busStopName":1
		}
	`, currentTime.Unix()))
	//Giving wrong busStopName type
	geofenceEvent, err = GeofenceEventFromJSON(geofenceJSONEvent)
	assert.Error(t, err)
	assert.Nil(t, geofenceEvent)
}

func TestGeofenceEvent_IsDwellEvent(t *testing.T) {
	geofenceEvent := createGeofenceEvent("GEOFENCE_ENTER", "", "")
	assert.False(t, geofenceEvent.IsDwellEvent())
	geofenceEvent = createGeofenceEvent("GEOFENCE_DWELL", "", "")
	assert.True(t, geofenceEvent.IsDwellEvent())
	geofenceEvent = createGeofenceEvent("GEOFENCE_EXIT", "", "")
	assert.False(t, geofenceEvent.IsDwellEvent())
}

func TestGeofenceEventLocationTypes(t *testing.T) {
	geofenceEvent := createGeofenceEvent("GEOFENCE_ENTER", "SITE", "")
	assert.False(t, geofenceEvent.IsForEmployeeLocation())
	assert.True(t, geofenceEvent.IsForSite())
	assert.False(t, geofenceEvent.IsForNodalPoint())
	geofenceEvent = createGeofenceEvent("GEOFENCE_DWELL", "NODAL_POINT", "")
	assert.False(t, geofenceEvent.IsForEmployeeLocation())
	assert.False(t, geofenceEvent.IsForSite())
	assert.True(t, geofenceEvent.IsForNodalPoint())
	geofenceEvent = createGeofenceEvent("GEOFENCE_EXIT", "EMPLOYEE_HOME", "")
	assert.True(t, geofenceEvent.IsForEmployeeLocation())
	assert.False(t, geofenceEvent.IsForSite())
	assert.False(t, geofenceEvent.IsForNodalPoint())
}

func TestGeofenceTypes(t *testing.T) {
	geofenceEvent := createGeofenceEvent("GEOFENCE_ENTER", "SITE", "Wider")
	assert.True(t, geofenceEvent.IsForWiderGeofence())
	assert.False(t, geofenceEvent.IsForNarrowGeofence())
	geofenceEvent = createGeofenceEvent("GEOFENCE_DWELL", "NODAL_POINT", "Narrow")
	assert.False(t, geofenceEvent.IsForWiderGeofence())
	assert.True(t, geofenceEvent.IsForNarrowGeofence())
}

func TestGetLocation(t *testing.T) {
	geofenceEvent := createGeofenceEvent("GEOFENCE_ENTER", "SITE", "Wider")
	geofenceEvent.Lat = 13.0
	geofenceEvent.Lng = 79.0
	location := geofenceEvent.GetLocation()
	assert.Equal(t, geofenceEvent.Lat, location.Lat)
	assert.Equal(t, geofenceEvent.Lng, location.Lng)
}

package socketstore

import (
	"fmt"
	"testing"
	"time"

	tst "github.com/MOOVE-Network/location_service/testutils"
	"github.com/stretchr/testify/assert"
)

func createGeofenceEvent(transitionType string) GeofenceEvent {
	return GeofenceEvent{
		TransitionType: transitionType,
	}
}

func TestGeofenceEventFromJSON(t *testing.T) {
	currentTime := time.Now()
	geofenceJSONEvent := []byte(fmt.Sprintf(`
		{
			"transitionType":"GEOFENCE_ENTER",
			"lat":13.036566,
			"lng":80.216382,
			"bearing":0,
			"speed":30,
			"accuracy":20.5,
			"unixSeconds":%d,
			"tripId":36521,
			"siteId":2,
			"tripRouteId":442671,
			"employeeId":21
		}
	`, currentTime.Unix()))
	geofenceEvent, err := GeofenceEventFromJSON(geofenceJSONEvent)
	tst.FailNowOnErr(t, err)
	assert.Equal(t, "GEOFENCE_ENTER", geofenceEvent.TransitionType)
	assert.Equal(t, 13.036566, geofenceEvent.Lat)
	assert.Equal(t, 80.216382, geofenceEvent.Lng)
	assert.Equal(t, float32(0.0), geofenceEvent.Bearing)
	assert.Equal(t, float32(30.0), geofenceEvent.Speed)
	assert.Equal(t, float32(20.5), geofenceEvent.Accuracy)
	assert.Equal(t, currentTime.Unix(), geofenceEvent.UnixSeconds)
	assert.Equal(t, 36521, geofenceEvent.TripID)
	assert.Equal(t, 2, geofenceEvent.SiteID)
	assert.Equal(t, 442671, geofenceEvent.TripRouteID)
	assert.Equal(t, 21, geofenceEvent.EmployeeID)
	//Don't know how to test time.
	// assert.Equal(t, currentTime.Second(), geofenceEvent.Timestamp.Second())
}

func TestGeofenceEvent_IsDwellEvent(t *testing.T) {
	geofenceEvent := createGeofenceEvent("GEOFENCE_ENTER")
	assert.False(t, geofenceEvent.IsDwellEvent())
	geofenceEvent = createGeofenceEvent("GEOFENCE_DWELL")
	assert.True(t, geofenceEvent.IsDwellEvent())
	geofenceEvent = createGeofenceEvent("GEOFENCE_EXIT")
	assert.False(t, geofenceEvent.IsDwellEvent())
}

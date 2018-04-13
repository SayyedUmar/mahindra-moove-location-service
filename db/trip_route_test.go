package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	tst "github.com/MOOVE-Network/location_service/testutils"
	"github.com/MOOVE-Network/location_service/utils"
	"github.com/jmoiron/sqlx"
)

func createTripRoutes(t *Trip, trs []TripRoute, tx *sqlx.Tx) (*Trip, error) {
	stmt, err := tx.Preparex(`insert into trip_routes(trip_id, status, employee_trip_id, 
						 scheduled_route_order, scheduled_start_location, scheduled_end_location)
						 values (?,?,?,?,?,?)`)
	if err != nil {
		return nil, err
	}
	for routeOrder, tr := range trs {
		eid := getEmployeeID(tx, users[routeOrder%len(employees)].Email)
		etID := createEmployeeTrip(tx, eid)
		_, err := stmt.Exec(t.ID, tr.Status, etID, routeOrder, tr.ScheduledStartLocation, tr.ScheduledEndLocation)
		if err != nil {
			return nil, err
		}
	}
	err = t.LoadTripRoutes(tx, true)
	if err != nil {
		return nil, err
	}
	return t, nil
}
func createEmployeeTrip(tx *sqlx.Tx, employeeID int) int {
	res, err := tx.Exec(`
		insert into employee_trips(employee_id, created_at, updated_at)
		values(?,?,?)
		`, employeeID, time.Now(), time.Now())
	if err != nil {
		panic(err)
	}
	etID, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}
	return int(etID)
}

func TestIsNotStarted(t *testing.T) {
	tx := createTx(t)
	defer tx.Rollback()
	trip, err := createTripWithRoutes(tx, 23, 42, "not_started", "driver_arrived", "on_board")
	tst.FailNowOnErr(t, err)

	trip1, err := GetTripByID(tx, trip.ID)
	tst.FailNowOnErr(t, err)

	assert.Equal(t, trip.ID, trip1.ID)
	assert.Equal(t, 3, len(trip1.TripRoutes))
	assert.True(t, trip1.TripRoutes[0].IsNotStarted())
	assert.False(t, trip1.TripRoutes[1].IsNotStarted())
	assert.False(t, trip1.TripRoutes[2].IsNotStarted())

	trip, err = createTripWithRoutes(tx, 23, 42, "missed", "canceled", "completed")
	tst.FailNowOnErr(t, err)

	trip1, err = GetTripByID(tx, trip.ID)
	tst.FailNowOnErr(t, err)

	assert.Equal(t, trip.ID, trip1.ID)
	assert.Equal(t, 3, len(trip1.TripRoutes))
	assert.False(t, trip1.TripRoutes[0].IsNotStarted())
	assert.False(t, trip1.TripRoutes[1].IsNotStarted())
	assert.False(t, trip1.TripRoutes[2].IsNotStarted())
}

func TestIsMissedOrCanceled(t *testing.T) {
	tx := createTx(t)
	defer tx.Rollback()
	trip, err := createTripWithRoutes(tx, 23, 42, "not_started", "driver_arrived", "on_board")
	tst.FailNowOnErr(t, err)

	trip1, err := GetTripByID(tx, trip.ID)
	tst.FailNowOnErr(t, err)

	assert.Equal(t, trip.ID, trip1.ID)
	assert.Equal(t, 3, len(trip1.TripRoutes))
	assert.False(t, trip1.TripRoutes[0].IsMissedOrCanceled())
	assert.False(t, trip1.TripRoutes[1].IsMissedOrCanceled())
	assert.False(t, trip1.TripRoutes[2].IsMissedOrCanceled())

	trip, err = createTripWithRoutes(tx, 23, 42, "missed", "canceled", "completed")
	tst.FailNowOnErr(t, err)

	trip1, err = GetTripByID(tx, trip.ID)
	tst.FailNowOnErr(t, err)

	assert.Equal(t, trip.ID, trip1.ID)
	assert.Equal(t, 3, len(trip1.TripRoutes))
	assert.True(t, trip1.TripRoutes[0].IsMissedOrCanceled())
	assert.True(t, trip1.TripRoutes[1].IsMissedOrCanceled())
	assert.False(t, trip1.TripRoutes[2].IsMissedOrCanceled())
}

func TestIsTripRouteNotStarted(t *testing.T) {
	tx := createTx(t)
	defer tx.Rollback()
	trip, err := createTripWithRoutes(tx, 23, 42, "not_started", "driver_arrived", "on_board")
	tst.FailNowOnErr(t, err)

	trip1, err := GetTripByID(tx, trip.ID)
	tst.FailNowOnErr(t, err)

	assert.Equal(t, trip.ID, trip1.ID)
	assert.Equal(t, 3, len(trip1.TripRoutes))

	isNotStarted, err := IsTripRouteNotStarted(tx, trip1.TripRoutes[0].ID)
	tst.FailNowOnErr(t, err)
	assert.True(t, isNotStarted)

	isNotStarted, err = IsTripRouteNotStarted(tx, trip1.TripRoutes[1].ID)
	tst.FailNowOnErr(t, err)
	assert.False(t, isNotStarted)

	isNotStarted, err = IsTripRouteNotStarted(tx, trip1.TripRoutes[2].ID)
	tst.FailNowOnErr(t, err)
	assert.False(t, isNotStarted)

	trip, err = createTripWithRoutes(tx, 23, 42, "missed", "canceled", "completed")
	tst.FailNowOnErr(t, err)

	trip1, err = GetTripByID(tx, trip.ID)
	tst.FailNowOnErr(t, err)

	assert.Equal(t, trip.ID, trip1.ID)
	assert.Equal(t, 3, len(trip1.TripRoutes))

	isNotStarted, err = IsTripRouteNotStarted(tx, trip1.TripRoutes[0].ID)
	tst.FailNowOnErr(t, err)
	assert.False(t, isNotStarted)

	isNotStarted, err = IsTripRouteNotStarted(tx, trip1.TripRoutes[1].ID)
	tst.FailNowOnErr(t, err)
	assert.False(t, isNotStarted)

	isNotStarted, err = IsTripRouteNotStarted(tx, trip1.TripRoutes[2].ID)
	tst.FailNowOnErr(t, err)
	assert.False(t, isNotStarted)
}

func TestUpdateDriverArrivedGeofenceInfo(t *testing.T) {
	tx := createTx(t)
	defer tx.Rollback()
	trip, err := createTripWithRoutes(tx, 23, 42, "not_started", "driver_arrived", "on_board")
	tst.FailNowOnErr(t, err)

	location := Location{utils.Location{Lat: 12.0, Lng: 79.0}}
	now := time.Now().Round(time.Second) //Rounding since time gets rounded to nearest second during insert. don't know why.

	for _, tr := range trip.TripRoutes {
		err = tr.UpdateDriverArrivedGeofenceInfo(tx, location, now)
		tst.FailNowOnErr(t, err)
	}

	type Temp struct {
		GeofenceDriverArrivedDate     time.Time `db:"geofence_driver_arrived_date"`
		GeofenceDriverArrivedLocation Location  `db:"geofence_driver_arrived_location"`
	}
	for _, tr := range trip.TripRoutes {
		var temp Temp

		row := tx.QueryRowx("SELECT geofence_driver_arrived_date, geofence_driver_arrived_location FROM trip_routes where id = ?", tr.ID)
		row.StructScan(&temp)
		assert.EqualValues(t, location, temp.GeofenceDriverArrivedLocation)
		assert.Equal(t, now.Unix(), temp.GeofenceDriverArrivedDate.Unix())
	}
}

func TestUpdateCompletedGeofenceInfo(t *testing.T) {
	tx := createTx(t)
	defer tx.Rollback()
	trip, err := createTripWithRoutes(tx, 23, 42, "not_started", "driver_arrived", "on_board")
	tst.FailNowOnErr(t, err)

	location := Location{utils.Location{Lat: 12.0, Lng: 79.0}}
	now := time.Now().Round(time.Second) //Rounding since time gets rounded to nearest second during insert. don't know why.

	for _, tr := range trip.TripRoutes {
		err = tr.UpdateCompletedGeofenceInfo(tx, location, now)
		tst.FailNowOnErr(t, err)
	}

	type Temp struct {
		GeofenceCompletedDate     time.Time `db:"geofence_completed_date"`
		GeofenceCompletedLocation Location  `db:"geofence_completed_location"`
	}
	for _, tr := range trip.TripRoutes {
		var temp Temp

		row := tx.QueryRowx("SELECT geofence_completed_date, geofence_completed_location FROM trip_routes where id = ?", tr.ID)
		row.StructScan(&temp)
		assert.EqualValues(t, location, temp.GeofenceCompletedLocation)
		assert.Equal(t, now.Unix(), temp.GeofenceCompletedDate.Unix())
	}
}

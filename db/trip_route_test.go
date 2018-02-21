package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	tst "github.com/MOOVE-Network/location_service/testutils"
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

	trip, err = createTripWithRoutes(tx, 23, 42, "missed", "cancelled", "completed")
	tst.FailNowOnErr(t, err)

	trip1, err = GetTripByID(tx, trip.ID)
	tst.FailNowOnErr(t, err)

	assert.Equal(t, trip.ID, trip1.ID)
	assert.Equal(t, 3, len(trip1.TripRoutes))
	assert.False(t, trip1.TripRoutes[0].IsNotStarted())
	assert.False(t, trip1.TripRoutes[1].IsNotStarted())
	assert.False(t, trip1.TripRoutes[2].IsNotStarted())
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

	trip, err = createTripWithRoutes(tx, 23, 42, "missed", "cancelled", "completed")
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

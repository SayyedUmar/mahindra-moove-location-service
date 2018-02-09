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

	trip1, err := createTripWithRoutes(tx, 23, 42, "missed", "cancelled", "completed")
	tst.FailNowOnErr(t, err)

	trips, err := GetTripsByStatus(tx, "active")
	tst.FailNowOnErr(t, err)

	assert.True(t, isInTripsArr(trips, trip))
	assert.True(t, isInTripsArr(trips, trip1))
	assert.Equal(t, 3, len(trips[0].TripRoutes))
	assert.Equal(t, 3, len(trips[1].TripRoutes))

	assert.True(t, trips[0].TripRoutes[0].IsNotStarted())
	assert.False(t, trips[0].TripRoutes[1].IsNotStarted())
	assert.False(t, trips[0].TripRoutes[2].IsNotStarted())
	assert.False(t, trips[1].TripRoutes[0].IsNotStarted())
	assert.False(t, trips[1].TripRoutes[1].IsNotStarted())
	assert.False(t, trips[1].TripRoutes[2].IsNotStarted())
}

func TestIsTripRouteNotStarted(t *testing.T) {
	tx := createTx(t)
	defer tx.Rollback()
	trip, err := createTripWithRoutes(tx, 23, 42, "not_started", "driver_arrived", "on_board")
	tst.FailNowOnErr(t, err)

	trip1, err := createTripWithRoutes(tx, 23, 42, "missed", "cancelled", "completed")
	tst.FailNowOnErr(t, err)

	trips, err := GetTripsByStatus(tx, "active")
	tst.FailNowOnErr(t, err)

	assert.True(t, isInTripsArr(trips, trip))
	assert.True(t, isInTripsArr(trips, trip1))
	assert.Equal(t, 3, len(trips[0].TripRoutes))
	assert.Equal(t, 3, len(trips[1].TripRoutes))

	isNotStarted, err := IsTripRouteNotStarted(tx, trips[0].TripRoutes[0].ID)
	tst.FailNowOnErr(t, err)
	assert.True(t, isNotStarted)

	isNotStarted, err = IsTripRouteNotStarted(tx, trips[0].TripRoutes[1].ID)
	tst.FailNowOnErr(t, err)
	assert.False(t, isNotStarted)

	isNotStarted, err = IsTripRouteNotStarted(tx, trips[0].TripRoutes[2].ID)
	tst.FailNowOnErr(t, err)
	assert.False(t, isNotStarted)

	isNotStarted, err = IsTripRouteNotStarted(tx, trips[1].TripRoutes[0].ID)
	tst.FailNowOnErr(t, err)
	assert.False(t, isNotStarted)

	isNotStarted, err = IsTripRouteNotStarted(tx, trips[1].TripRoutes[1].ID)
	tst.FailNowOnErr(t, err)
	assert.False(t, isNotStarted)

	isNotStarted, err = IsTripRouteNotStarted(tx, trips[1].TripRoutes[2].ID)
	tst.FailNowOnErr(t, err)
	assert.False(t, isNotStarted)
}

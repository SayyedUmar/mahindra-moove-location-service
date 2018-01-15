package db

import (
	"testing"
	"time"

	tst "github.com/MOOVE-Network/location_service/testutils"
	"github.com/MOOVE-Network/location_service/utils"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

const TarkaLabsLat float64 = 13.0363178
const TarkaLabsLng float64 = 80.2142206

const HomeOneLat float64 = 12.8947997
const HomeOneLng float64 = 80.2010107

const HomeTwoLat float64 = 12.9217982
const HomeTwoLng float64 = 80.1999803

const HomeThreeLat float64 = 12.9851054
const HomeThreeLng float64 = 80.1983123

func createTrip(t *Trip, tx *sqlx.Tx) (*Trip, error) {
	now := time.Now()

	res := tx.MustExec(`insert into trips(trip_type, driver_id, vehicle_id, status, created_at, updated_at)
				 values (?,?,?,?,?,?)`, t.TripType, t.DriverID, t.VehicleID, t.Status, now, now)
	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return GetTripByID(tx, int(lastID))
}

func createTripWithRoutes(tx *sqlx.Tx, driverID, vehicleID int, status ...string) (*Trip, error) {
	statusStr := "not_started"
	if len(status) == 0 {
		status = []string{statusStr}
	}
	statusIdx := 0
	trip, err := createTrip(&Trip{TripType: TripTypeCheckIn,
		DriverID:  driverID,
		VehicleID: vehicleID,
		Status:    "active"}, tx)
	if err != nil {
		return nil, err
	}

	tr1 := TripRoute{
		ScheduledStartLocation: Location{utils.Location{Lat: HomeOneLat, Lng: HomeOneLng}},
		ScheduledEndLocation:   Location{utils.Location{Lat: HomeTwoLat, Lng: HomeTwoLng}},
		Status:                 status[statusIdx],
	}
	if statusIdx < len(status)-1 {
		statusIdx++
	}

	tr2 := TripRoute{
		ScheduledStartLocation: Location{utils.Location{Lat: HomeTwoLat, Lng: HomeTwoLng}},
		ScheduledEndLocation:   Location{utils.Location{Lat: HomeThreeLat, Lng: HomeThreeLng}},
		Status:                 status[statusIdx],
	}
	if statusIdx < len(status)-1 {
		statusIdx++
	}

	tr3 := TripRoute{
		ScheduledStartLocation: Location{utils.Location{Lat: HomeThreeLat, Lng: HomeThreeLng}},
		ScheduledEndLocation:   Location{utils.Location{Lat: TarkaLabsLat, Lng: TarkaLabsLng}},
		Status:                 status[statusIdx],
	}
	if statusIdx < len(status)-1 {
		statusIdx++
	}

	return createTripRoutes(trip, []TripRoute{tr1, tr2, tr3}, tx)
}

func TestGetTripByID(t *testing.T) {
	tx := createTx(t)
	defer tx.Rollback()
	trip, err := createTrip(&Trip{TripType: TripTypeCheckIn,
		DriverID:  23,
		VehicleID: 42,
		Status:    "active"}, tx)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	assert.Equal(t, 23, trip.DriverID)
	assert.Equal(t, 42, trip.VehicleID)
	assert.Equal(t, "active", trip.Status)
}

func TestGetTripByIDLoadRoutes(t *testing.T) {
	tx := createTx(t)
	defer tx.Rollback()

	trip, err := createTripWithRoutes(tx, 23, 42)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	assert.Equal(t, 3, len(trip.TripRoutes))
	assert.False(t, trip.AllCheckedIn())
}

func TestTrip_AllCheckedIn(t *testing.T) {
	tx := createTx(t)
	defer tx.Rollback()
	trip, err := createTripWithRoutes(tx, 23, 42, "canceled", "on_board", "on_board")
	tst.FailNowOnErr(t, err)

	assert.Equal(t, 3, len(trip.TripRoutes))
	assert.True(t, trip.AllCheckedIn())
}

func isInTripsArr(trips []*Trip, tr *Trip) bool {
	found := false
	for _, t := range trips {
		if t.ID == tr.ID {
			return true
		}
	}
	return found
}

func TestGetTripsByStatus(t *testing.T) {
	tx := createTx(t)
	defer tx.Rollback()
	trip1, err := createTripWithRoutes(tx, 23, 42, "canceled", "on_board", "on_board")
	tst.FailNowOnErr(t, err)

	trip2, err := createTripWithRoutes(tx, 24, 43, "driver_arrived", "not_started", "not_started")
	tst.FailNowOnErr(t, err)

	trip3, err := createTripWithRoutes(tx, 25, 44, "not_started", "not_started", "not_started")
	tst.FailNowOnErr(t, err)

	trips, err := GetTripsByStatus(tx, "active")
	tst.FailNowOnErr(t, err)

	assert.Equal(t, len(trips), 3)
	assert.Equal(t, len(trips[0].TripRoutes), 3)
	assert.True(t, isInTripsArr(trips, trip1))
	assert.Equal(t, len(trips[1].TripRoutes), 3)
	assert.True(t, isInTripsArr(trips, trip2))
	assert.Equal(t, len(trips[2].TripRoutes), 3)
	assert.True(t, isInTripsArr(trips, trip3))
}

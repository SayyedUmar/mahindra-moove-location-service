package db

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"

	tst "github.com/MOOVE-Network/location_service/testutils"
	"github.com/MOOVE-Network/location_service/utils"
	"github.com/stretchr/testify/assert"
)

func tripLocationsRowCount(q sqlx.Queryer) (int, error) {
	var rowCount int
	r := q.QueryRowx("select count(1) as count from trip_locations")
	err := r.Scan(&rowCount)
	return rowCount, err
}

func TestTripLocation_Save(t *testing.T) {
	tx := createTx(t)
	defer tx.Rollback()
	trip, err := createTripWithRoutes(tx)
	tst.FailNowOnErr(t, err)
	tl1 := TripLocation{
		TripID:    trip.ID,
		Location:  Location{utils.Location{Lat: TarkaLabsLat, Lng: TarkaLabsLng}},
		Time:      time.Now(),
		Speed:     "20",
		Distance:  0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	beforeCount, err := tripLocationsRowCount(tx)
	tst.FailNowOnErr(t, err)

	insertStmt := InsertTripLocationStatement(tx)

	err = tl1.Save(insertStmt)
	tst.FailNowOnErr(t, err)

	afterCount, err := tripLocationsRowCount(tx)
	tst.FailNowOnErr(t, err)

	assert.Equal(t, afterCount, beforeCount+1)
}

func TestLatestTripLocation(t *testing.T) {
	tx := createTx(t)
	defer tx.Rollback()
	trip, err := createTripWithRoutes(tx)
	tst.FailNowOnErr(t, err)
	tl1 := TripLocation{
		TripID:    trip.ID,
		Location:  Location{utils.Location{Lat: TarkaLabsLat, Lng: TarkaLabsLng}},
		Time:      time.Now(),
		Speed:     "20",
		Distance:  0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	time.Sleep(200 * time.Millisecond)
	tl2 := TripLocation{
		TripID:    trip.ID,
		Location:  Location{utils.Location{Lat: TarkaLabsLat, Lng: TarkaLabsLng}},
		Time:      time.Now(),
		Speed:     "33",
		Distance:  0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	insertStmt := InsertTripLocationStatement(tx)

	err = tl1.Save(insertStmt)
	tst.FailNowOnErr(t, err)
	err = tl2.Save(insertStmt)
	tst.FailNowOnErr(t, err)

	tl, err := LatestTripLocation(tx, trip.ID)
	tst.FailNowOnErr(t, err)
	assert.Equal(t, tl.Speed, "33")
}

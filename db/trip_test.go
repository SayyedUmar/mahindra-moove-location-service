package db

import (
	"testing"
	"time"

	"gopkg.in/guregu/null.v3"

	tst "github.com/MOOVE-Network/location_service/testutils"
	"github.com/MOOVE-Network/location_service/utils"
	"github.com/icrowley/fake"
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

func ensureDriver(tx *sqlx.Tx, driverID int) {
	var dID int
	row := tx.QueryRow("select 1 from drivers where id=?", driverID)
	if row.Scan(&dID) != nil {
		driversInsertQuery := "insert into drivers (id, created_at, updated_at) values(?,?,?)"
		usersInsertQuery := `
		insert into users (f_name, l_name, email, encrypted_password, uid, sign_in_count, created_at, updated_at, provider,  entity_id, entity_type) 
					values(?, 		?,		?    , ?                 , ?  , ?            , ?         , ?         , ?       ,  ?        , ?          )`

		tx.MustExec(driversInsertQuery, driverID, time.Now(), time.Now())
		email := fake.EmailAddress()
		fName := fake.FirstName()
		lName := fake.LastName()
		tx.MustExec(usersInsertQuery, fName, lName, email, fake.SimplePassword(), email, 0, time.Now(), time.Now(), "provider", driverID, "driver")
	}
}

func createTrip(t *Trip, tx *sqlx.Tx) (*Trip, error) {
	now := time.Now()
	ensureDriver(tx, t.DriverID)
	res := tx.MustExec(`insert into trips(trip_type, driver_id, vehicle_id, status, scheduled_date, created_at, updated_at)
				 values (?,?,?,?,?,?,?)`, t.TripType, t.DriverID, t.VehicleID, t.Status, t.ScheduledStartDate, now, now)
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
	currentTime := time.Now().Round(time.Second) //Rounding since time gets rounded to nearest second during insert. don't know why.
	trip, err := createTrip(&Trip{TripType: TripTypeCheckIn,
		DriverID:           23,
		VehicleID:          42,
		Status:             "active",
		ScheduledStartDate: null.TimeFrom(currentTime),
	}, tx)
	tst.FailNowOnErr(t, err)

	assert.Equal(t, 23, trip.DriverID)
	assert.Equal(t, 42, trip.VehicleID)
	assert.True(t, trip.DriverUserID > 0)
	assert.Equal(t, "active", trip.Status)
	assert.True(t, trip.ScheduledStartDate.Valid)
	assert.Equal(t, currentTime.Unix(), trip.ScheduledStartDate.Time.Unix())
}

func TestGetTripByIDLoadRoutes(t *testing.T) {
	tx := createTx(t)
	defer tx.Rollback()

	trip, err := createTripWithRoutes(tx, 23, 42)
	tst.FailNowOnErr(t, err)

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

	trips, err := GetTripsByStatuses(tx, "active")
	tst.FailNowOnErr(t, err)

	assert.Equal(t, len(trips), 3)
	assert.Equal(t, len(trips[0].TripRoutes), 3)
	assert.True(t, isInTripsArr(trips, trip1))
	assert.Equal(t, len(trips[1].TripRoutes), 3)
	assert.True(t, isInTripsArr(trips, trip2))
	assert.Equal(t, len(trips[2].TripRoutes), 3)
	assert.True(t, isInTripsArr(trips, trip3))
}

func TestGetTripsByStatusFillsScheduledStartDate(t *testing.T) {
	tx := createTx(t)
	defer tx.Rollback()
	currentTime := time.Now().Round(time.Second) //Rounding since time gets rounded to nearest second during insert. don't know why.
	_, err := createTrip(&Trip{TripType: TripTypeCheckIn,
		DriverID:           23,
		VehicleID:          42,
		Status:             "active",
		ScheduledStartDate: null.TimeFrom(currentTime),
	}, tx)
	tst.FailNowOnErr(t, err)
	trips, err := GetTripsByStatuses(tx, "active")
	tst.FailNowOnErr(t, err)

	assert.Equal(t, 1, len(trips))
	assert.Equal(t, "active", trips[0].Status)
	assert.True(t, trips[0].ScheduledStartDate.Valid)
	assert.Equal(t, currentTime.Unix(), trips[0].ScheduledStartDate.Time.Unix())
}

func TestUpdateScheduleStartTripTimeAndLocation(t *testing.T) {
	tx := createTx(t)
	defer tx.Rollback()
	currentTime := time.Now().Round(time.Second) //Rounding since time gets rounded to nearest second during insert. don't know why.
	currentLocation := Location{
		utils.Location{
			Lat: 13.01,
			Lng: 70.01,
		},
	}
	trip, err := createTrip(&Trip{TripType: TripTypeCheckIn,
		DriverID:  23,
		VehicleID: 42,
		Status:    "active",
	}, tx)
	tst.FailNowOnErr(t, err)
	driverShouldStartTripTime := currentTime.Add(time.Duration(time.Hour * 1))
	err = trip.UpdateDriverShouldStartTripTimeAndLocation(tx, driverShouldStartTripTime, currentLocation, currentTime)
	tst.FailNowOnErr(t, err)

	getDriverShouldStartTripTimeAndLocationQuery := `
		select driver_should_start_trip_time, driver_should_start_trip_location, driver_should_start_trip_timestamp from trips where id=?
	`
	type TempStruct struct {
		DiverShouldStartTripTime            null.Time `db:"driver_should_start_trip_time"`
		DiverShouldStartTripLocation        Location  `db:"driver_should_start_trip_location"`
		DiverShouldStartTripCalculationTime null.Time `db:"driver_should_start_trip_timestamp"`
	}
	row := tx.QueryRowx(getDriverShouldStartTripTimeAndLocationQuery, trip.ID)
	var tempStruct TempStruct
	err = row.StructScan(&tempStruct)
	tst.FailNowOnErr(t, err)

	assert.True(t, tempStruct.DiverShouldStartTripTime.Valid)
	assert.True(t, tempStruct.DiverShouldStartTripCalculationTime.Valid)
	assert.Equal(t, driverShouldStartTripTime.Unix(), tempStruct.DiverShouldStartTripTime.Time.Unix())
	assert.Equal(t, currentLocation, tempStruct.DiverShouldStartTripLocation)
	assert.Equal(t, currentTime.Unix(), tempStruct.DiverShouldStartTripCalculationTime.Time.Unix())
}

func TestGetDriverShouldStartTripLocation(t *testing.T) {
	tx := createTx(t)
	defer tx.Rollback()

	trip, err := createTrip(&Trip{TripType: TripTypeCheckIn,
		DriverID:  23,
		VehicleID: 42,
		Status:    "active",
	}, tx)
	tst.FailNowOnErr(t, err)

	//Testing for error, if DriverShouldStartTripLocation is nil than GetDriverShouldStartTripLocation should return error.
	startTripLocation, err := trip.GetDriverShouldStartTripLocation(tx)
	assert.Error(t, err)

	location := Location{
		utils.Location{
			Lat: 13.01,
			Lng: 79.01,
		},
	}

	trip.UpdateDriverShouldStartTripTimeAndLocation(tx, time.Now(), location, time.Now())
	startTripLocation, err = trip.GetDriverShouldStartTripLocation(tx)
	tst.FailNowOnErr(t, err)

	assert.Equal(t, location, *startTripLocation)
}

func TestHasTripStarted(t *testing.T) {
	trip := Trip{
		Status: "active",
	}
	assert.True(t, trip.HasStarted())

	for _, status := range []string{"assigned", "assign_requested", "assign_request_expired", "completed", "canceled"} {
		trip.Status = status
		assert.False(t, trip.HasStarted())
	}
}

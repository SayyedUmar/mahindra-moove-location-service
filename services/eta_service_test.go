package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	tst "github.com/MOOVE-Network/location_service/testutils"
	"gopkg.in/guregu/null.v3"

	"strconv"

	"github.com/MOOVE-Network/location_service/db"
	"github.com/MOOVE-Network/location_service/services"
	"github.com/MOOVE-Network/location_service/services/mocks"
	"github.com/MOOVE-Network/location_service/utils"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}
func TestGetETAForTripShould_NotifyETAForASimpleCheckinTrip(t *testing.T) {
	// setup Mocks
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockDurationService := mocks.NewMockDurationService(mockController)
	services.SetDurationService(mockDurationService)
	mockNotificationService := mocks.NewMockNotificationService(mockController)
	services.SetNotificationService(mockNotificationService)

	// setup current location
	currentLocation := db.Location{utils.Location{2, 2}}
	trip := simpleCheckinTrip()

	duration := services.DurationMetrics{Duration: 20 * time.Minute}
	mockDurationService.EXPECT().GetDuration(currentLocation, db.Location{utils.Location{3, 3}}, mockClock{}.Now()).Return(duration, nil)
	data := make(map[string]interface{})
	data["duration"] = int64(duration.Duration.Minutes())
	data["push_type"] = "driver_location_update_1"
	mockNotificationService.EXPECT().SendNotification(strconv.Itoa(4212), data, "user")

	_, err := services.GetETAForTrip(&trip, currentLocation, mockClock{})
	if err != nil {
		t.Log("Error getting ETA for trip")
		t.Log(err)
		t.Fail()
	}
}

func TestGetETAForTripShould_NotifyETAForACheckinTripWithEmpNotStarted(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockDurationService := mocks.NewMockDurationService(mockController)
	services.SetDurationService(mockDurationService)
	mockNotificationService := mocks.NewMockNotificationService(mockController)
	services.SetNotificationService(mockNotificationService)

	// setup current location
	currentLocation := db.Location{utils.Location{2, 2}}
	trip := checkinTripNotAllBoard()

	duration := services.DurationMetrics{Duration: 20 * time.Minute}

	durationCall1 := mockDurationService.EXPECT().GetDuration(currentLocation, db.Location{utils.Location{3, 3}}, mockClock{}.Now()).Return(duration, nil)
	notificationData1 := make(map[string]interface{})
	notificationData1["duration"] = int64(duration.Duration.Minutes())
	notificationData1["push_type"] = "driver_location_update_1"
	userCall1 := mockNotificationService.EXPECT().SendNotification(strconv.Itoa(4223), notificationData1, "user")

	mockDurationService.EXPECT().GetDuration(db.Location{utils.Location{3, 3}}, db.Location{utils.Location{4, 4}}, mockClock{}.Now().Add(duration.Duration)).Return(duration, nil).After(durationCall1)
	notificationData2 := make(map[string]interface{})
	notificationData2["duration"] = int64((duration.Duration + duration.Duration).Minutes())
	notificationData2["push_type"] = "driver_location_update_1"
	mockNotificationService.EXPECT().SendNotification(strconv.Itoa(4212), notificationData2, "user").After(userCall1)

	_, err := services.GetETAForTrip(&trip, currentLocation, mockClock{})
	if err != nil {
		t.Log("Error getting ETA for trip")
		t.Log(err)
		t.Fail()
	}
}

func TestGetETAForTripShould_NotifyETAForACheckinTripWithOffset(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockDurationService := mocks.NewMockDurationService(mockController)
	services.SetDurationService(mockDurationService)
	mockNotificationService := mocks.NewMockNotificationService(mockController)
	services.SetNotificationService(mockNotificationService)

	// setup current location
	currentLocation := db.Location{utils.Location{2, 2}}
	trip := checkinTripNotAllBoardWithOffset()

	duration := services.DurationMetrics{Duration: 20 * time.Minute}

	durationCall1 := mockDurationService.EXPECT().GetDuration(currentLocation, db.Location{utils.Location{4, 4}}, mockClock{}.Now()).Return(duration, nil)
	notificationData1 := make(map[string]interface{})
	notificationData1["duration"] = int64(duration.Duration.Minutes())
	notificationData1["push_type"] = "driver_location_update_1"
	mockNotificationService.EXPECT().SendNotification(strconv.Itoa(4233), notificationData1, "user")

	mockDurationService.EXPECT().GetDuration(db.Location{utils.Location{4, 4}}, db.Location{utils.Location{5, 5}}, mockClock{}.Now().Add(duration.Duration)).Return(duration, nil).After(durationCall1)

	notificationData2 := make(map[string]interface{})
	notificationData2["duration"] = int64((duration.Duration + duration.Duration).Minutes())
	notificationData2["push_type"] = "driver_location_update_1"
	mockNotificationService.EXPECT().SendNotification(strconv.Itoa(4212), notificationData2, "user")

	notificationData3 := make(map[string]interface{})
	notificationData3["duration"] = int64((duration.Duration + duration.Duration).Minutes())
	notificationData3["push_type"] = "driver_location_update_1"
	mockNotificationService.EXPECT().SendNotification(strconv.Itoa(4223), notificationData3, "user")

	_, err := services.GetETAForTrip(&trip, currentLocation, mockClock{})
	if err != nil {
		t.Log("Error getting ETA for trip")
		t.Log(err)
		t.Fail()
	}
}

func TestGetETAForTripShould_NotifyETAForACheckinTripWithOneOnBoard(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockDurationService := mocks.NewMockDurationService(mockController)
	services.SetDurationService(mockDurationService)
	mockNotificationService := mocks.NewMockNotificationService(mockController)
	services.SetNotificationService(mockNotificationService)

	// setup current location
	currentLocation := db.Location{utils.Location{2, 2}}
	trip := checkinTripNotAllBoardWithOffsetWithOneOnBoard()

	duration := services.DurationMetrics{Duration: 20 * time.Minute}

	durationCall1 := mockDurationService.EXPECT().GetDuration(currentLocation, db.Location{utils.Location{3, 3}}, mockClock{}.Now()).Return(duration, nil)
	notificationData1 := make(map[string]interface{})
	notificationData1["duration"] = int64(duration.Duration.Minutes())
	notificationData1["push_type"] = "driver_location_update_1"
	mockNotificationService.EXPECT().SendNotification(strconv.Itoa(4223), notificationData1, "user")

	durationCall2 := mockDurationService.EXPECT().GetDuration(db.Location{utils.Location{3, 3}}, db.Location{utils.Location{4, 4}}, mockClock{}.Now().Add(duration.Duration)).Return(duration, nil).After(durationCall1)
	notificationData2 := make(map[string]interface{})
	notificationData2["duration"] = int64((duration.Duration + duration.Duration).Minutes())
	notificationData2["push_type"] = "driver_location_update_1"
	mockNotificationService.EXPECT().SendNotification(strconv.Itoa(4233), notificationData2, "user")

	mockDurationService.EXPECT().GetDuration(db.Location{utils.Location{4, 4}}, db.Location{utils.Location{5, 5}}, mockClock{}.Now().Add(duration.Duration*2)).Return(duration, nil).After(durationCall2)
	notificationData3 := make(map[string]interface{})
	notificationData3["duration"] = int64((duration.Duration * 3).Minutes())
	notificationData3["push_type"] = "driver_location_update_1"
	mockNotificationService.EXPECT().SendNotification(strconv.Itoa(4212), notificationData3, "user")

	_, err := services.GetETAForTrip(&trip, currentLocation, mockClock{})
	if err != nil {
		t.Log("Error getting ETA for trip")
		t.Log(err)
		t.Fail()
	}
}

type mockClock struct {
}

func (mc mockClock) Now() time.Time {
	return time.Date(2018, 01, 01, 9, 0, 0, 0, time.Local)
}

// #region checkin
func simpleCheckinTrip() db.Trip {
	return db.Trip{
		ID:           42,
		TripType:     db.TripTypeCheckIn,
		DriverID:     43,
		DriverUserID: 400,
		VehicleID:    23,
		Status:       "active",
		TripRoutes: []db.TripRoute{
			db.TripRoute{
				EmployeeUserID:         4212,
				ID:                     421,
				TripID:                 42,
				Status:                 "on_board",
				Trip:                   db.Trip{DriverUserID: 400},
				ScheduledRouteOrder:    1,
				ScheduledStartLocation: db.Location{utils.Location{1, 1}},
				ScheduledEndLocation:   db.Location{utils.Location{3, 3}},
			},
		},
	}
}

func checkinTripNotAllBoard() db.Trip {
	return db.Trip{
		ID:           42,
		TripType:     db.TripTypeCheckIn,
		DriverID:     43,
		DriverUserID: 400,
		VehicleID:    23,
		Status:       "active",
		TripRoutes: []db.TripRoute{
			db.TripRoute{
				EmployeeUserID:         4212,
				ID:                     421,
				TripID:                 42,
				Status:                 "on_board",
				Trip:                   db.Trip{DriverUserID: 400},
				ScheduledRouteOrder:    1,
				ScheduledStartLocation: db.Location{utils.Location{1, 1}},
				ScheduledEndLocation:   db.Location{utils.Location{3, 3}},
			},
			db.TripRoute{
				EmployeeUserID:         4223,
				ID:                     422,
				TripID:                 42,
				Status:                 "not_started",
				Trip:                   db.Trip{DriverUserID: 400},
				ScheduledRouteOrder:    2,
				ScheduledStartLocation: db.Location{utils.Location{3, 3}},
				ScheduledEndLocation:   db.Location{utils.Location{4, 4}},
			},
		},
	}
}

func TestFindWhenShouldDriverStartTrip_ForCheckInTrip(t *testing.T) {
	driverLocation := db.Location{
		utils.Location{
			Lat: 1,
			Lng: 1,
		},
	}
	trip := makeAssignedCheckinTrip()

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockDurationService := mocks.NewMockDurationService(mockController)
	services.SetDurationService(mockDurationService)

	dm := services.DurationMetrics{
		Duration: time.Duration(time.Minute * 10),
	}

	clock := mockClock{}
	mockDurationService.EXPECT().GetDuration(driverLocation, trip.TripRoutes[0].ScheduledStartLocation, clock.Now()).Return(dm, nil).Times(1)

	newStartTime, err := services.FindWhenShouldDriverStartTrip(trip, &driverLocation, clock)
	tst.FailNowOnErr(t, err)

	assert.Equal(t, trip.ScheduledStartDate.Time.Add(-dm.Duration), *newStartTime)

	//Testing for Errors.

	//Testing if duration service gives error than function should return error.
	mockDurationService.EXPECT().GetDuration(driverLocation, trip.TripRoutes[0].ScheduledStartLocation, clock.Now()).Return(services.DurationMetrics{}, errors.New("Some Error")).Times(1)
	_, err = services.FindWhenShouldDriverStartTrip(trip, &driverLocation, clock)
	assert.EqualError(t, err, "Some Error")

	//Should give error as there is no schedule start date for first pickup.
	trip.ScheduledStartDate = null.NewTime(time.Time{}, false)

	mockDurationService.EXPECT().GetDuration(driverLocation, trip.TripRoutes[0].ScheduledStartLocation, clock.Now()).Times(0)

	_, err = services.FindWhenShouldDriverStartTrip(trip, &driverLocation, clock)
	assert.Error(t, err)

	//Checking for trip with no trip routes. should give error.
	trip.ScheduledStartDate = null.TimeFrom(clock.Now())
	trip.TripRoutes = []db.TripRoute{}

	mockDurationService.EXPECT().GetDuration(driverLocation, gomock.Any(), clock.Now()).Times(0)

	_, err = services.FindWhenShouldDriverStartTrip(trip, &driverLocation, clock)
	assert.Error(t, err)
}
func TestNotifyDriverShouldStartTrip(t *testing.T) {
	trip := makeAssignedCheckinTrip()

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockNotificationService := mocks.NewMockNotificationService(mockController)
	services.SetNotificationService(mockNotificationService)

	clock := mockClock{}

	newStartTime := clock.Now()
	calculationTime := clock.Now()
	data := make(map[string]interface{})
	data["push_type"] = "driver_should_start_trip"
	data["trip_id"] = trip.ID
	data["driver_should_start_trip_time"] = newStartTime.Unix()
	data["driver_should_start_trip_timestamp"] = calculationTime.Unix()
	data["driver_should_start_trip_calc_time"] = calculationTime.Unix()
	mockNotificationService.EXPECT().SendNotification(strconv.Itoa(trip.DriverUserID), data, "user").Return(nil).Times(1)

	sent, err := services.NotifyDriverShouldStartTrip(trip, &newStartTime, &calculationTime)
	tst.FailNowOnErr(t, err)
	assert.True(t, sent)

	//Testing for Errors.
	//Returns false if notification services returns error.
	mockNotificationService.EXPECT().SendNotification(strconv.Itoa(trip.DriverUserID), data, "user").Return(errors.New("Some Error")).Times(1)

	sent, err = services.NotifyDriverShouldStartTrip(trip, &newStartTime, &calculationTime)
	assert.EqualError(t, err, "Some Error")
	assert.False(t, sent)
}

func TestGetETAForBusTripShould_NotifyETAForACheckinTripWithEmpNotStarted(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockDurationService := mocks.NewMockDurationService(mockController)
	services.SetDurationService(mockDurationService)
	mockNotificationService := mocks.NewMockNotificationService(mockController)
	services.SetNotificationService(mockNotificationService)
	clock := mockClock{}

	// setup current location
	currentLocation := db.Location{utils.Location{1, 1}}
	trip := checkinBusTripNoneBoardWithOffset()

	duration := services.DurationMetrics{Duration: 20 * time.Minute}

	durationCall1 := mockDurationService.EXPECT().GetDuration(currentLocation, db.Location{utils.Location{3, 3}}, clock.Now()).Return(duration, nil)
	notificationData1 := make(map[string]interface{})
	notificationData1["duration"] = int64(duration.Duration.Minutes())
	notificationData1["push_type"] = "driver_location_update_1"
	log.Debugf("Expecting notification data : %v\n", notificationData1)
	notif1 := mockNotificationService.EXPECT().SendNotification("4213", notificationData1, "user").Return(nil)
	// mockNotificationService.EXPECT().SendNotification(strconv.Itoa(4213), notificationData1, "user").Return(nil).Times(1)
	mockNotificationService.EXPECT().SendNotification(strconv.Itoa(4223), notificationData1, "user").After(notif1)

	durationCall2 := mockDurationService.EXPECT().GetDuration(db.Location{utils.Location{3, 3}}, db.Location{utils.Location{4, 4}}, clock.Now().Add(duration.Duration)).Return(duration, nil).After(durationCall1)
	notificationData2 := make(map[string]interface{})
	notificationData2["duration"] = int64((duration.Duration + duration.Duration).Minutes())
	notificationData2["push_type"] = "driver_location_update_1"
	mockNotificationService.EXPECT().SendNotification(strconv.Itoa(4233), notificationData2, "user")
	mockNotificationService.EXPECT().SendNotification(strconv.Itoa(4243), notificationData2, "user")

	mockDurationService.EXPECT().GetDuration(db.Location{utils.Location{4, 4}}, db.Location{utils.Location{5, 5}}, clock.Now().Add(duration.Duration*2)).Return(duration, nil).After(durationCall2)
	notificationData3 := make(map[string]interface{})
	notificationData3["duration"] = int64((duration.Duration * 3).Minutes())
	notificationData3["push_type"] = "driver_location_update_1"
	mockNotificationService.EXPECT().SendNotification(strconv.Itoa(4253), notificationData3, "user")

	etaResp, err := services.GetETAForTrip(&trip, currentLocation, clock)
	if err != nil {
		t.Log("Error getting ETA for trip")
		t.Log(err)
		t.Fail()
	}

	desiredResponse := services.ETAResponse{
		ID:        trip.ID,
		UpdatedAt: clock.Now(),
		TripRoutes: []services.ETATripRoute{
			services.ETATripRoute{
				ID:             421,
				Status:         "not_started",
				EmployeeUserID: 4213,
				PickupTime:     services.NotNullTime(clock.Now().Add(duration.Duration)),
				ETAInMinutes:   duration.Duration.Minutes(),
			}, services.ETATripRoute{
				ID:             422,
				Status:         "not_started",
				EmployeeUserID: 4223,
				PickupTime:     services.NotNullTime(clock.Now().Add(duration.Duration)),
				ETAInMinutes:   duration.Duration.Minutes(),
			}, services.ETATripRoute{
				ID:             423,
				Status:         "not_started",
				EmployeeUserID: 4233,
				PickupTime:     services.NotNullTime(clock.Now().Add(duration.Duration * 2)),
				ETAInMinutes:   duration.Duration.Minutes() * 2,
			}, services.ETATripRoute{
				ID:             424,
				Status:         "not_started",
				EmployeeUserID: 4243,
				PickupTime:     services.NotNullTime(clock.Now().Add(duration.Duration * 2)),
				ETAInMinutes:   duration.Duration.Minutes() * 2,
			}, services.ETATripRoute{
				ID:             425,
				Status:         "not_started",
				EmployeeUserID: 4253,
				PickupTime:     services.NotNullTime(clock.Now().Add(duration.Duration * 3)),
				ETAInMinutes:   duration.Duration.Minutes() * 3,
			},
		},
	}
	assert.Equal(t, desiredResponse, *etaResp)
}

func TestGetETAForBusTripShould_NotifyETAForACheckoutTripWithEmpOnBoardAndDriverArrived(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockDurationService := mocks.NewMockDurationService(mockController)
	services.SetDurationService(mockDurationService)
	mockNotificationService := mocks.NewMockNotificationService(mockController)
	services.SetNotificationService(mockNotificationService)
	clock := mockClock{}

	// setup current location
	currentLocation := db.Location{utils.Location{1, 1}}
	trip := checkoutBusTripMixedOnBoardAndDriverArrivedWithOffset()

	duration := services.DurationMetrics{Duration: 20 * time.Minute}

	durationCall1 := mockDurationService.EXPECT().GetDuration(currentLocation, db.Location{utils.Location{3, 3}}, clock.Now()).Return(duration, nil)
	notificationData1 := make(map[string]interface{})
	notificationData1["duration"] = int64(duration.Duration.Minutes())
	notificationData1["push_type"] = "driver_location_update_1"
	log.Debugf("Expecting notification data : %v\n", notificationData1)
	notif1 := mockNotificationService.EXPECT().SendNotification("4213", notificationData1, "user").Return(nil)
	// mockNotificationService.EXPECT().SendNotification(strconv.Itoa(4213), notificationData1, "user").Return(nil).Times(1)
	mockNotificationService.EXPECT().SendNotification(strconv.Itoa(4223), notificationData1, "user").After(notif1)

	durationCall2 := mockDurationService.EXPECT().GetDuration(db.Location{utils.Location{3, 3}}, db.Location{utils.Location{4, 4}}, clock.Now().Add(duration.Duration)).Return(duration, nil).After(durationCall1)
	notificationData2 := make(map[string]interface{})
	notificationData2["duration"] = int64((duration.Duration + duration.Duration).Minutes())
	notificationData2["push_type"] = "driver_location_update_1"
	mockNotificationService.EXPECT().SendNotification(strconv.Itoa(4233), notificationData2, "user")
	mockNotificationService.EXPECT().SendNotification(strconv.Itoa(4243), notificationData2, "user")

	mockDurationService.EXPECT().GetDuration(db.Location{utils.Location{4, 4}}, db.Location{utils.Location{5, 5}}, clock.Now().Add(duration.Duration*2)).Return(duration, nil).After(durationCall2)
	notificationData3 := make(map[string]interface{})
	notificationData3["duration"] = int64((duration.Duration * 3).Minutes())
	notificationData3["push_type"] = "driver_location_update_1"
	mockNotificationService.EXPECT().SendNotification(strconv.Itoa(4253), notificationData3, "user")

	etaResp, err := services.GetETAForTrip(&trip, currentLocation, clock)
	if err != nil {
		t.Log("Error getting ETA for trip")
		t.Log(err)
		t.Fail()
	}
	log.Debugf("eta of tripRoutes : %d", len(etaResp.TripRoutes))
	desiredResponse := services.ETAResponse{
		ID:        trip.ID,
		UpdatedAt: clock.Now(),
		TripRoutes: []services.ETATripRoute{
			services.ETATripRoute{
				ID:             421,
				Status:         "on_board",
				EmployeeUserID: 4213,
				DropoffTime:    services.NotNullTime(clock.Now().Add(duration.Duration)),
				ETAInMinutes:   duration.Duration.Minutes(),
			}, services.ETATripRoute{
				ID:             422,
				Status:         "driver_arrived",
				EmployeeUserID: 4223,
				DropoffTime:    services.NotNullTime(clock.Now().Add(duration.Duration)),
				ETAInMinutes:   duration.Duration.Minutes(),
			}, services.ETATripRoute{
				ID:             423,
				Status:         "on_board",
				EmployeeUserID: 4233,
				DropoffTime:    services.NotNullTime(clock.Now().Add(duration.Duration * 2)),
				ETAInMinutes:   duration.Duration.Minutes() * 2,
			}, services.ETATripRoute{
				ID:             424,
				Status:         "on_board",
				EmployeeUserID: 4243,
				DropoffTime:    services.NotNullTime(clock.Now().Add(duration.Duration * 2)),
				ETAInMinutes:   duration.Duration.Minutes() * 2,
			}, services.ETATripRoute{
				ID:             425,
				Status:         "driver_arrived",
				EmployeeUserID: 4253,
				DropoffTime:    services.NotNullTime(clock.Now().Add(duration.Duration * 3)),
				ETAInMinutes:   duration.Duration.Minutes() * 3,
			},
		},
	}
	assert.Equal(t, desiredResponse, *etaResp)
}

// #endregion

func checkinTripNotAllBoardWithOffset() db.Trip {
	return db.Trip{
		ID:           42,
		TripType:     db.TripTypeCheckIn,
		DriverID:     43,
		DriverUserID: 400,
		VehicleID:    23,
		Status:       "active",
		TripRoutes: []db.TripRoute{
			db.TripRoute{
				EmployeeUserID:         4212,
				ID:                     421,
				TripID:                 42,
				Status:                 "on_board",
				Trip:                   db.Trip{DriverUserID: 400},
				ScheduledRouteOrder:    1,
				ScheduledStartLocation: db.Location{utils.Location{1, 1}},
				ScheduledEndLocation:   db.Location{utils.Location{3, 3}},
			},
			db.TripRoute{
				EmployeeUserID:         4223,
				ID:                     422,
				TripID:                 42,
				Status:                 "on_board",
				Trip:                   db.Trip{DriverUserID: 400},
				ScheduledRouteOrder:    2,
				ScheduledStartLocation: db.Location{utils.Location{3, 3}},
				ScheduledEndLocation:   db.Location{utils.Location{4, 4}},
			},
			db.TripRoute{
				EmployeeUserID:         4233,
				ID:                     423,
				TripID:                 42,
				Status:                 "not_started",
				Trip:                   db.Trip{DriverUserID: 400},
				ScheduledRouteOrder:    3,
				ScheduledStartLocation: db.Location{utils.Location{4, 4}},
				ScheduledEndLocation:   db.Location{utils.Location{5, 5}},
			},
		},
	}
}

func checkinBusTripNoneBoardWithOffset() db.Trip {
	return db.Trip{
		ID:           42,
		TripType:     db.TripTypeCheckIn,
		DriverID:     43,
		DriverUserID: 400,
		VehicleID:    23,
		Status:       "active",
		TripRoutes: []db.TripRoute{
			db.TripRoute{
				EmployeeUserID:         4213,
				ID:                     421,
				TripID:                 42,
				Status:                 "not_started",
				Trip:                   db.Trip{DriverUserID: 400},
				ScheduledRouteOrder:    0,
				ScheduledStartLocation: db.Location{utils.Location{3, 3}},
				ScheduledEndLocation:   db.Location{utils.Location{2, 2}},
			},
			db.TripRoute{
				EmployeeUserID:         4223,
				ID:                     422,
				TripID:                 42,
				Status:                 "not_started",
				Trip:                   db.Trip{DriverUserID: 400},
				ScheduledRouteOrder:    0,
				ScheduledStartLocation: db.Location{utils.Location{3, 3}},
				ScheduledEndLocation:   db.Location{utils.Location{2, 2}},
			},
			db.TripRoute{
				EmployeeUserID:         4233,
				ID:                     423,
				TripID:                 42,
				Status:                 "not_started",
				Trip:                   db.Trip{DriverUserID: 400},
				ScheduledRouteOrder:    1,
				ScheduledStartLocation: db.Location{utils.Location{4, 4}},
				ScheduledEndLocation:   db.Location{utils.Location{2, 2}},
			},
			db.TripRoute{
				EmployeeUserID:         4243,
				ID:                     424,
				TripID:                 42,
				Status:                 "not_started",
				Trip:                   db.Trip{DriverUserID: 400},
				ScheduledRouteOrder:    1,
				ScheduledStartLocation: db.Location{utils.Location{4, 4}},
				ScheduledEndLocation:   db.Location{utils.Location{2, 2}},
			},
			db.TripRoute{
				EmployeeUserID:         4253,
				ID:                     425,
				TripID:                 42,
				Status:                 "not_started",
				Trip:                   db.Trip{DriverUserID: 400},
				ScheduledRouteOrder:    2,
				ScheduledStartLocation: db.Location{utils.Location{5, 5}},
				ScheduledEndLocation:   db.Location{utils.Location{2, 2}},
			},
		},
	}
}

func checkoutBusTripMixedOnBoardAndDriverArrivedWithOffset() db.Trip {
	return db.Trip{
		ID:           42,
		TripType:     db.TripTypeCheckOut,
		DriverID:     43,
		DriverUserID: 400,
		VehicleID:    23,
		Status:       "active",
		TripRoutes: []db.TripRoute{
			db.TripRoute{
				EmployeeUserID:         4213,
				ID:                     421,
				TripID:                 42,
				Status:                 "on_board",
				Trip:                   db.Trip{DriverUserID: 400},
				ScheduledRouteOrder:    0,
				ScheduledEndLocation:   db.Location{utils.Location{3, 3}},
				ScheduledStartLocation: db.Location{utils.Location{2, 2}},
			},
			db.TripRoute{
				EmployeeUserID:         4223,
				ID:                     422,
				TripID:                 42,
				Status:                 "driver_arrived",
				Trip:                   db.Trip{DriverUserID: 400},
				ScheduledRouteOrder:    0,
				ScheduledEndLocation:   db.Location{utils.Location{3, 3}},
				ScheduledStartLocation: db.Location{utils.Location{2, 2}},
			},
			db.TripRoute{
				EmployeeUserID:         4233,
				ID:                     423,
				TripID:                 42,
				Status:                 "on_board",
				Trip:                   db.Trip{DriverUserID: 400},
				ScheduledRouteOrder:    1,
				ScheduledEndLocation:   db.Location{utils.Location{4, 4}},
				ScheduledStartLocation: db.Location{utils.Location{2, 2}},
			},
			db.TripRoute{
				EmployeeUserID:         4243,
				ID:                     424,
				TripID:                 42,
				Status:                 "on_board",
				Trip:                   db.Trip{DriverUserID: 400},
				ScheduledRouteOrder:    1,
				ScheduledEndLocation:   db.Location{utils.Location{4, 4}},
				ScheduledStartLocation: db.Location{utils.Location{2, 2}},
			},
			db.TripRoute{
				EmployeeUserID:         4253,
				ID:                     425,
				TripID:                 42,
				Status:                 "driver_arrived",
				Trip:                   db.Trip{DriverUserID: 400},
				ScheduledRouteOrder:    2,
				ScheduledEndLocation:   db.Location{utils.Location{5, 5}},
				ScheduledStartLocation: db.Location{utils.Location{2, 2}},
			},
		},
	}
}

func checkinTripNotAllBoardWithOffsetWithOneOnBoard() db.Trip {
	return db.Trip{
		ID:           42,
		TripType:     db.TripTypeCheckIn,
		DriverID:     43,
		DriverUserID: 400,
		VehicleID:    23,
		Status:       "active",
		TripRoutes: []db.TripRoute{
			db.TripRoute{
				EmployeeUserID:         4212,
				ID:                     421,
				TripID:                 42,
				Status:                 "on_board",
				Trip:                   db.Trip{DriverUserID: 400},
				ScheduledRouteOrder:    1,
				ScheduledStartLocation: db.Location{utils.Location{1, 1}},
				ScheduledEndLocation:   db.Location{utils.Location{3, 3}},
			},
			db.TripRoute{
				EmployeeUserID:         4223,
				ID:                     422,
				TripID:                 42,
				Status:                 "not_started",
				Trip:                   db.Trip{DriverUserID: 400},
				ScheduledRouteOrder:    2,
				ScheduledStartLocation: db.Location{utils.Location{3, 3}},
				ScheduledEndLocation:   db.Location{utils.Location{4, 4}},
			},
			db.TripRoute{
				EmployeeUserID:         4233,
				ID:                     423,
				TripID:                 42,
				Status:                 "not_started",
				Trip:                   db.Trip{DriverUserID: 400},
				ScheduledRouteOrder:    3,
				ScheduledStartLocation: db.Location{utils.Location{4, 4}},
				ScheduledEndLocation:   db.Location{utils.Location{5, 5}},
			},
		},
	}
}

func makeAssignedCheckinTrip() *db.Trip {
	return &db.Trip{
		ID:                 42,
		TripType:           db.TripTypeCheckIn,
		DriverID:           43,
		DriverUserID:       400,
		VehicleID:          23,
		Status:             "assigned",
		ScheduledStartDate: null.TimeFrom(mockClock{}.Now()),
		TripRoutes: []db.TripRoute{
			db.TripRoute{
				EmployeeUserID:         4212,
				ID:                     421,
				TripID:                 42,
				Status:                 "not_started",
				Trip:                   db.Trip{DriverUserID: 400},
				ScheduledRouteOrder:    1,
				ScheduledStartLocation: db.Location{utils.Location{1, 1}},
				ScheduledEndLocation:   db.Location{utils.Location{3, 3}},
			},
		},
	}
}

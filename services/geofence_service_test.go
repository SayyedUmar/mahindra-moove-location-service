package services_test

import (
	"database/sql"
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MOOVE-Network/location_service/db"
	"github.com/MOOVE-Network/location_service/services"
	"github.com/MOOVE-Network/location_service/services/mocks"
	tst "github.com/MOOVE-Network/location_service/testutils"

	"github.com/golang/mock/gomock"
)

func TestSendDriverArrivingNotification(t *testing.T) {
	// setup Mocks
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockNotificationService := mocks.NewMockNotificationService(mockController)
	services.SetNotificationService(mockNotificationService)

	tripID := 123
	employeeID := 789
	driver := db.Driver{
		User: db.User{
			FirstName: sql.NullString{
				Valid:  true,
				String: "Kalpesh",
			},

			LastName: sql.NullString{
				Valid:  true,
				String: "Patel",
			},
		},
	}

	notificationMap := make(map[string]interface{})
	notificationMap["push_type"] = "driver_arriving"
	notificationMap["employee_trip_id"] = strconv.Itoa(tripID)
	notificationMap["driver_name"] = "Kalpesh Patel"

	mockNotificationService.EXPECT().SendNotification(strconv.Itoa(employeeID), notificationMap, "user").Times(1).Return(nil)
	err := services.SendDriverArrivingNotification(tripID, employeeID, &driver)
	tst.FailNowOnErr(t, err)

	driver.User.FirstName.String = "Rahul"
	driver.User.LastName.String = "Patel"

	//Driver with valid first and last names.
	notificationMap["driver_name"] = "Rahul Patel"
	mockNotificationService.EXPECT().SendNotification(strconv.Itoa(employeeID), notificationMap, "user").Times(1).Return(nil)
	err = services.SendDriverArrivingNotification(tripID, employeeID, &driver)
	tst.FailNowOnErr(t, err)

	//Driver with nil last name.
	driver.User.LastName.Valid = false
	notificationMap["driver_name"] = "Rahul"
	mockNotificationService.EXPECT().SendNotification(strconv.Itoa(employeeID), notificationMap, "user").Times(1).Return(nil)
	err = services.SendDriverArrivingNotification(tripID, employeeID, &driver)
	tst.FailNowOnErr(t, err)

	//Driver with nil first name.
	driver.User.LastName.Valid = true
	driver.User.FirstName.Valid = false
	notificationMap["driver_name"] = "Patel"
	mockNotificationService.EXPECT().SendNotification(strconv.Itoa(employeeID), notificationMap, "user").Times(1).Return(nil)
	err = services.SendDriverArrivingNotification(tripID, employeeID, &driver)
	tst.FailNowOnErr(t, err)

	//Driver with nil first and last names.
	driver.User.LastName.Valid = false
	driver.User.FirstName.Valid = false
	delete(notificationMap, "driver_name")
	mockNotificationService.EXPECT().SendNotification(strconv.Itoa(employeeID), notificationMap, "user").Times(1).Return(nil)
	err = services.SendDriverArrivingNotification(tripID, employeeID, &driver)
	tst.FailNowOnErr(t, err)

	notifError := errors.New("Error")
	mockNotificationService.EXPECT().SendNotification(strconv.Itoa(employeeID), notificationMap, "user").Times(1).Return(notifError)
	err = services.SendDriverArrivingNotification(tripID, employeeID, &driver)
	assert.Equal(t, notifError, err)
}

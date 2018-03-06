package services

import (
	"fmt"
	"strconv"

	"github.com/MOOVE-Network/location_service/db"
	logger "github.com/sirupsen/logrus"
)

//SendDriverArrivingNotification sends driver arriving notification to employee.
func SendDriverArrivingNotification(tripID int, employeeID int, driver *db.Driver) error {
	logger.Infof("Sending notification to trip %d, employeeID %d and driver %+v", tripID, employeeID, driver)
	data := make(map[string]interface{})
	data["push_type"] = "driver_arriving"
	if driver.User.FirstName.Valid && driver.User.LastName.Valid {
		data["driver_name"] = fmt.Sprintf("%s %s", driver.User.FirstName.String, driver.User.LastName.String)
	} else if driver.User.FirstName.Valid {
		data["driver_name"] = driver.User.FirstName.String
	} else if driver.User.LastName.Valid {
		data["driver_name"] = driver.User.LastName.String
	}
	data["employee_trip_id"] = strconv.Itoa(tripID)
	err := notificationService.SendNotification(strconv.Itoa(employeeID), data, "user")
	return err
}

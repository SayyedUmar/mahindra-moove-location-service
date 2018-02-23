package services

import (
	"fmt"
	"strconv"

	"github.com/MOOVE-Network/location_service/db"
	logger "github.com/sirupsen/logrus"
)

//SendDriverArrivingNotification sends driver arriving notification to employee.
func SendDriverArrivingNotification(tripID int, employeeID int, driver *db.Driver) {
	logger.Infoln("Sending notification to trip:", tripID, " employeeID:", employeeID)
	data := make(map[string]interface{})
	data["push_type"] = "driver_arriving"
	data["driver_name"] = fmt.Sprintf("%s %s", driver.User.FirstName.String, driver.User.LastName.String)
	data["employee_trip_id"] = strconv.Itoa(tripID)
	notificationService.SendNotification(strconv.Itoa(employeeID), data, "user")
}

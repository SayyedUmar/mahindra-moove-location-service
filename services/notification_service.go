package services

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/MOOVE-Network/location_service/db"
	fcm "github.com/appleboy/go-fcm"
	log "github.com/sirupsen/logrus"
)

var notificationService NotificationService

// InitNotificationService initializes the notification service with an FCM instance
func InitNotificationService(apiKey string, topicPrefix string) {
	if apiKey == "" || topicPrefix == "" {
		panic(errors.New("FCM API key and/or Topic Prefix is not set"))
	}
	ns, err := CreateFCMNotificationService(apiKey, topicPrefix)
	if err != nil {
		panic(err)
	}
	notificationService = ns
}

// SetNotificationService sets current notification service
func SetNotificationService(ns NotificationService) {
	notificationService = ns
}

// GetNotificationService returns the current NotificationService
// Please note that this can be nil
func GetNotificationService() NotificationService {
	return notificationService
}

// NotificationService provides methods to notify clients on topics
type NotificationService interface {
	SendNotification(receiverID string, data map[string]interface{}, receiverType string) error
}

// FCMNotificationService implements methods of NotificationService using the Firebase Cloud Messaging API
type FCMNotificationService struct {
	client      *fcm.Client
	topicPrefix string
}

// CreateFCMNotificationService creates a FCMNotification service using the FCM API Key and a topic prefix
func CreateFCMNotificationService(apiKey string, topicPrefix string) (*FCMNotificationService, error) {
	client, err := fcm.NewClient(apiKey)
	if err != nil {
		return nil, err
	}
	return &FCMNotificationService{client: client, topicPrefix: topicPrefix}, nil
}

// SendNotification sends the given data via FCM
func (ns *FCMNotificationService) SendNotification(receiverID string, data map[string]interface{}, receiverType string) error {
	topic := fmt.Sprintf("/topics/%s_%s_%s", ns.topicPrefix, receiverType, receiverID)
	// TODO: Check if this message is compatible for iOS
	msg := &fcm.Message{
		To:               topic,
		Data:             data,
		ContentAvailable: true,
		Priority:         "high",
	}
	res, err := ns.client.Send(msg)
	if err != nil {
		return err
	}
	if res.Error != nil {
		return errors.New("Sending notification failed")
	}
	return nil
}

// NotifyTripRouteToEmployee takes duration metrics for a trip route and sends a notification to the employee
func NotifyTripRouteToEmployee(tr *db.TripRoute, dm *DurationMetrics, offset time.Duration, ns NotificationService) {
	data := getPushNotificationData(dm, offset)
	log.Debugf("notification data for trip %d \n", tr.TripID)
	log.Debug(data)
	empID := strconv.Itoa(tr.EmployeeUserID)
	err := ns.SendNotification(empID, data, "user")
	if err != nil {
		log.Error("Unable to send notification ", err)
	}
}

// NotifyTripRouteToDriver takes duration metrics for a trip route and sends a notification to the driver
func NotifyTripRouteToDriver(tr *db.TripRoute, dm *DurationMetrics, offset time.Duration, ns NotificationService) {

	data := getPushNotificationData(dm, offset)
	log.Debugf("notification data for trip %d \n", tr.TripID)
	log.Debug(data)
	driverID := strconv.Itoa(tr.Trip.DriverUserID)
	err := ns.SendNotification(driverID, data, "driver")
	if err != nil {
		log.Error("Unable to send notification ", err)
	}
}

func getPushNotificationData(dm *DurationMetrics, offset time.Duration) map[string]interface{} {
	data := make(map[string]interface{})
	data["duration"] = int64((dm.Duration + offset).Minutes())
	data["push_type"] = "driver_location_update_1"
	return data
}

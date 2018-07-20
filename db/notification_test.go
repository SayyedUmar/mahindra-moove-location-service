package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShouldInsertTSSNotification(t *testing.T) {
	tx := createTx(t)
	defer tx.Rollback()
	tssNotification, err := CreateTripShouldStartNotification(tx, 42, 23)
	if err != nil {
		t.Fatalf("Could not create notification %v", err)
	}
	if tssNotification.ID <= 0 {
		t.Fatalf("Notification ID was not set %v", tssNotification)
	}


	t.Logf("Notification : %v", tssNotification)

	//There should not be two unresolved notifications.
	_, err = CreateTripShouldStartNotification(tx, 42, 23)
	t.Logf("Error : %v", err)
	assert.NotNil(t, err)

	oldNotifID := tssNotification.ID
	//A new notification should generate if older notification for same trip and driver is resolved.
	tx.MustExec("update notifications set resolved_status=true where id=?", tssNotification.ID)
	tssNotification, err = CreateTripShouldStartNotification(tx, 42, 23)
	if err != nil {
		t.Fatalf("Could not create notification %v", err)
	}
	if tssNotification.ID <= 0 {
		t.Fatalf("Notification ID was not set %v", tssNotification)
	}

	if tssNotification.ID == oldNotifID {
		t.Fatalf("New notification is not created when older is resolved for %v", tssNotification)
	}
}

func TestShouldInsertDOSNotification(t *testing.T) {
	tx := createTx(t)
	defer tx.Rollback()
	dosNotification, err := CreateDriverOverSpeedingNotification(tx, 42, 23)
	if err != nil {
		t.Fatalf("Could not create notification %v", err)
	}
	if dosNotification.ID <= 0 {
		t.Fatalf("Notification ID was not set %v", dosNotification)
	}
}

func TestHasUnresolvedNotification(t *testing.T) {
	tx := createTx(t)
	defer tx.Rollback()

	notification := &Notification{
		TripID:          100,
		DriverID:        101,
		Message:         "notification",
		Receiver:        int(Operator),
		Status:          0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		ResolvedStatus:  false,
		NewNotification: true,
		Sequence:        2,
	}

	notificationID, err := insertNotification(tx, notification)
	if err != nil {
		t.Fatalf("Could not create notification %v, ", notification)
	}
	notification.ID = notificationID

	hasUnresolvedNotif := HasUnresolvedNotifications(tx, notification.TripID, notification.DriverID, notification.Message)
	assert.True(t, hasUnresolvedNotif)
	assert.True(t, notification.HasUnresolved(tx))

	tx.MustExec("update notifications set resolved_status = true where id = ?", notification.ID)

	hasUnresolvedNotif = HasUnresolvedNotifications(tx, notification.TripID, notification.DriverID, notification.Message)
	assert.False(t, hasUnresolvedNotif)
	assert.False(t, notification.HasUnresolved(tx))
}

func TestHasUnresolvedNotificationsByEmployee(t *testing.T) {
	tx := createTx(t)
	defer tx.Rollback()

	notification := &Notification{
		TripID:          100,
		DriverID:        101,
		EmployeeID:      102,
		Message:         "notification",
		Receiver:        int(Operator),
		Status:          0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		ResolvedStatus:  false,
		NewNotification: true,
		Sequence:        2,
	}

	notificationID, err := insertNotification(tx, notification)
	if err != nil {
		t.Fatalf("Could not create notification %v, ", notification)
	}
	notification.ID = notificationID

	hasUnresolvedNotif := HasUnresolvedNotificationsByEmployee(tx, notification.TripID, notification.EmployeeID, notification.Message)
	assert.True(t, hasUnresolvedNotif)
	assert.True(t, notification.HasUnresolvedByEmployee(tx))

	tx.MustExec("update notifications set resolved_status = true where id = ?", notification.ID)

	hasUnresolvedNotif = HasUnresolvedNotificationsByEmployee(tx, notification.TripID, notification.EmployeeID, notification.Message)
	assert.False(t, hasUnresolvedNotif)
	assert.False(t, notification.HasUnresolvedByEmployee(tx))
}

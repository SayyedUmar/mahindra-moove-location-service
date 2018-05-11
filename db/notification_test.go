package db

import (
	"testing"
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

package db

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// Notification maps to the notifications table on the database
type Notification struct {
	ID              int64     `db:"id"`
	TripID          int       `db:"trip_id"`
	DriverID        int       `db:"driver_id"`
	Message         string    `db:"message"`
	Receiver        int       `db:"receiver"`
	Status          int       `db:"status"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
	ResolvedStatus  bool      `db:"resolved_status"`
	NewNotification bool      `db:"new_notification"`
	Sequence        int       `db:"sequence"`
}

const insertNotificationStmt = `
insert into notifications (trip_id, driver_id, message, 
													receiver, status, created_at, 
													updated_at, resolved_status,
													new_notification, sequence)
values (:trip_id, :driver_id, :message,
				:receiver, :status, :created_at,
				:updated_at, :resolved_status,
				:new_notification, :sequence)`

const TRIP_SHOULD_START = "trip_should_start"
const TSS_SEQUENCE = 3

type Receiver int

const (
	Operator Receiver = iota
	Employer
	Both
)

// CreateTripShouldStartNotification method creates the trip_should_start notification in the notifications table
func CreateTripShouldStartNotification(db *sqlx.Tx, tripID int, driverID int) (*Notification, error) {
	tssNotification := &Notification{
		TripID:          tripID,
		DriverID:        driverID,
		Message:         TRIP_SHOULD_START,
		Receiver:        int(Operator),
		Status:          0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		ResolvedStatus:  false,
		NewNotification: true,
		Sequence:        TSS_SEQUENCE,
	}
	res, err := db.NamedExec(insertNotificationStmt, tssNotification)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("Unable to get last inserted id %v", err)
	}
	tssNotification.ID = id
	return tssNotification, nil
}

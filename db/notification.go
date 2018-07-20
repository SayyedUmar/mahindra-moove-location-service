package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// Notification maps to the notifications table on the database
type Notification struct {
	ID              int64     `db:"id"`
	TripID          int       `db:"trip_id"`
	DriverID        int       `db:"driver_id"`
	EmployeeID      int       `db:"employee_id"`
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
insert into notifications (trip_id, driver_id, employee_id, message, 
													receiver, status, created_at, 
													updated_at, resolved_status,
													new_notification, sequence)
values (:trip_id, :driver_id, :employee_id, :message,
				:receiver, :status, :created_at,
				:updated_at, :resolved_status,
				:new_notification, :sequence)`

const TRIP_SHOULD_START = "trip_should_start"
const TSS_SEQUENCE = 3

const DRIVER_OVER_SPEEDING = "driver_over_speeding"
const DOS_SEQUENCE = 3

const FIRST_PICKUP_DELAYED = "first_pickup_delayed"
const FPD_SEQUENCE = 1

const SITE_ARRIVAL_DELAY = "site_arrival_delay"
const SAD_SEQUENCE = 1

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
	if !tssNotification.HasUnresolved(db) {
		notificationID, err := insertNotification(db, tssNotification)
		if err != nil {
			return nil, err
		}
		tssNotification.ID = notificationID
		return tssNotification, nil
	}
	return nil, fmt.Errorf("A trip_should_start notification already exists for trip %d and driver %d", tripID, driverID)
}

// CreateFirstPickupDelayedNotification creates a first_pickup_delayed notification
func CreateFirstPickupDelayedNotification(db *sqlx.Tx, tripID int, driverID int) (*Notification, error) {
	fpdNotification := &Notification{
		TripID:          tripID,
		DriverID:        driverID,
		Message:         FIRST_PICKUP_DELAYED,
		Receiver:        int(Operator),
		Status:          0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		ResolvedStatus:  false,
		NewNotification: true,
		Sequence:        FPD_SEQUENCE,
	}
	if !fpdNotification.HasUnresolved(db) {
		notificationID, err := insertNotification(db, fpdNotification)
		if err != nil {
			return nil, err
		}
		fpdNotification.ID = notificationID
		return fpdNotification, nil
	}
	return nil, fmt.Errorf("A first_pickup_delayed notification already exists for trip %d and driver %d", tripID, driverID)
}

// CreateSiteArrivalDelayNotification creates a site_arrival_delay notification
func CreateSiteArrivalDelayNotification(db *sqlx.Tx, tripID, driverID, employeeID int) (*Notification, error) {
	sadNotification := &Notification{
		TripID:          tripID,
		DriverID:        driverID,
		EmployeeID:      employeeID,
		Message:         SITE_ARRIVAL_DELAY,
		Receiver:        int(Operator),
		Status:          0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		ResolvedStatus:  false,
		NewNotification: true,
		Sequence:        SAD_SEQUENCE,
	}
	if !sadNotification.HasUnresolvedByEmployee(db) {
		notificationID, err := insertNotification(db, sadNotification)
		if err != nil {
			return nil, err
		}
		sadNotification.ID = notificationID
		return sadNotification, nil
	}
	return nil, fmt.Errorf("A site_arrival_delay notification already exists for trip %d and employee %d", tripID, employeeID)
}

// CreateDriverOverSpeedingNotification method creates the driver_over_speeding notification in the notifications table
func CreateDriverOverSpeedingNotification(db *sqlx.Tx, tripID int, driverID int) (*Notification, error) {
	dosNotification := &Notification{
		TripID:          tripID,
		DriverID:        driverID,
		Message:         DRIVER_OVER_SPEEDING,
		Receiver:        int(Operator),
		Status:          0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		ResolvedStatus:  false,
		NewNotification: true,
		Sequence:        DOS_SEQUENCE,
	}
	notificationID, err := insertNotification(db, dosNotification)
	if err != nil {
		return nil, err
	}
	dosNotification.ID = notificationID
	return dosNotification, nil
}

// HasUnresolved returns true if the given notifications message, trip_id, driver_id
// combination has any unresolved notifications
func (n *Notification) HasUnresolved(db *sqlx.Tx) bool {
	return HasUnresolvedNotifications(db, n.TripID, n.DriverID, n.Message)
}

// HasUnresolvedByEmployee returns true if the given notifications message, trip_id, employee_id
// combination has any unresolved notifications
func (n *Notification) HasUnresolvedByEmployee(db *sqlx.Tx) bool {
	return HasUnresolvedNotificationsByEmployee(db, n.TripID, n.EmployeeID, n.Message)
}

// HasUnresolvedNotifications returns true if trip has unresolved notifications for a driver
func HasUnresolvedNotifications(db *sqlx.Tx, tripID int, driverID int, message string) bool {
	checkNotificationQuery := `select id from notifications 
	where trip_id=? and driver_id=? and message=?
	and resolved_status=false`
	row := db.QueryRowx(checkNotificationQuery, tripID, driverID, message)
	var id int
	err := row.Scan(&id)
	if err != nil && err == sql.ErrNoRows {
		return false
	}
	if err == nil {
		return true
	}
	return false
}

// HasUnresolvedNotificationsByEmployee returns true if trip has unresolved notifications for an employee
func HasUnresolvedNotificationsByEmployee(db *sqlx.Tx, tripID int, employeeID int, message string) bool {
	checkNotificationQuery := `select id from notifications 
	where trip_id=? and employee_id=? and message=?
	and resolved_status=false`
	row := db.QueryRowx(checkNotificationQuery, tripID, employeeID, message)
	var id int
	err := row.Scan(&id)
	if err != nil && err == sql.ErrNoRows {
		return false
	}
	if err == nil {
		return true
	}
	return false
}

func insertNotification(db *sqlx.Tx, notification *Notification) (int64, error) {
	res, err := db.NamedExec(insertNotificationStmt, notification)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Unable to get last inserted id %v", err)
	}
	return id, nil
}

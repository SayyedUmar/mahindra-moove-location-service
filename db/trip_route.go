package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	null "gopkg.in/guregu/null.v3"
)

// OnBoardStatuses contains of list of statuses
// when employees are either on board or have been cancelled
var OnBoardStatuses = map[string]bool{
	"canceled":       true,
	"on_board":       true,
	"missed":         true,
	"driver_arrived": true,
}

// TripRoute represents the database structure of TripRoute
type TripRoute struct {
	ID                     int            `db:"id"`
	TripID                 int            `db:"trip_id"`
	Status                 string         `db:"status"`
	ScheduledRouteOrder    int            `db:"scheduled_route_order"`
	ScheduledStartLocation Location       `db:"scheduled_start_location"`
	ScheduledEndLocation   Location       `db:"scheduled_end_location"`
	EmployeeID             int            `db:"employee_id"`
	EmployeeUserID         int            `db:"employee_user_id"`
	BusStopName            sql.NullString `db:"bus_stop_name"`
	PickUpTime             null.Time      `db:"pick_up_time"`
	DropOffTime            null.Time      `db:"drop_off_time"`
	Date                   null.Time      `db:"date"`
	Trip                   Trip
}

const tripRoutesByIDQuery = `
	select tr.id, tr.trip_id, tr.status, u.id as employee_user_id, et.employee_id,
	tr.scheduled_route_order, tr.scheduled_start_location, tr.scheduled_end_location,
	tr.bus_stop_name, et.date
	from trip_routes tr
	join employee_trips et on et.id = tr.employee_trip_id
	join employees e on e.id = et.employee_id
	join users u on u.entity_id=e.id and u.entity_type="Employee"
	where tr.id=?`

const updateGeofenceDriverArriveQuery = `
	update trip_routes set geofence_driver_arrived_date=?,
	geofence_driver_arrived_location=? where id=?
`
const updateGeofenceCompletedQuery = `
	update trip_routes set geofence_completed_date=?,
	geofence_completed_location=? where id=?
`

// IsMissedOrCanceled is returns true if employee has canceled the trip or didn't show up for pickup. false otherwise.
func (tr *TripRoute) IsMissedOrCanceled() bool {
	return tr.Status == "canceled" || tr.Status == "missed"
}

// IsOnBoard is considered on board if he is on board or driver has arrived
func (tr *TripRoute) IsOnBoard() bool {
	return tr.Status == "on_board" || tr.Status == "driver_arrived"
}

//IsNotStarted checks if the employee trip is yet to start.
func (tr *TripRoute) IsNotStarted() bool {
	return tr.Status == "not_started"
}

//IsTripRouteNotStarted check is for given tripRoute id trip route is in started state or not.
func IsTripRouteNotStarted(db sqlx.Queryer, id int) (bool, error) {
	tr, err := getTripRouteByID(db, id)

	if err != nil {
		return false, err
	}

	return tr.IsNotStarted(), nil
}

//getTripRouteByID returns TripRoute for given tripRoute id
//Caution: this does not load TripRoute.Trip
func getTripRouteByID(db sqlx.Queryer, id int) (*TripRoute, error) {
	rows, err := db.Query(tripRoutesByIDQuery, id)
	if err != nil {
		panic(err)
	}
	var tr TripRoute
	// err := row.StructScan(&tr)

	/*
		ID                     int            `db:"id"`
		TripID                 int            `db:"trip_id"`
		Status                 string         `db:"status"`
		ScheduledRouteOrder    int            `db:"scheduled_route_order"`
		ScheduledStartLocation Location       `db:"scheduled_start_location"`
		ScheduledEndLocation   Location       `db:"scheduled_end_location"`
		EmployeeID             int            `db:"employee_id"`
		EmployeeUserID         int            `db:"employee_user_id"`
		BusStopName            sql.NullString `db:"bus_stop_name"`
		PickUpTime             null.Time      `db:"pick_up_time"`
		DropOffTime            null.Time      `db:"drop_off_time"`
		Date                   null.Time      `db:"date"`
		Trip                   Trip
	*/

	for rows.Next() {
		err := rows.Scan(&tr.ID,
			&tr.TripID,
			&tr.Status,
			&tr.ScheduledRouteOrder,
			&tr.ScheduledStartLocation,
			&tr.ScheduledEndLocation,
			// &tr.EmployeeID,
			&tr.EmployeeUserID,
			&tr.BusStopName,
			// &tr.PickUpTime,
			// &tr.DropOffTime,
			&tr.Date)
		// &tr.Trip)
		if err != nil {
			fmt.Println("Error during stuct scan of trip route", err)
		}
		// fmt.Println("================================", tl)
	}

	// if err != nil {
	// 	fmt.Println("Error during stuct scan of trip route")
	// 	return nil, err
	// }

	return &tr, nil
}

func (tr *TripRoute) IsSiteArrivalDelayed(db sqlx.Queryer, newDropTime time.Time) bool {
	if !tr.Date.Valid {
		return false
	}
	val, err := ConfigGetTripDelayNotification(db)
	if err != nil {
		log.Errorf("Unable to fetch TripDelayedNotification configuration - %s. Defaulting to %d", err, val)
	}
	return newDropTime.Sub(tr.Date.Time) > time.Duration(val)*time.Minute
}
func (tr *TripRoute) TriggerSiteArrivalDelayNotification(db *sqlx.Tx, newDropTime time.Time) error {
	if !tr.IsSiteArrivalDelayed(db, newDropTime) {
		return nil
	}
	_, err := CreateSiteArrivalDelayNotification(db, tr.TripID, tr.Trip.DriverID, tr.EmployeeID)
	return err
}

func (tr *TripRoute) UpdateDriverArrivedGeofenceInfo(db sqlx.Execer, location Location, time time.Time) error {
	currentLocation, err := location.ToYaml()
	if err != nil {
		log.Errorf("Can not convert location to yaml %+v", location)
		return err
	}
	_, err = db.Exec(updateGeofenceDriverArriveQuery, time, currentLocation, tr.ID)
	return err
}

func (tr *TripRoute) UpdateCompletedGeofenceInfo(db *sqlx.Tx, location Location, time time.Time) error {
	currentLocation, err := location.ToYaml()
	if err != nil {
		log.Errorf("Can not convert location to yaml %+v", location)
		return err
	}
	_, err = db.Exec(updateGeofenceCompletedQuery, time, currentLocation, tr.ID)
	return err
}

func SaveEta(db sqlx.Execer, trId int, pickUpTime null.Time, dropOffTime null.Time) error {
	if !pickUpTime.IsZero() && !dropOffTime.IsZero() {
		_, err := db.Exec(`update trip_routes
							set pick_up_time=?, drop_off_time=?
							where id=?`, pickUpTime, dropOffTime, trId)
		return err
	}
	if !pickUpTime.IsZero() && dropOffTime.IsZero() {
		_, err := db.Exec(`update trip_routes
							set pick_up_time=?
							where id=?`, pickUpTime, trId)
		return err
	}
	if pickUpTime.IsZero() && !dropOffTime.IsZero() {
		_, err := db.Exec(`update trip_routes
							set drop_off_time=?
							where id=?`, dropOffTime, trId)
		return err
	}
	return nil
}

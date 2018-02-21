package db

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
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
	EmployeeUserID         int            `db:"employee_user_id"`
	BusStopName            sql.NullString `db:"bus_stop_name"`
	Trip                   Trip
}

var tripRoutesByIDQuery = `
	select tr.id, tr.trip_id, tr.status, u.id as employee_user_id,
	tr.scheduled_route_order, tr.scheduled_start_location, tr.scheduled_end_location,
	tr.bus_stop_name
	from trip_routes tr
	join employee_trips et on et.id = tr.employee_trip_id
	join employees e on e.id = et.employee_id
	join users u on u.entity_id=e.id and u.entity_type="Employee"
	where tr.id=?`

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
	row := db.QueryRowx(tripRoutesByIDQuery, id)
	var tr TripRoute
	err := row.StructScan(&tr)

	if err != nil {
		fmt.Println("Error during stuct scan of trip route")
		return nil, err
	}

	return &tr, nil
}

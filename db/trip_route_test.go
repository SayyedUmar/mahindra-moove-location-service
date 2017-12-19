package db

import (
	"time"

	"github.com/jmoiron/sqlx"
)

func createTripRoutes(t *Trip, trs []TripRoute, tx *sqlx.Tx) (*Trip, error) {
	stmt, err := tx.Preparex(`insert into trip_routes(trip_id, status, employee_trip_id, 
						 scheduled_route_order, scheduled_start_location, scheduled_end_location)
						 values (?,?,?,?,?,?)`)
	if err != nil {
		return nil, err
	}
	for routeOrder, tr := range trs {
		eid := getEmployeeID(tx, users[routeOrder%len(employees)].Email)
		etID := createEmployeeTrip(tx, eid)
		_, err := stmt.Exec(t.ID, tr.Status, etID, routeOrder, tr.ScheduledStartLocation, tr.ScheduledEndLocation)
		if err != nil {
			return nil, err
		}
	}
	err = t.LoadTripRoutes(tx, true)
	if err != nil {
		return nil, err
	}
	return t, nil
}
func createEmployeeTrip(tx *sqlx.Tx, employeeID int) int {
	res, err := tx.Exec(`
		insert into employee_trips(employee_id, created_at, updated_at)
		values(?,?,?)
		`, employeeID, time.Now(), time.Now())
	if err != nil {
		panic(err)
	}
	etID, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}
	return int(etID)
}

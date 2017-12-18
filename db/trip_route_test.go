package db

import (
	"github.com/jmoiron/sqlx"
)

func createTripRoutes(t *Trip, trs []TripRoute, tx *sqlx.Tx) (*Trip, error) {
	stmt, err := tx.Preparex(`insert into trip_routes(trip_id, status, 
						 scheduled_route_order, scheduled_start_location, scheduled_end_location)
						 values (?,?,?,?,?)`)
	if err != nil {
		return nil, err
	}
	for routeOrder, tr := range trs {
		_, err := stmt.Exec(t.ID, tr.Status, routeOrder, tr.ScheduledStartLocation, tr.ScheduledEndLocation)
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

package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	// TripTypeCheckIn is a checkin trip with value 0
	TripTypeCheckIn = iota
	// TripTypeCheckOut is a checkin trip with value 0
	TripTypeCheckOut = iota
)

// Trip structure maps to the trips table
type Trip struct {
	ID             int    `db:"id"`
	TripType       int    `db:"trip_type"`
	DriverID       int    `db:"driver_id"`
	VehicleID      int    `db:"vehicle_id"`
	Status         string `db:"status"`
	TripRoutes     []TripRoute
	isRoutesLoaded bool
}

// GetTripByID returns a trip if found otherwise returns an error
func GetTripByID(db sqlx.Queryer, id int) (*Trip, error) {
	row := db.QueryRowx("select id, trip_type, driver_id, vehicle_id, status from trips where id=?", id)
	var t Trip
	err := row.StructScan(&t)
	if err != nil {
		fmt.Println("Error during stuct scan")
		return nil, err
	}
	err = t.LoadTripRoutes(db, false)
	if err != nil {
		fmt.Println("Error during loading trip routes")
		return nil, err
	}
	return &t, nil
}

// LoadTripRoutes loads / refreshes the trip routes for a given trip
func (t *Trip) LoadTripRoutes(db sqlx.Queryer, force bool) error {
	if t.isRoutesLoaded && !force {
		return nil
	}
	rows, err := db.Queryx("select id, trip_id, status, scheduled_route_order, scheduled_start_location, scheduled_end_location from trip_routes where trip_id=? order by scheduled_route_order asc", t.ID)
	if err != nil {
		return err
	}
	var tripRoutes []TripRoute
	for rows.Next() {
		var tr TripRoute
		err := rows.StructScan(&tr)
		if err != nil {
			return err
		}
		tr.Trip = t
		tripRoutes = append(tripRoutes, tr)
	}
	t.isRoutesLoaded = true
	t.TripRoutes = tripRoutes
	return nil
}

func (t *Trip) IsCheckIn() bool {
	return t.TripType == TripTypeCheckIn
}

// AllCheckedIn returns if all the employees are on board
// or if the driver has arrived at the location
func (t *Trip) AllCheckedIn() bool {
	allOnBoard := false
	for _, tr := range t.TripRoutes {
		_, found := OnBoardStatuses[tr.Status]
		if !found {
			break
		}
		allOnBoard = true
	}
	if allOnBoard && t.Status != "canceled" {
		return true
	}
	return false
}

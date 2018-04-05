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

var tripRoutesForTripQuery = `
	select tr.id, tr.trip_id, tr.status, u.id as employee_user_id,
	tr.scheduled_route_order, tr.scheduled_start_location, tr.scheduled_end_location,
	tr.pick_up_time, tr.drop_off_time
	from trip_routes tr
	join employee_trips et on et.id = tr.employee_trip_id
	join employees e on e.id = et.employee_id
	join users u on u.entity_id=e.id and u.entity_type="Employee"
	where tr.trip_id=? order by scheduled_route_order asc`

var tripRoutesForTripIDsQuery = `
	select tr.id, tr.trip_id, tr.status, u.id as employee_user_id,
	tr.scheduled_route_order, tr.scheduled_start_location, tr.scheduled_end_location,
	tr.pick_up_time, tr.drop_off_time
	from trip_routes tr
	join employee_trips et on et.id = tr.employee_trip_id
	join employees e on e.id = et.employee_id
	join users u on u.entity_id=e.id and u.entity_type="Employee"
	where tr.trip_id in (?) order by tr.trip_id, scheduled_route_order asc`

var tripByIDQuery = `
	select t.id, t.trip_type, t.driver_id, u.id as driver_user_id, t.vehicle_id, t.status 
	from trips t 
	join drivers d on d.id = t.driver_id
	join users u on u.entity_id=d.id and u.entity_type="Driver"
	where t.id=?`

var tripsByStatusQuery = `
	select t.id, t.trip_type, t.driver_id, u.id as driver_user_id, t.vehicle_id, t.status 
	from trips t 
	join drivers d on d.id = t.driver_id
	join users u on u.entity_id=d.id and u.entity_type="Driver"
	where t.status=?`

var updateActualMileageStmt = `
	update trips set actual_mileage=? where id=?
`

// Trip structure maps to the trips table
type Trip struct {
	ID             int    `db:"id"`
	TripType       int    `db:"trip_type"`
	DriverID       int    `db:"driver_id"`
	DriverUserID   int    `db:"driver_user_id"`
	VehicleID      int    `db:"vehicle_id"`
	Status         string `db:"status"`
	ActualMileage  int    `db:"actual_mileage"`
	TripRoutes     []TripRoute
	isRoutesLoaded bool
}

// GetTripByID returns a trip if found otherwise returns an error
func GetTripByID(db sqlx.Queryer, id int) (*Trip, error) {
	row := db.QueryRowx(tripByIDQuery, id)
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

func (t *Trip) SetActualMileage(db sqlx.Execer, mileage int) error {
	_, err := db.Exec(updateActualMileageStmt, mileage, t.ID)
	return err
}

// LoadTripRoutes loads / refreshes the trip routes for a given trip
func (t *Trip) LoadTripRoutes(db sqlx.Queryer, force bool) error {
	if t.isRoutesLoaded && !force {
		return nil
	}
	rows, err := db.Queryx(tripRoutesForTripQuery, t.ID)
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
		tr.Trip = *t
		tripRoutes = append(tripRoutes, tr)
	}
	t.isRoutesLoaded = true
	t.TripRoutes = tripRoutes
	return nil
}

// GetTripsByStatus loads trips with a given status and also eager loads trip routes along with it
func GetTripsByStatus(db RebindQueryer, status string) ([]*Trip, error) {
	var trips []*Trip
	tripMap := make(map[int]*Trip)
	rows, err := db.Queryx(tripsByStatusQuery, status)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var t Trip
		err = rows.StructScan(&t)
		if err != nil {
			return nil, err
		}
		tripMap[t.ID] = &t
	}
	var tripIDs []int
	for id, t := range tripMap {
		tripIDs = append(tripIDs, id)
		trips = append(trips, t)
	}

	q, args, err := sqlx.In(tripRoutesForTripIDsQuery, tripIDs)
	if err != nil {
		return nil, err
	}
	q = db.Rebind(q)
	trRows, err := db.Queryx(q, args...)
	if err != nil {
		return nil, err
	}
	for trRows.Next() {
		var tr TripRoute
		err := trRows.StructScan(&tr)
		if err != nil {
			return nil, err
		}
		tripMap[tr.TripID].TripRoutes = append(tripMap[tr.TripID].TripRoutes, tr)
	}
	return trips, nil
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
			return false
		}
	}
	allOnBoard = true
	if allOnBoard && t.Status != "canceled" {
		return true
	}
	return false
}

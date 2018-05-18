package db

import (
	"fmt"

	"gopkg.in/guregu/null.v3"

	"time"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

const (
	// TripTypeCheckIn is a checkin trip with value 0
	TripTypeCheckIn = iota
	// TripTypeCheckOut is a checkin trip with value 0
	TripTypeCheckOut = iota
)

var tripRoutesForTripQuery = `
	select tr.id, tr.trip_id, tr.status, u.id as employee_user_id, et.employee_id,
	tr.scheduled_route_order, tr.scheduled_start_location, tr.scheduled_end_location,
	tr.bus_stop_name, tr.pick_up_time, tr.drop_off_time, et.date
	from trip_routes tr
	join employee_trips et on et.id = tr.employee_trip_id
	join employees e on e.id = et.employee_id
	join users u on u.entity_id=e.id and u.entity_type="Employee"
	where tr.trip_id=? order by scheduled_route_order asc`

var tripRoutesForTripIDsQuery = `
	select tr.id, tr.trip_id, tr.status, u.id as employee_user_id, et.employee_id,
	tr.scheduled_route_order, tr.scheduled_start_location, tr.scheduled_end_location,
	tr.bus_stop_name, tr.pick_up_time, tr.drop_off_time, et.date
	from trip_routes tr
	join employee_trips et on et.id = tr.employee_trip_id
	join employees e on e.id = et.employee_id
	join users u on u.entity_id=e.id and u.entity_type="Employee"
	where tr.trip_id in (?) order by tr.trip_id, scheduled_route_order asc`

var tripByIDQuery = `
	select t.id, t.trip_type, t.driver_id, u.id as driver_user_id, t.vehicle_id, t.status, t.scheduled_date 
	from trips t 
	join drivers d on d.id = t.driver_id
	join users u on u.entity_id=d.id and u.entity_type="Driver"
	where t.id=?`

var tripsByStatusQuery = `
	select t.id, t.trip_type, t.driver_id, u.id as driver_user_id, t.vehicle_id, t.status, t.scheduled_date 
	from trips t 
	join drivers d on d.id = t.driver_id
	join users u on u.entity_id=d.id and u.entity_type="Driver"
	where t.status in (?)`

var updateActualMileageStmt = `
	update trips set actual_mileage=? where id=?
`

var updateDriverShouldStartTripTimeAndLocationStmt = `
	update trips set driver_should_start_trip_time=?, driver_should_start_trip_location=?, 
	driver_should_start_trip_timestamp=? where id=?
`
var getDriverShouldStartTripLocationQuery = `
	select driver_should_start_trip_location from trips where id=?
`

// Trip structure maps to the trips table
type Trip struct {
	ID                 int       `db:"id"`
	TripType           int       `db:"trip_type"`
	DriverID           int       `db:"driver_id"`
	DriverUserID       int       `db:"driver_user_id"`
	VehicleID          int       `db:"vehicle_id"`
	Status             string    `db:"status"`
	ActualMileage      int       `db:"actual_mileage"`
	ScheduledStartDate null.Time `db:"scheduled_date"`
	TripRoutes         []TripRoute
	isRoutesLoaded     bool
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

// HasStarted returns true if the trip is in active state.
func (t *Trip) HasStarted() bool {
	return t.Status == "active"
}

// IsTerminal returns true if the trip's state is cancelled or completed
func (t *Trip) IsTerminal() bool {
	return t.Status == "completed" || t.Status == "cancelled"
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

// GetTripsByStatuses loads trips with given statuses and also eager loads trip routes along with it
func GetTripsByStatuses(db RebindQueryer, statuses ...string) ([]*Trip, error) {
	var trips []*Trip
	tripMap := make(map[int]*Trip)

	q, args, err := sqlx.In(tripsByStatusQuery, statuses)
	if err != nil {
		return nil, err
	}
	q = db.Rebind(q)
	rows, err := db.Queryx(q, args...)
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

	q, args, err = sqlx.In(tripRoutesForTripIDsQuery, tripIDs)
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

// GetFirstTripRoute returns the first trip route based on ScheduledRouteOrder
func (t *Trip) GetFirstTripRoute() *TripRoute {
	for _, tr := range t.TripRoutes {
		if tr.ScheduledRouteOrder == 0 {
			return &tr
		}
	}
	return nil
}

func (t *Trip) IsFirstPickupDelayed(db sqlx.Queryer, newPickupTime time.Time) bool {
	if !t.ScheduledStartDate.Valid {
		return false
	}
	val, err := ConfigGetTripDelayNotification(db)
	if err != nil {
		log.Errorf("Unable to fetch TripDelayedNotification configuration - %s. Defaulting to %d", err, val)
	}
	return newPickupTime.Sub(t.ScheduledStartDate.Time) > time.Duration(val)*time.Minute
}

func (t *Trip) TriggerFirstPickupDelayedNotification(db *sqlx.DB, pickupTime time.Time) error {
	if !t.IsFirstPickupDelayed(db, pickupTime) {
		return nil
	}
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("Could not create transaction %v", err)
	}
	CreateFirstPickupDelayedNotification(tx, t.ID, t.DriverID)
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Error committing transaction for First Pickup delayed %d", t.ID)
	}
	return nil
}
func (t *Trip) GetTripRoute(trID int) *TripRoute {
	for _, tr := range t.TripRoutes {
		if tr.ID == trID {
			return &tr
		}
	}
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
			return false
		}
	}
	allOnBoard = true
	if allOnBoard && t.Status != "canceled" {
		return true
	}
	return false
}

//UpdateDriverShouldStartTripTimeAndLocation updates driver_should_start_trip_time, driver_should_start_trip_location
//and driver_should_start_trip_timestamp for a given trip t.
func (t *Trip) UpdateDriverShouldStartTripTimeAndLocation(db sqlx.Execer, scheduledTime time.Time, location Location, calculationTime time.Time) error {
	_, err := db.Exec(updateDriverShouldStartTripTimeAndLocationStmt, scheduledTime, location, calculationTime, t.ID)
	return err
}

func (t *Trip) GetDriverShouldStartTripLocation(db sqlx.Queryer) (*Location, error) {
	row := db.QueryRowx(getDriverShouldStartTripLocationQuery, t.ID)
	var location Location
	err := row.Scan(&location)
	if err != nil {
		return nil, err
	}
	return &location, nil
}

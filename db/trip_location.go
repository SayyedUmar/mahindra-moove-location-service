package db

import (
	"time"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

// TripLocation represents a record in the trip_locations table
type TripLocation struct {
	ID        int       `db:"id"`
	TripID    int       `db:"trip_id"`
	Location  Location  `db:"location"`
	Time      time.Time `db:"time"`
	Speed     string    `db:"speed"`
	Distance  int       `db:"distance"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

const GetTripLocationsByTripIDQuery = `
	select id, trip_id, location, time, speed, distance, created_at, updated_at from trip_locations where trip_id=? order by time asc
`

// InsertTripLocationStatement returns a prepared statement that can be used to create
// a TripLocation struct in the database
func InsertTripLocationStatement(db Preparex) *sqlx.NamedStmt {
	stmt, err := db.PrepareNamed(`insert into
		trip_locations(trip_id, location, time, speed, distance, created_at, updated_at)
		values(:trip_id, :location, :time, :speed, :distance, :created_at, :updated_at)`)
	if err != nil {
		log.Panic(err)
	}
	return stmt
}

// Save saves the trip location to the database
func (tl *TripLocation) Save(db NamedExecer) error {
	_, err := db.Exec(tl)
	return err
}

// LatestTripLocation returns the latest trip location of the trip
func LatestTripLocation(q sqlx.Queryer, tripID int) (*TripLocation, error) {
	var tl TripLocation
	row := q.QueryRowx("select * from trip_locations where trip_id=? order by id desc limit 1", tripID)
	err := row.StructScan(&tl)
	return &tl, err
}

func GetTripLocationsByTrip(q Selectable, tripID int) ([]TripLocation, error) {
	var tripLocations []TripLocation
	err := q.Select(&tripLocations, GetTripLocationsByTripIDQuery, tripID)
	return tripLocations, err
}

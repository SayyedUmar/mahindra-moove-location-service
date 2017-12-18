package db

import (
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"time"
)

type TripLocation struct {
	TripId    int       `db:"trip_id"`
	Location  string    `db:"location"`
	Time      time.Time `db:"time"`
	Speed     string    `db:"speed"`
	Distance  int       `db:"distance"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func InsertTripLocationStatement(db Preparex) *sqlx.NamedStmt {
	stmt, err := db.PrepareNamed(`insert into
													 trip_locations(trip_id, location, time,
													 								speed, distance,
																					created_at, updated_at)
													 values(:trip_id, :location,
													 				:time, :speed, :distance,
																	:created_at, :updated_at)`)
	if err != nil {
		log.Panic(err)
	}
	return stmt
}

func (tl *TripLocation) Save(db NamedExecer) error {
	_, err := db.Exec(tl)
	return err
}

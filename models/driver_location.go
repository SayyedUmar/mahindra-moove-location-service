package models

import (
	"time"

	"github.com/jmoiron/sqlx"
	null "gopkg.in/guregu/null.v3"
)

type DriverLocation struct {
	ID         string      `db:"ID"`
	RecordedAt time.Time   `db:"recorded_at"`
	TripID     null.Int    `db:"trip_id"`
	UserID     null.String `db:"user_id"`
	Location   Point       `db:"location"`
	Distance   int         `db:"distance"`
	Speed      float64     `db:"speed"`
	Accuracy   float64     `db:"accuracy"`
	CreatedAt  time.Time   `db:"created_at"`
}

const InsertDriverLocationStmt = `
	insert into 
		driver_locations(recorded_at, trip_id, user_id, location, 
										 distance, speed, accuracy, created_at)
		values(:recorded_at, :trip_id, :user_id, :location,
					 :distance, :speed, :accuracy, :created_at)
`

func DriverLocationPrepareInsertStmt(db *sqlx.DB) (*sqlx.NamedStmt, error) {
	return db.PrepareNamed(InsertDriverLocationStmt)
}

func (dl *DriverLocation) Insert(stmt *sqlx.NamedStmt) error {
	_, err := stmt.Exec(dl)
	return err
}

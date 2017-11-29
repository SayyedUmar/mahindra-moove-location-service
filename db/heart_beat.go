package db

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type HeartBeat struct {
	UserID    int
	Lat       float64
	Lng       float64
	UpdatedAt time.Time
}

func (hb *HeartBeat) Save(db sqlx.Execer) error {
	if hb.Lat == 0 && hb.Lng == 0 {
		return hb.SaveOnlyTime(db)
	} else {
		return hb.SaveWithLatLng(db)
	}
}

func (hb *HeartBeat) SaveWithLatLng(db sqlx.Execer) error {
	_, err := db.Exec(`update users set
						lat = ?, lng = ?, updated_at=?
						where id=?`, hb.Lat, hb.Lng, hb.UpdatedAt, hb.UserID)
	return err
}

func (hb *HeartBeat) SaveOnlyTime(db sqlx.Execer) error {
	_, err := db.Exec(`update users set updated_at=?  where id=?`, hb.UpdatedAt, hb.UserID)
	return err
}

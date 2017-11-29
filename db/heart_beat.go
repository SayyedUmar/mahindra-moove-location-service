package db

import (
	"github.com/MOOVE-Network/location_service/utils"
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
	currentLocation, err := utils.ToLocation(hb.Lat, hb.Lng)
	if err != nil {
		return hb.SaveOnlyTime(db)
	}
	_, err = db.Exec(`update users set
						last_active_time=?, current_location=?
						where id=?`, hb.UpdatedAt, currentLocation, hb.UserID)
	return err
}

func (hb *HeartBeat) SaveOnlyTime(db sqlx.Execer) error {
	_, err := db.Exec(`update users set last_active_time=? where id=?`, hb.UpdatedAt, hb.UserID)
	return err
}

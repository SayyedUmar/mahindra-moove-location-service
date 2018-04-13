package db

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

//Driver structure that maps to drivers table
type Driver struct {
	ID     int            `db:"id"`
	Status sql.NullString `db:"status"`
	User   User
}

var getDriverByIDQuery = `
	select id, status
	from drivers 
	where id=? `

var getDriverByTripIDQuery = `
	select d.id, d.status
	from drivers d
	join trips t on t.driver_id = d.id
	where t.id = ?
`

var getUserByDriverIDQuery = `
	select u.id, u.f_name, u.m_name, u.l_name, u.entity_id, u.entity_type
	from users u
	join drivers d on d.id = u.entity_id and u.entity_type = "Driver"
	where d.id=?`

var getLastDriverLocationQuery = `select current_location from users where id=?`

//GetDriverByID retuns Driver struct along with user for a give driver id if found otherwise will return error
func GetDriverByID(db sqlx.Queryer, driverID int) (*Driver, error) {
	row := db.QueryRowx(getDriverByIDQuery, driverID)
	return loadDiver(db, row)
}

func loadDiver(db sqlx.Queryer, row *sqlx.Row) (*Driver, error) {
	var driver Driver
	err := row.StructScan(&driver)
	if err != nil {
		fmt.Println("Error during loading driver")
		return nil, err
	}

	row = db.QueryRowx(getUserByDriverIDQuery, driver.ID)
	var user User
	err = row.StructScan(&user)
	if err != nil {
		fmt.Println("Error during loading user for driver")
		return nil, err
	}

	driver.User = user

	return &driver, nil
}

//GetDriverByTripID returns driver struct along with User for given trip id if found otherwise will return error.
func GetDriverByTripID(db sqlx.Queryer, tripID int) (*Driver, error) {
	row := db.QueryRowx(getDriverByTripIDQuery, tripID)
	return loadDiver(db, row)
}

//DriverLocation Returns last known location for the driver.
func DriverLocation(db sqlx.Queryer, driverUserID int) (*Location, error) {
	row := db.QueryRowx(getLastDriverLocationQuery, driverUserID)
	type LocationWrapper struct {
		Location Location `db:"current_location"`
	}
	var locationWrapper LocationWrapper
	err := row.StructScan(&locationWrapper)
	if err != nil {
		return nil, err
	}
	return &locationWrapper.Location, nil
}

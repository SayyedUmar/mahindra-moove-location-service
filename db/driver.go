package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

//Driver structure that maps to drivers table
type Driver struct {
	ID     int    `db:"id"`
	Status string `db:"status"`
	User   User
}

var getDriverByIDQuery = `
	select id, status
	from drivers 
	where id=? `

var getUserByDriverIDQuery = `
	select u.id, u.f_name, u.m_name, u.l_name, u.entity_id, u.entity_type
	from users u
	join drivers d on d.id = u.entity_id and u.entity_type = "Driver"
	where d.id=?`

//GetDriverByID retuns Driver struct along with user for a give driver id if found otherwise will return error
func GetDriverByID(db sqlx.Queryer, driverID int) (*Driver, error) {
	row := db.QueryRowx(getDriverByIDQuery, driverID)
	var driver Driver
	err := row.StructScan(&driver)
	if err != nil {
		fmt.Println("Error during loading driver")
		return nil, err
	}

	row = db.QueryRowx(getUserByDriverIDQuery, driverID)
	var user User
	err = row.StructScan(&user)
	if err != nil {
		fmt.Println("Error during loading user for driver")
		return nil, err
	}

	driver.User = user

	return &driver, nil
}

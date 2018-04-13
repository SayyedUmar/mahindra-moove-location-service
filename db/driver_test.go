package db

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	tst "github.com/MOOVE-Network/location_service/testutils"
	"github.com/MOOVE-Network/location_service/utils"
	"github.com/icrowley/fake"
	"github.com/jmoiron/sqlx"
)

func createDriver(tx *sqlx.Tx, status string) (*Driver, error) {
	driversInsertQuery := "insert into drivers (status, created_at, updated_at) values(?,?,?)"
	usersInsertQuery := `
		insert into users (f_name, l_name, email, encrypted_password, uid, sign_in_count, created_at, updated_at, provider,  entity_id, entity_type) 
					values(?,		?,		?    , ?                 , ?  , ?            , ?         , ?         , ?       ,  ?        , ?          )`

	result := tx.MustExec(driversInsertQuery, status, time.Now(), time.Now())
	driverID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	driver := Driver{
		ID: int(driverID),
		Status: sql.NullString{
			String: status,
			Valid:  true,
		},
	}

	email := fake.EmailAddress()
	fName := fake.FirstName()
	lName := fake.LastName()
	result = tx.MustExec(usersInsertQuery, fName, lName, email, fake.SimplePassword(), email, 0, time.Now(), time.Now(), "provider", driverID, "Driver")
	userID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	user := User{
		ID: int(userID),
		FirstName: sql.NullString{
			String: fName,
			Valid:  true,
		},
		LastName: sql.NullString{
			String: lName,
			Valid:  true,
		},
		EntityID: sql.NullInt64{
			Int64: int64(driver.ID),
			Valid: true,
		},
		EntityType: sql.NullString{
			String: "Driver",
			Valid:  true,
		},
	}

	driver.User = user
	return &driver, nil
}

func TestGetDriverByID(t *testing.T) {
	tx := createTx(t)
	defer tx.Rollback()

	driver, err := createDriver(tx, "on_duty")
	tst.FailNowOnErr(t, err)

	driver1, err := GetDriverByID(tx, driver.ID)
	tst.FailNowOnErr(t, err)

	assert.Equal(t, driver.ID, driver1.ID)
	assert.Equal(t, driver.Status, driver1.Status)
	assert.Equal(t, driver.User.ID, driver1.User.ID)
	assert.Equal(t, driver.User.FirstName, driver1.User.FirstName)
	assert.Equal(t, driver.User.MiddleName, driver1.User.MiddleName)
	assert.Equal(t, driver.User.LastName, driver1.User.LastName)
	assert.Equal(t, driver.User.EntityID, driver1.User.EntityID)
	assert.Equal(t, driver.User.EntityType, driver1.User.EntityType)
}

func TestGetDriverByTripID(t *testing.T) {
	tx := createTx(t)
	defer tx.Rollback()

	driver, err := createDriver(tx, "on_duty")
	tst.FailNowOnErr(t, err)

	trip, err := createTrip(&Trip{
		DriverID: driver.ID,
	}, tx)
	tst.FailNowOnErr(t, err)

	driver1, err := GetDriverByTripID(tx, trip.ID)
	tst.FailNowOnErr(t, err)

	assert.Equal(t, driver.ID, driver1.ID)
	assert.Equal(t, driver.Status, driver1.Status)
	assert.Equal(t, driver.User.ID, driver1.User.ID)
	assert.Equal(t, driver.User.FirstName, driver1.User.FirstName)
	assert.Equal(t, driver.User.MiddleName, driver1.User.MiddleName)
	assert.Equal(t, driver.User.LastName, driver1.User.LastName)
	assert.Equal(t, driver.User.EntityID, driver1.User.EntityID)
	assert.Equal(t, driver.User.EntityType, driver1.User.EntityType)
}

func TestGetDriverLocation(t *testing.T) {
	tx := createTx(t)
	defer tx.Rollback()

	driver, err := createDriver(tx, "on_duty")
	tst.FailNowOnErr(t, err)

	location := Location{
		utils.Location{
			Lat: 13.01,
			Lng: 79.01,
		},
	}
	tx.Exec("update users set current_location = ? where id = ?", location, driver.User.ID)

	driverLocation, err := DriverLocation(tx, driver.User.ID)
	tst.FailNowOnErr(t, err)

	assert.Equal(t, location, *driverLocation)
}

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"

	tst "github.com/MOOVE-Network/location_service/testutils"
	"github.com/jmoiron/sqlx"
)

func TestGetBufferDurationForDelayTripNotification(t *testing.T) {
	tx := createTx(t)
	defer tx.Rollback()
	_, err := createConfigurator(tx, "buffer_duration_for_delayed_trip_notification", "20")
	tst.FailNowOnErr(t, err)

	value, err := GetBufferDurationForDelayTripNotification(tx)
	tst.FailNowOnErr(t, err)

	assert.Equal(t, 20, value)
}

func createConfigurator(tx *sqlx.Tx, name string, value string) (*Configuration, error) {
	var insertConfigurationStmnt = `insert into configurators(request_type, value) values (?, ?)`
	res := tx.MustExec(insertConfigurationStmnt, name, value)
	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &Configuration{
		ID:          lastID,
		RequestType: name,
		Value:       value,
	}, nil
}

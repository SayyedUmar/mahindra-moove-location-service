package db

import (
	"strconv"

	"github.com/jmoiron/sqlx"
)

type Configuration struct {
	ID          int64  `db:"id"`
	RequestType string `db:"request_type"`
	Value       string `db:"value"`
}

var findByRequestTypeQuery = `select id, request_type, value 
	from configurators where request_type=?`

//GetBufferDurationForDelayTripNotification returns value of buffer_duration_for_delayed_trip_notification in minutes.
func GetBufferDurationForDelayTripNotification(db sqlx.Queryer) (int, error) {
	configuration, err := getConfiguration(db, "buffer_duration_for_delayed_trip_notification")
	if err != nil {
		return 0, err
	}
	timeInMinutes, err := strconv.Atoi(configuration.Value)
	if err != nil {
		return 0, err
	}
	return timeInMinutes, nil
}

//GetMinDistanceToCalculateStartTripEta returns value of min_distance_to_calc_start_trip_eta in meters.
func GetMinDistanceToCalculateStartTripEta(db sqlx.Queryer) (int, error) {
	configuration, err := getConfiguration(db, "min_distance_to_calc_start_trip_eta")
	if err != nil {
		return 0, err
	}
	timeInMinutes, err := strconv.Atoi(configuration.Value)
	if err != nil {
		return 0, err
	}
	return timeInMinutes, nil
}

//GetMaxTimeToCalculateStartTripEta returns value of max_time_to_calc_start_trip_eta in meters.
func GetMaxTimeToCalculateStartTripEta(db sqlx.Queryer) (int, error) {
	configuration, err := getConfiguration(db, "max_time_to_calc_start_trip_eta")
	if err != nil {
		return 0, err
	}

	timeInMinutes, err := strconv.Atoi(configuration.Value)
	if err != nil {
		return 0, err
	}
	return timeInMinutes, nil
}

func getConfiguration(db sqlx.Queryer, requestType string) (*Configuration, error) {
	row := db.QueryRowx(findByRequestTypeQuery, requestType)
	var configuration Configuration
	err := row.StructScan(&configuration)
	if err != nil {
		return nil, err
	}
	return &configuration, nil
}

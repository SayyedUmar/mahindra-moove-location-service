package db

import (
	"database/sql/driver"
	"errors"

	"github.com/MOOVE-Network/location_service/utils"
)

// Value function implements the Valuer interface for Location
func (loc Location) Value() (driver.Value, error) {
	yml, err := loc.ToYaml()
	if err != nil {
		return driver.Value(""), err
	}
	return yml, nil
}

// Scan function implements the Scanner interface for Location
func (loc *Location) Scan(src interface{}) error {
	bytearr, ok := src.([]byte)
	if !ok {
		return errors.New("invalid format detected")
	}
	l, err := utils.LocationFromYaml(string(bytearr))
	if err != nil {
		return err
	}
	loc.Location = *l
	return nil
}

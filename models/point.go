package models

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func (pt Point) Value() (driver.Value, error) {
	return fmt.Sprintf("(%f, %f)", pt.Lng, pt.Lat), nil
}
func (pt *Point) Scan(value interface{}) error {
	if value == nil {
		// its okay to have a nil represented as 0, 0
		pt.Lat = 0
		pt.Lng = 0
		return nil
	}
	if sv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := sv.(string); ok {
			latLngStr := strings.Split(strings.Trim(v, " ()"), ",")
			if len(latLngStr) != 2 {
				return fmt.Errorf("Unable to scan %s to point", v)
			}
			lat, latErr := strconv.ParseFloat(strings.Trim(latLngStr[1], " "), 64)
			lng, lngErr := strconv.ParseFloat(strings.Trim(latLngStr[0], " "), 64)
			if latErr != nil || lngErr != nil {
				return fmt.Errorf("Unable to scan %s to point", v)
			}
			pt.Lat = lat
			pt.Lng = lng
			return nil
		}
	} else {
		str, err := driver.String.ConvertValue(value)
		fmt.Println("unable to convert to string ", str, err)
	}
	return fmt.Errorf("Unable to scan to Point")
}

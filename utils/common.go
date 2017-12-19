package utils

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

// Location encodes a Latitude and Longitude
type Location struct {
	Lat float64 `yaml:":lat"`
	Lng float64 `yaml:":lng"`
}

// ToString returns a string representation of Location
func (l *Location) ToString() string {
	return fmt.Sprintf("%f,%f", l.Lat, l.Lng)
}

// ToYaml encodes location in rails friendly yaml hash
func (l *Location) ToYaml() (string, error) {
	return ToYamlLocation(l.Lat, l.Lng)
}

// LocationFromYaml returns a Location struct parsing rails friendly yaml
func LocationFromYaml(yml string) (*Location, error) {
	var loc Location
	err := yaml.Unmarshal([]byte(yml), &loc)
	if err != nil {
		return nil, err
	}
	return &loc, nil
}

// ToYamlLocation encodes lat and lng in a rails friendly Hash in YAML format
func ToYamlLocation(lat float64, lng float64) (string, error) {
	// content, err := yaml.Marshal(map[string]float64{":lat": lat, ":lng": lng})
	content, err := yaml.Marshal(Location{Lat: lat, Lng: lng})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("---\n%s", string(content)), nil
}

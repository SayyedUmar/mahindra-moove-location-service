package utils

import (
	"fmt"
	yaml "gopkg.in/yaml.v2"
)

func ToLocation(lat float64, lng float64) (string, error) {
	content, err := yaml.Marshal(map[string]float64{":lat": lat, ":lng": lng})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("---\n%s", string(content)), nil
}

package services

import (
	"time"

	"github.com/MOOVE-Network/location_service/utils"
)

// RoadsService Structure to access Match API
type RoadsService interface {
	Match(points []utils.Location, timestamps []time.Time) (*MatchResponse, error)
}

// MatchResponse Structure to receive response of  Match API
type MatchResponse struct {
	Code      string  `json:"code"`
	Matchings []Route `json:"matchings"`
	Distance  float64 `json:"distance"`
}

func (mr *MatchResponse) calculateTotalMileage() {
	mr.Distance = 0
	for _, r := range mr.Matchings {
		mr.Distance += r.Distance
	}
}

// Route Structure to define a route
type Route struct {
	Distance   float64 `json:"distance"`
	Duration   float64 `json:"duration"`
	Geometry   string  `json:"geometry"`
	Confidence float64 `json:"confidence"`
}

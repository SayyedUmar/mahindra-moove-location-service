package services

import (
	"time"

	"math"

	"github.com/MOOVE-Network/location_service/utils"
	geo "github.com/paulmach/go.geo"
	log "github.com/sirupsen/logrus"
)

// PassThroughRoadsService Structure to do a passthrough a Road Service
type PassThroughRoadsService struct {
}

// Match Function to do a passthrough a Road Service
func (pts *PassThroughRoadsService) Match(points []utils.Location, timestamps []time.Time) (*MatchResponse, error) {
	ps := geo.NewPointSet()
	for _, loc := range points {
		ps.Push(geo.NewPointFromLatLng(loc.Lat, loc.Lng))
	}
	path := &geo.Path{*ps}
	cleanPs := geo.NewPointSet()
	for i, pt := range path.Points() {
		log.Infof("Lat : %f, Lng: %f, Direction: %f", pt.Lat(), pt.Lng(), path.DirectionAt(i))
		if math.Abs(path.DirectionAt(i)) < 0.5 {
			cleanPs.Push(&pt)
		}
	}
	simplifiedPath := &geo.Path{*cleanPs}
	totalMileage := simplifiedPath.GeoDistance(true)
	geometry := simplifiedPath.Encode()
	return &MatchResponse{
		Code:      "Ok",
		Matchings: []Route{Route{Geometry: geometry}},
		Distance:  totalMileage,
	}, nil
}

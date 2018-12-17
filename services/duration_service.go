package services

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/MOOVE-Network/location_service/db"
	"github.com/MOOVE-Network/location_service/utils"
	log "github.com/sirupsen/logrus"
	"googlemaps.github.io/maps"
)

var durationService DurationService

// InitDurationService initializes the global
func InitDurationService(apiKey string) {
	gds, err := MakeGoogleDurationService(apiKey)
	if err != nil {
		panic(err)
	}
	durationService = gds
}

// SetDurationService sets the current duration service
func SetDurationService(ds DurationService) {
	durationService = ds
}

// GetDurationService returns the currently wrapped duration service
// This can be nil
func GetDurationService() DurationService {
	return durationService
}

// DurationMetrics Define a structure for duration metrics
type DurationMetrics struct {
	DepartureTime    time.Time
	ArrivalTime      time.Time
	Duration         time.Duration
	DistanceInMeters int
	StartLocation    db.Location
	EndLocation      db.Location
}

// DurationService is an interface used by GetDuration
type DurationService interface {
	GetDuration(start db.Location, end db.Location, startAt time.Time) (DurationMetrics, error)
}

// GoogleDurationService is an implementation of DurationService using Google Maps API
type GoogleDurationService struct {
	client *maps.Client
}

// GetDuration using google maps Directions API for Driving mode and avoiding tolls
func (g *GoogleDurationService) GetDuration(start db.Location, end db.Location, startAt time.Time) (DurationMetrics, error) {
	startAtSecs := strconv.FormatInt(startAt.Unix(), 10)
	dirRequest := &maps.DirectionsRequest{
		Origin:        start.ToString(),
		Destination:   end.ToString(),
		DepartureTime: startAtSecs,
		Mode:          maps.TravelModeDriving,
		TrafficModel:  maps.TrafficModelBestGuess,
		Avoid:         []maps.Avoid{maps.AvoidTolls},
	}
	log.Debug("Requesting directions with parameters")
	log.Debug(dirRequest)
	routes, _, err := g.client.Directions(context.Background(), dirRequest)
	if err != nil {
		return DurationMetrics{}, err
	}
	leg := routes[0].Legs[0]
	dm := DurationMetrics{
		DepartureTime:    leg.DepartureTime,
		ArrivalTime:      leg.ArrivalTime,
		Duration:         leg.DurationInTraffic,
		DistanceInMeters: leg.Distance.Meters,
		StartLocation:    db.Location{utils.Location{Lat: leg.StartLocation.Lat, Lng: leg.StartLocation.Lng}},
		EndLocation:      db.Location{utils.Location{Lat: leg.EndLocation.Lat, Lng: leg.EndLocation.Lng}},
	}
	return dm, nil
}

// MakeGoogleDurationService creates a google duration service given an API Key
func MakeGoogleDurationService(apiKey string) (*GoogleDurationService, error) {
	if apiKey == "" {
		return nil, errors.New("API Key is empty")
	}
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &GoogleDurationService{client: client}, nil
}

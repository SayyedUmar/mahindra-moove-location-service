package services

import (
	"context"
	"strconv"
	"time"

	"github.com/MOOVE-Network/location_service/db"
	"googlemaps.github.io/maps"
)

// DurationService is an interface used by GetDuration
type DurationService interface {
	GetDuration(start db.Location, end db.Location, startAt ...time.Time) (time.Duration, error)
}

// GoogleDurationService is an implementation of DurationService using Google Maps API
type GoogleDurationService struct {
	client *maps.Client
}

// GetDuration using google maps Directions API for Driving mode and avoiding tolls
func (g *GoogleDurationService) GetDuration(start db.Location, end db.Location, startAt time.Time) (time.Duration, error) {
	startAtSecs := strconv.FormatInt(startAt.Unix(), 10)
	dirRequest := &maps.DirectionsRequest{
		Origin:        start.ToString(),
		Destination:   end.ToString(),
		DepartureTime: startAtSecs,
		Mode:          maps.TravelModeDriving,
		TrafficModel:  maps.TrafficModelBestGuess,
		Avoid:         []maps.Avoid{maps.AvoidTolls},
	}
	routes, _, err := g.client.Directions(context.Background(), dirRequest)
	if err != nil {
		return 1 * time.Second, err
	}
	return routes[0].Legs[0].DurationInTraffic, nil
}

// MakeGoogleDurationService creates a google duration service given an API Key
func MakeGoogleDurationService(apiKey string) (*GoogleDurationService, error) {
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &GoogleDurationService{client: client}, nil
}

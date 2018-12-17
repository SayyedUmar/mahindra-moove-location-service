package services

import (
	"context"
	"errors"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"net/http"

	"github.com/MOOVE-Network/location_service/utils"
	geo "github.com/paulmach/go.geo"
	log "github.com/sirupsen/logrus"
	"googlemaps.github.io/maps"
)

var googleRoadsService RoadsService

// GoogleRoadsService Structure add google roads service with the same contract as the OSRM service
type GoogleRoadsService struct {
	client *maps.Client
}

// InitGoogleRoadsService function to initiate google roads service with the same contract as the OSRM service
func InitGoogleRoadsService(apiKey string) {
	grs, err := MakeGoogleRoadsService(apiKey)
	if err != nil {
		panic(err)
	}
	googleRoadsService = grs
}

// GetGoogleRoadsService Function to get google roads service with the same contract as the OSRM service
func GetGoogleRoadsService() RoadsService {
	return googleRoadsService
}

// MakeGoogleRoadsService Function to create a Google Road Service
func MakeGoogleRoadsService(apiKey string) (*GoogleRoadsService, error) {
	if apiKey == "" {
		return nil, errors.New("API Key is empty")
	}
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &GoogleRoadsService{client: client}, nil
}

func (rs *GoogleRoadsService) snapToRoads(points []utils.Location) (*Route, error) {
	var latLngs []maps.LatLng
	ps := geo.NewPointSet()
	for _, loc := range points {
		ps.Push(geo.NewPointFromLatLng(loc.Lat, loc.Lng))
	}
	path := &geo.Path{*ps}
	for i, pt := range path.Points() {
		log.Infof("Lat : %f, Lng: %f, Direction: %f", pt.Lat(), pt.Lng(), path.DirectionAt(i))
		if math.Abs(path.DirectionAt(i)) < 1.7 {
			latLngs = append(latLngs, maps.LatLng{Lat: pt.Lat(), Lng: pt.Lng()})
		}
	}

	strRequest := maps.SnapToRoadRequest{
		Path:        latLngs,
		Interpolate: true,
	}
	// Only do this for debugging : makeRequest(strRequest)

	strResponse, err := rs.client.SnapToRoad(context.Background(), &strRequest)
	if err != nil {
		return nil, err
	}
	pts := geo.NewPointSet()
	for _, pt := range strResponse.SnappedPoints {
		pts = pts.Push(geo.NewPointFromLatLng(pt.Location.Lat, pt.Location.Lng))
	}
	// log.Infof("found points %v ", pts)
	path = &geo.Path{PointSet: *pts}

	encodedString := path.Encode()
	route := Route{
		Geometry:   encodedString,
		Distance:   path.GeoDistance(true),
		Confidence: 0.8,
	}
	return &route, nil
}

// Match Function match goole road service with OSRM Srvice
func (rs *GoogleRoadsService) Match(points []utils.Location, _ []time.Time) (*MatchResponse, error) {
	var pointSets [][]utils.Location
	if len(points) > 100 {
		for i := 0; i <= len(points)/100; i++ {
			var pointSet []utils.Location
			for j := 0; j < 100; j++ {
				idx := i*100 + j
				if idx < len(points) {
					pointSet = append(pointSet, points[idx])
				} else {
					break
				}
			}
			pointSets = append(pointSets, pointSet)
		}
	} else {
		pointSets = append(pointSets, points)
	}
	totalMileage := 0.0
	var routes []Route
	for _, pointSet := range pointSets {
		route, err := rs.snapToRoads(pointSet)
		if err != nil {
			log.Errorf("error while calculation route for point set %s", err.Error())
			return &MatchResponse{Code: "Invalid"}, err
		}
		routes = append(routes, *route)
		totalMileage += route.Distance
	}
	log.Info(routes)
	return &MatchResponse{
		Code:      "Ok",
		Matchings: routes,
		Distance:  totalMileage,
	}, nil
}

// LocationToLatLngs Function that converts location to lat,long
func LocationToLatLngs(loc utils.Location) maps.LatLng {
	return maps.LatLng{Lat: loc.Lat, Lng: loc.Lng}
}

// httpRequest Function to use request SnapToRoad
func httpRequest(strRequest maps.SnapToRoadRequest) (*http.Request, error) {
	baseUrl := "https://roads.googleapis.com/v1/snapToRoads"
	req, err := http.NewRequest("GET", baseUrl, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	var locs []string
	for _, loc := range strRequest.Path {
		locs = append(locs, loc.String())
	}
	pathString := strings.Join(locs, "|")
	query.Set("path", pathString)
	query.Set("interpolate", strconv.FormatBool(strRequest.Interpolate))
	query.Set("key", os.Getenv("LOCATION_MAPS_API_KEY"))
	req.URL.RawQuery = query.Encode()
	return req, nil
}

func makeRequest(strRequest maps.SnapToRoadRequest) {
	req, err := httpRequest(strRequest)
	if err != nil {
		log.Errorf("Error creating request : %s", err.Error())
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("Error peformign reqeust %s ", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("Error reading body : %s", err.Error())
			return
		}
		log.Error("Success Body is %s ", string(body))
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("Error reading body : %s", err.Error())
			return
		}
		log.Error("Body is %s", string(body))
		log.Error("Status is %s", resp.Status)
	}

}

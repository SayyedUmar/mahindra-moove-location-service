package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/MOOVE-Network/location_service/utils"
)

type OSRMClient struct {
	baseURL string
}

func NewOSRMClient(baseUrl string) *OSRMClient {
	return &OSRMClient{baseURL: baseUrl}
}

func (c *OSRMClient) Match(locations []utils.Location, timestamps []time.Time) (*MatchResponse, error) {
	var locs []string
	for _, loc := range locations {
		locs = append(locs, fmt.Sprintf("%f,%f", loc.Lng, loc.Lat))
	}
	locsString := strings.Join(locs, ";")
	osrmUrl, err := url.Parse(fmt.Sprintf("%s/match/v1/car/%s", c.baseURL, locsString))
	if err != nil {
		return nil, err
	}
	var tsInts []int64
	for _, t := range timestamps {
		tsInts = append(tsInts, t.Unix())
	}
	var tsString []string
	for _, ts := range tsInts {
		tsString = append(tsString, strconv.FormatInt(ts, 10))
	}
	timestampsString := strings.Join(tsString, ";")

	query := osrmUrl.Query()
	query.Add("overview", "full")
	query.Add("timestamps", timestampsString)
	osrmUrl.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", osrmUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	var matchResp MatchResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&matchResp)
	if err != nil {
		return nil, err
	}
	matchResp.calculateTotalMileage()
	return &matchResp, nil
}

type MatchResponse struct {
	Code      string  `json:"code"`
	Matchings []Route `json:"matchings"`
	Distance  float64 `json:"distance"`
}

func (mr *MatchResponse) calculateTotalMileage() {
	for _, r := range mr.Matchings {
		mr.Distance += r.Distance
	}
}

type Route struct {
	Distance   float64 `json:"distance"`
	Duration   float64 `json:"duration"`
	Geometry   string  `json:"geometry"`
	Confidence float64 `json:"confidence"`
}

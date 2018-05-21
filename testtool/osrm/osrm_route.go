package osrm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"

	geo "github.com/paulmach/go.geo"
)

const ApiVersion = "v1"
const Profile = "driving"
const RouteService = "route"
const TripService = "trip"

type OSRMClient struct {
	baseUrl string
}
type RouteOptions struct {
	Locations    []geo.Pointer
	Overview     string
	Steps        bool
	Alternatives bool
}

func (opts RouteOptions) buildRequest(baseUrl, service string) (*http.Request, error) {
	if service != RouteService && service != TripService {
		return nil, fmt.Errorf("Service should be one of %s, %s", RouteService, TripService)
	}
	if len(opts.Locations) < 2 {
		return nil, fmt.Errorf("you need to have atleast 2 locations")
	}
	var locs []string
	for _, loc := range opts.Locations {
		locs = append(locs, fmt.Sprintf("%f,%f", loc.Point().Lng(), loc.Point().Lat()))
	}
	locations := strings.Join(locs, ";")
	reqUrl := fmt.Sprintf("%s/%s/%s/%s/%s", baseUrl, service, ApiVersion, Profile, locations)
	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("overview", opts.Overview)
	q.Add("steps", strconv.FormatBool(opts.Steps))
	if service == RouteService {
		q.Add("alternatives", strconv.FormatBool(opts.Alternatives))
	}
	if service == TripService {
		q.Add("source", "first")
		q.Add("destination", "last")
		q.Add("roundtrip", "false")
	}
	req.URL.RawQuery = q.Encode()
	return req, nil
}

func NewOSRMClient(baseUrl string) *OSRMClient {
	return &OSRMClient{baseUrl: baseUrl}
}

type RouteResponse struct {
	Routes []struct {
		Geometry string  `json:"geometry"`
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
	} `json:"routes"`
}

type TripResponse struct {
	Waypoints []Waypoint `json:"waypoints"`
	Trips     []struct {
		Geometry string  `json:"geometry"`
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
	}
}
type ByWaypointIndex []Waypoint

func (wp ByWaypointIndex) Len() int           { return len(wp) }
func (wp ByWaypointIndex) Swap(i, j int)      { wp[i], wp[j] = wp[j], wp[i] }
func (wp ByWaypointIndex) Less(i, j int) bool { return wp[i].TripsIndex < wp[j].TripsIndex }

func (tr *TripResponse) sortWaypoints() {
	sortableWPs := ByWaypointIndex(tr.Waypoints)
	sort.Sort(ByWaypointIndex(tr.Waypoints))
	tr.Waypoints = sortableWPs
}

type Waypoint struct {
	WaypointIndex int       `json:"waypoint_index"`
	Location      []float64 `json:"location"`
	TripsIndex    int       `json:"trips_index"`
}

func (wp Waypoint) Point() *geo.Point {
	return geo.NewPointFromLatLng(wp.Location[1], wp.Location[0])
}

func (client *OSRMClient) GetTrip(opts RouteOptions) (*TripResponse, error) {
	req, err := opts.buildRequest(client.baseUrl, TripService)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var tripResp TripResponse
	err = json.Unmarshal(body, &tripResp)
	if err != nil {
		return nil, err
	}
	tripResp.sortWaypoints()
	return &tripResp, nil
}

func (client *OSRMClient) GetRoute(opts RouteOptions) (*RouteResponse, error) {
	req, err := opts.buildRequest(client.baseUrl, RouteService)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var routeResp RouteResponse
	json.Unmarshal(body, &routeResp)
	if err != nil {
		return nil, err
	}
	return &routeResp, nil
}

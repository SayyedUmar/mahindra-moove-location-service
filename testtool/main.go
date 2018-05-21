package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/MOOVE-Network/location_service/socketstore"

	"github.com/MOOVE-Network/location_service/db"
	"github.com/MOOVE-Network/location_service/testtool/osrm"
	"github.com/MOOVE-Network/location_service/testtool/socketclient"
	"github.com/paulmach/go.geo"
)

func init() {
	flag.Parse()
}

type Location struct {
	Lat float64
	Lng float64
}

func (l Location) Point() *geo.Point {
	return geo.NewPointFromLatLng(l.Lat, l.Lng)
}

func main() {
	if flag.NArg() < 1 {
		fmt.Println("You should pass in a trip id")
		os.Exit(1)
	}

	tripID, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		fmt.Println("Trip id should be a number", flag.Arg(0))
		panic(err)
	}
	client := osrm.NewOSRMClient(getOSRMURL())
	conn := db.InitSQLConnection()
	db.SetActiveDB(conn)
	defer conn.Close()

	t, err := db.GetTripByID(conn, tripID)
	if err != nil {
		panic(err)
	}
	var locations []geo.Pointer
	for _, tr := range t.TripRoutes {
		locations = append(locations, Location{
			tr.ScheduledStartLocation.Location.Lat,
			tr.ScheduledStartLocation.Location.Lng,
		})
	}
	locations = append(locations, Location{
		t.TripRoutes[len(t.TripRoutes)-1].ScheduledEndLocation.Lat,
		t.TripRoutes[len(t.TripRoutes)-1].ScheduledEndLocation.Lng,
	})
	ro := osrm.RouteOptions{
		Locations:    locations,
		Overview:     "full",
		Steps:        false,
		Alternatives: false,
	}
	resp, err := client.GetRoute(ro)
	if err != nil {
		panic(err)
	}
	geom := resp.Routes[0].Geometry
	finish := make(chan bool)
	ptChan := emitPoints(geom, finish)
	wsConn := socketclient.SetupWSConnection(t.DriverUserID)
outer:
	for {
		select {
		case pt := <-ptChan:
			lu := socketstore.LocationUpdate{
				EventType:   "LOCATION",
				Lat:         pt.Lat(),
				Lng:         pt.Lng(),
				TripID:      t.ID,
				UserID:      strconv.Itoa(t.DriverUserID),
				Speed:       30,
				Distance:    0,
				UnixSeconds: time.Now().Unix() * 1000,
			}
			fmt.Println("writing json")
			wsConn.WriteJSON(lu)
		case <-finish:
			break outer
		}
	}
	wsConn.Close()
	fmt.Println("Done ... !!!")

}

func emitPoints(geom string, finish chan bool) chan geo.Point {
	path := geo.NewPathFromEncoding(geom)
	outChan := make(chan geo.Point)
	go func() {
		for _, pt := range path.Points() {
			outChan <- pt
			time.Sleep(200 * time.Millisecond)
		}
		finish <- true
	}()

	return outChan
}

func getOSRMURL() string {
	osrmURL := os.Getenv("OSRM_URL")
	if osrmURL == "" {
		osrmURL = "http://ec2-13-127-26-106.ap-south-1.compute.amazonaws.com:5000"
	}
	return osrmURL
}

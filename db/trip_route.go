package db

// OnBoardStatuses contains of list of statuses
// when employees are either on board or have been cancelled
var OnBoardStatuses = map[string]bool{
	"canceled":       true,
	"on_board":       true,
	"missed":         true,
	"driver_arrived": true,
}

// TripRoute represents the database structure of TripRoute
type TripRoute struct {
	ID                     int      `db:"id"`
	TripID                 int      `db:"trip_id"`
	Status                 string   `db:"status"`
	ScheduledRouteOrder    int      `db:"scheduled_route_order"`
	ScheduledStartLocation Location `db:"scheduled_start_location"`
	ScheduledEndLocation   Location `db:"scheduled_end_location"`
	Trip                   *Trip
}

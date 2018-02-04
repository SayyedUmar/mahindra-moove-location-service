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
	EmployeeUserID         int      `db:"employee_user_id"`
	Trip                   Trip
}

// IsOnBoard is considered on board if he is on board or driver has arrived
func (tr *TripRoute) IsOnBoard() bool {
	return tr.Status == "on_board" || tr.Status == "driver_arrived"
}

package migrations

const AddDriverShouldStartTripColumnsInTrip = `
ALTER TABLE trips
ADD COLUMN driver_should_start_trip_time datetime default NULL,
ADD COLUMN driver_should_start_trip_time text default NULL;
`

func init() {
	initGM()
	GlobalMigrations.Register(NewMigration("20180413162723", AddDriverShouldStartTripColumnsInTrip))
}

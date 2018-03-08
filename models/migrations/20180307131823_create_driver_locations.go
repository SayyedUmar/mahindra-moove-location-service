package migrations

const CreateDriverLocations = `
	create extension if not exists "uuid-ossp";
	create table driver_locations(
		id uuid primary key default uuid_generate_v1mc(),
		recorded_at timestamp not null,
		trip_id int,
		user_id text,
		location point,
		distance int,
		speed double precision,
		accuracy double precision,
		created_at timestamp
	);
	create index driver_locations_trip_id_idx on driver_locations (trip_id);
	create index driver_locations_user_id_idx on driver_locations (user_id);
	create index driver_locations_recorded_at_idx on driver_locations (recorded_at);
`

func init() {
	initGM()
	GlobalMigrations.Register(NewMigration("20180307131823", CreateDriverLocations))
}

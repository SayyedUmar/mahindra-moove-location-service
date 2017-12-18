package services

import "github.com/MOOVE-Network/location_service/db"

func GetETAForTrip(trip *db.Trip) error {
	if err := trip.LoadTripRoutes(db.CurrentDB(), false); err != nil {
		return err
	}
	return nil
}

func GetETAForTripID(tripID int) error {
	trip, err := db.GetTripByID(db.CurrentDB(), tripID)
	if err != nil {
		return err
	}
	return GetETAForTrip(trip)
}

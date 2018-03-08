package models

import (
	"testing"
	"time"

	null "gopkg.in/guregu/null.v3"
)

func TestDriverLocationPrepareInsertStmt_ShouldPrepareStatement(t *testing.T) {
	_, err := DriverLocationPrepareInsertStmt(CurrentDB())
	if err != nil {
		t.Logf("Could not prepare statement for insert %s", err)
		t.FailNow()
	}
}

func TestDriverLocation_Insert_ShouldInsertDriverLocation(t *testing.T) {
	stmt, err := DriverLocationPrepareInsertStmt(CurrentDB())
	if err != nil {
		t.Logf("Could not prepare statement for insert %s", err)
		t.FailNow()
	}
	dl := DriverLocation{
		RecordedAt: time.Now(),
		TripID:     null.IntFrom(23),
		UserID:     null.StringFrom("53"),
		Location:   Point{Lat: 40, Lng: 30},
		Distance:   300,
		Speed:      400,
		Accuracy:   0.8,
		CreatedAt:  time.Now(),
	}
	err = dl.Insert(stmt)
	if err != nil {
		t.Logf("Could not insert driver location - %s", err)
		t.FailNow()
	}
}

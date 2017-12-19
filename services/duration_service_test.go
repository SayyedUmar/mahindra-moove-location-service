package services

import (
	"os"
	"testing"

	"time"

	"github.com/MOOVE-Network/location_service/db"
	tst "github.com/MOOVE-Network/location_service/testutils"
	"github.com/MOOVE-Network/location_service/utils"
	"github.com/stretchr/testify/assert"
)

const TarkaLabsLat float64 = 13.0363178
const TarkaLabsLng float64 = 80.2142206

const HomeOneLat float64 = 12.8947997
const HomeOneLng float64 = 80.2010107

func TestGoogleDurationService_GetDuration_Actual(t *testing.T) {
	mapsAPIKey := os.Getenv("LOCATION_MAPS_API_KEY")
	gds, err := MakeGoogleDurationService(mapsAPIKey)
	tst.FailNowOnErr(t, err)
	startLocation := db.Location{utils.Location{Lat: TarkaLabsLat, Lng: TarkaLabsLng}}
	endLocation := db.Location{utils.Location{Lat: HomeOneLat, Lng: HomeOneLng}}
	duration, err := gds.GetDuration(startLocation, endLocation, time.Now())
	tst.FailNowOnErr(t, err)
	assert.True(t, duration.Seconds() > 1)
}

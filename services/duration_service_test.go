package services

import (
	"os"
	"strings"
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

// This tests only exists for sanity and provides
// no valuable assertions
func TestGoogleDurationService_GetDuration_Actual(t *testing.T) {
	mapsAPIKey := os.Getenv("LOCATION_MAPS_API_KEY")
	if mapsAPIKey == "" {
		t.Skip("Skipping duration test since environment variable is unset")
		return
	}
	gds, err := MakeGoogleDurationService(mapsAPIKey)
	tst.FailNowOnErr(t, err)
	startLocation := db.Location{utils.Location{Lat: TarkaLabsLat, Lng: TarkaLabsLng}}
	endLocation := db.Location{utils.Location{Lat: HomeOneLat, Lng: HomeOneLng}}
	dm, err := gds.GetDuration(startLocation, endLocation, time.Now())
	if err != nil {
		assert.True(t, strings.HasPrefix(err.Error(), "maps: OVER_QUERY_LIMIT - You have exceeded your daily request quota for this API. We recommend registering for a key at the Google Developers Console"))
	} else {
		tst.FailNowOnErr(t, err)
		assert.True(t, dm.Duration.Seconds() > 1)
	}
}

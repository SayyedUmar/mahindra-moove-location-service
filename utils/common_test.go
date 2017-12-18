package utils

import (
	"testing"

	"math"

	"github.com/stretchr/testify/assert"
)

const TarkaLabsLat float64 = 13.0363178
const TarkaLabsLng float64 = 80.2142206

func TestLocation_ToYaml(t *testing.T) {
	tarkaLabsLoc := Location{Lat: TarkaLabsLat, Lng: TarkaLabsLng}
	yml, err := tarkaLabsLoc.ToYaml()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	expectedYml := "---\n:lat: 13.036318\n:lng: 80.21422\n"
	t.Log(yml)
	assert.Equal(t, expectedYml, yml)
}
func TestToYamlLocation(t *testing.T) {
	yml, err := ToYamlLocation(TarkaLabsLat, TarkaLabsLng)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	expectedYml := "---\n:lat: 13.036318\n:lng: 80.21422\n"
	t.Log(yml)
	assert.Equal(t, expectedYml, yml)
}

func TestLocationFromYaml(t *testing.T) {
	ymlString := "---\n:lat: 13.036318\n:lng: 80.21422\n"
	loc, err := LocationFromYaml(ymlString)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if math.Abs(loc.Lat-TarkaLabsLat) > 0.00001 {
		t.Logf("expected lat to be %f found %f", TarkaLabsLat, loc.Lat)
		t.Log(loc)
		t.FailNow()
	}
	if math.Abs(loc.Lng-TarkaLabsLng) > 0.00001 {
		t.Logf("expected lng to be %f found %f", TarkaLabsLng, loc.Lng)
		t.FailNow()
	}
}

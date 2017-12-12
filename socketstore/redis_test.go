package socketstore

import (
	"fmt"
	"os"
	"testing"
)

func TestPipelining(t *testing.T) {
	p := GetClient().Pipeline()
	p.RPush("values", 1)
	p.RPush("values", 2)
	p.RPush("values", 3)
	_, err := p.Exec()

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	vals, err := GetClient().LRange("values", 0, -1).Result()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(vals)
	if len(vals) != 3 {
		t.Errorf("expected %d found %d", 3, len(vals))
		t.FailNow()
	}
}

func TestMain(m *testing.M) {
	setupRedis()
	_, err := GetClient().FlushDB().Result()
	if err != nil {
		fmt.Println("unable to connect or flush redis")
		fmt.Println(err)
		os.Exit(-1)
	}
	os.Exit(m.Run())
}

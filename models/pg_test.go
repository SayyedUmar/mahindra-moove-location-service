package models

import (
	"os"
	"testing"
)

func handleExit() {
	if e := recover(); e != nil {
		if exit, ok := e.(exit); ok {
			os.Exit(exit.code)
		}
		panic(e) // not an Exit, bubble up
	}
}

type exit struct{ code int }

func TestMain(m *testing.M) {
	defer handleExit()
	if os.Getenv("LOCATION_PG_DATABASE_URL") == "" {
		os.Setenv("LOCATION_PG_DATABASE_URL", "postgres://moove:m00ve4wd@clustering-pg-test.ct3tozaserta.ap-southeast-1.rds.amazonaws.com/location_service?sslmode=disable")
	}
	db := InitSQLConnection()
	SetActiveDB(db)
	RunMigrations()
	defer db.Close()
	panic(exit{code: m.Run()})
}

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
	os.Setenv("LOCATION_PG_DATABASE_URL", "postgres://localhost/location_service_test?sslmode=disable")
	db := InitSQLConnection()
	SetActiveDB(db)
	RunMigrations()
	defer db.Close()
	panic(exit{code: m.Run()})
}

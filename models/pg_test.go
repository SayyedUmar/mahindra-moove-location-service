package models

import (
	"fmt"
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
		const (
			host     = "moove-pg-uat10.cjny84emnsh9.ap-south-1.rds.amazonaws.com"
			port     = 5432
			user     = "MOOVE_DEV"
			password = ""
			dbname   = "moove-pg-uat10"
		)
		dbURL := "MOOVE_DEV:NG$Pir7ySMJ9m&p9@moove-pg-uat10.cjny84emnsh9.ap-south-1.rds.amazonaws.com/location_service?sslmode=disable"
		//fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
		os.Setenv("LOCATION_PG_DATABASE_URL", dbURL)
	}
	db := InitSQLConnection()
	SetActiveDB(db)
	RunMigrations()
	defer db.Close()
	panic(exit{code: m.Run()})
}

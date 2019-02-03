package models

import (
	"os"

	"github.com/MOOVE-Network/location_service/models/migrations"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

var db *sqlx.DB

func InitSQLConnection() *sqlx.DB {
	dbUrl := os.Getenv("LOCATION_PG_DATABASE_URL")
	if dbUrl == "" {
		dbUrl = "postgres://clustering-pg-test.ct3tozaserta.ap-southeast-1.rds.amazonaws.com/location_service?sslmode=disable"
	}
	localDb, err := sqlx.Open("postgres", dbUrl)
	if err != nil {
		panic("something is wrong with open")
	}
	return localDb
}

func CurrentDB() *sqlx.DB {
	return db
}

func SetActiveDB(active *sqlx.DB) {
	db = active
}

func RunMigrations() {
	err := migrations.GlobalMigrations.Apply(CurrentDB())
	if err != nil {
		log.Fatal("Unable to run migrations")
		panic(err)
	}
}

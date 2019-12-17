package models

import (
	"fmt"

	"github.com/MOOVE-Network/location_service/models/migrations"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

var db *sqlx.DB

func InitSQLConnection() *sqlx.DB {
	// dbUrl := os.Getenv("LOCATION_PG_DATABASE_URL")
	// if dbUrl == "" {
	// 	dbUrl = "postgres://localhost/location_service_dev?sslmode=disable"
	// }
	const (
		host     = "moove-pg-uat10.cjny84emnsh9.ap-south-1.rds.amazonaws.com"
		port     = 5432
		user     = "MOOVE_DEV"
		password = ""
		dbname   = "moove-pg-uat10"
	)
	dbURL := "MOOVE_DEV:NG$Pir7ySMJ9m&p9@moove-pg-uat10.cjny84emnsh9.ap-south-1.rds.amazonaws.com/location_service?sslmode=disable"
		//fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
		panic("postgres Database connecting..")
	localDb, err := sqlx.Open("postgres", dbURL)
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

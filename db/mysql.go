package db

import (
	"os"
	// mysql driver import for database/sql
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var database *sqlx.DB

func InitSQLConnection() *sqlx.DB {
	dbUrl := os.Getenv("LOCATION_DATABASE_URL")
	if dbUrl == "" {
		dbUrl = "root:@/moove_development?parseTime=true"
	}
	localDb, err := sqlx.Open("mysql", dbUrl)
	if err != nil {
		panic("something is wrong with open")
	}
	return localDb
}

func CurrentDB() *sqlx.DB {
	return database
}

func SetActiveDB(active *sqlx.DB) {
	database = active
}

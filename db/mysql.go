package db

import (
	"os"
	// mysql driver import for database/sql
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var database *sqlx.DB

// InitSQLConnection initiate a db connection and returns db handler
func InitSQLConnection() *sqlx.DB {
	dbUrl := os.Getenv("LOCATION_DATABASE_URL")
	if dbUrl == "" {
		dbUrl = "MOOVE_DEV:m00ve4wd@tcp(grab.ct3tozaserta.ap-southeast-1.rds.amazonaws.com:3306)/grab?parseTime=true"
	}
	localDb, err := sqlx.Open("mysql", dbUrl)
	if err != nil {
		panic("something is wrong with open")
	}
	return localDb
}

// CurrentDB get the database handler of current db
func CurrentDB() *sqlx.DB {
	return database
}

// SetActiveDB activate a db with database handler
func SetActiveDB(active *sqlx.DB) {
	database = active
}

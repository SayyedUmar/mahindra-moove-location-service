package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"os"
)

var database *sqlx.DB

func InitSQLConnection() *sqlx.DB {
	dbUrl := os.Getenv("CLUSTERING_DATABASE_URL")
	if dbUrl == "" {
		dbUrl = "root:@/moove_development"
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

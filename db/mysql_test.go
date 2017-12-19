package db

import (
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
)

type exit struct{ code int }

func createTx(t *testing.T) *sqlx.Tx {
	tx, err := CurrentDB().Beginx()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	return tx
}

// exit code handler
func handleExit() {
	if e := recover(); e != nil {
		if exit, ok := e.(exit); ok {
			os.Exit(exit.code)
		}
		panic(e) // not an Exit, bubble up
	}
}

func TestMain(m *testing.M) {
	defer handleExit()
	os.Setenv("LOCATION_DATABASE_URL", "root:@/moove_test?multiStatements=true&autocommit=false&parseTime=true")
	db := InitSQLConnection()
	_, err := sqlx.LoadFile(db, "../data/mysql_structure.sql")
	if err != nil {
		panic(err)
	}

	SetActiveDB(db)
	defer db.Close()
	panic(exit{code: m.Run()})
}

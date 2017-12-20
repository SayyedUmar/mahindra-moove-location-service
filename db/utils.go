package db

import (
	"database/sql"

	"github.com/MOOVE-Network/location_service/utils"
	"github.com/jmoiron/sqlx"
)

type Location struct{ utils.Location }

type Preparex interface {
	sqlx.Preparer
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
}

type NamedExecer interface {
	Exec(interface{}) (sql.Result, error)
}

type RebindQueryer interface {
	sqlx.Queryer
	Rebind(string) string
}

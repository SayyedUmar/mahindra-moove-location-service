package db

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type Preparex interface {
	sqlx.Preparer
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
}

type NamedExecer interface {
	Exec(interface{}) (sql.Result, error)
}

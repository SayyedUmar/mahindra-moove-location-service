package db

import "database/sql"

type User struct {
	ID         int            `db:"id"`
	FirstName  sql.NullString `db:"f_name"`
	MiddleName sql.NullString `db:"m_name"`
	LastName   sql.NullString `db:"l_name"`
	EntityType sql.NullString `db:"entity_type"`
	EntityID   sql.NullInt64  `db:"entity_id"`
}

package db

type User struct {
	ID         int     `db:"id"`
	FirstName  string  `db:"f_name"`
	MiddleName *string `db:"m_name"`
	LastName   string  `db:"l_name"`
	EntityType string  `db:"entity_type"`
	EntityID   int     `db:"entity_id"`
}

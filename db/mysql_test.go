package db

import (
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
)

type exit struct{ code int }

type bareUser struct {
	ID                int
	Email             string
	EncryptedPassword string
	SignInCount       int
	Provider          string
	Uid               string
}
type bareEmployee struct {
	ID         int
	EmployeeID string
}

var users = []bareUser{
	bareUser{Email: "vagmi@tarkalabs.com", EncryptedPassword: "blah", SignInCount: 0, Provider: "email", Uid: "vagmi@tarkalabs.com"},
	bareUser{Email: "sudhakar@tarkalabs.com", EncryptedPassword: "blah", SignInCount: 0, Provider: "email", Uid: "sudhakar@tarkalabs.com"},
	bareUser{Email: "dhruva@tarkalabs.com", EncryptedPassword: "blah", SignInCount: 0, Provider: "email", Uid: "dhruva@tarkalabs.com"},
}
var employees = []bareEmployee{
	bareEmployee{EmployeeID: "vagmi@tarkalabs.com"},
	bareEmployee{EmployeeID: "sudhakar@tarkalabs.com"},
	bareEmployee{EmployeeID: "dhruva@tarkalabs.com"},
}

func getEmployeeID(db sqlx.Queryer, email string) int {
	var employeeID int
	row := db.QueryRowx(`select id from employees where employee_id=?`, email)
	err := row.Scan(&employeeID)
	if err != nil {
		panic(err)
	}
	return employeeID
}

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
func seedData(db *sqlx.DB) {
	for _, e := range employees {
		db.MustExec(`
			insert into employees(employee_id, created_at, updated_at)
			values(?,?,?)
			`, e.EmployeeID, time.Now(), time.Now())

	}
	for _, u := range users {
		var employeeID int
		row := db.QueryRow(`select id from employees where employee_id=?`, u.Email)
		err := row.Scan(&employeeID)
		if err != nil {
			panic(err)
		}
		db.MustExec(`
		insert into users(email, encrypted_password, sign_in_count, 
						  created_at, updated_at, provider, 
						  uid, entity_type, entity_id)
		values(?,?,?,?,?,?,?,?,?)
		`, u.Email, u.EncryptedPassword, u.SignInCount,
			time.Now(), time.Now(), u.Provider,
			u.Uid, "Employee", employeeID)
	}
}

func TestMain(m *testing.M) {
	defer handleExit()
	os.Setenv("LOCATION_DATABASE_URL", "MOOVE_DEV:NG$Pir7ySMJ9m&p9@tcp(vaayu-uat.cjny84emnsh9.ap-south-1.rds.amazonaws.com:3306)/moove_db_uat?parseTime=true")
	db := InitSQLConnection()
	_, err := sqlx.LoadFile(db, "../data/mysql_structure.sql")
	seedData(db)
	if err != nil {
		panic(err)
	}
	SetActiveDB(db)
	defer db.Close()
	panic(exit{code: m.Run()})
}

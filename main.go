package main

import (
	"bytes"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

var Db *sqlx.DB

type Token struct {
	Token  string `json:"token"`
	Expiry int64  `json:"expiry"`
}
type Identity struct {
	Uid          string
	Tokens       string
	ParsedTokens map[string]Token
}

func FetchIdentityByUID(db *sqlx.DB, uid string) *Identity {
	var i Identity
	row := db.QueryRowx("select uid, tokens from users where uid=?", uid)
	if row != nil {
		err := row.StructScan(&i)
		if err != nil {
			return nil
		}
		i.parseTokens()
		return &i
	}
	return &i
}

func (i *Identity) parseTokens() {
	buffer := bytes.NewBuffer([]byte(i.Tokens))
	decoder := json.NewDecoder(buffer)
	err := decoder.Decode(&i.ParsedTokens)
	if err != nil {
		panic(err)
	}
}

func (i *Identity) isValid(clientId, token string) bool {
	return i.isValidToken(clientId, token) && i.isValidTimeStamp(clientId, token)
}

func (i *Identity) isValidToken(clientId string, token string) bool {
	hash := i.ParsedTokens[clientId].Token
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(token))
	if err != nil {
		log.Debug("passwords do not match", err)
		return false
	}
	return true
}
func (i *Identity) isValidTimeStamp(clientId, token string) bool {
	expiry := i.ParsedTokens[clientId].Expiry
	t := time.Unix(expiry, -1)
	log.Info(t)
	return time.Now().Before(t)
}

func initDbConnection() {
	db, err := sqlx.Open("mysql", "root@/moove_development")
	if err != nil {
		panic(err)
	}
	Db = db
}

func main() {
	initDbConnection()
	defer func() {
		err := Db.Close()
		if err != nil {
			panic(err)
		}
	}()
	identity := FetchIdentityByUID(Db, "moove.dinesh1651@gmail.com")
	if identity.isValid("GWyFWywJ3mf5TJe77DBIHw", "v-pWLlS4pwf_L9qbx-FtYw") {
		log.Info("Yay! passwords match")
	}
}

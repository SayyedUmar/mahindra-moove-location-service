package identity

import (
	"bytes"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Token struct {
	Token  string `json:"token"`
	Expiry int64  `json:"expiry"`
}
type Identity struct {
	Id           int
	Uid          string
	Tokens       string
	ParsedTokens map[string]Token
}

func FetchIdentityByUID(db *sqlx.DB, uid string) *Identity {
	var i Identity
	row := db.QueryRowx("select id, uid, tokens from users where uid=?", uid)
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

func (i *Identity) IsValid(clientId, token string) bool {
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
	return time.Now().Before(t)
}

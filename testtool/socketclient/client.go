package socketclient

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MOOVE-Network/location_service/db"
	"github.com/MOOVE-Network/location_service/identity"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type AuthHeader struct {
	ClientID string
	UID      string
	Token    string
}

func (ah AuthHeader) GoString() string {
	return fmt.Sprintf("%s - %s - %s", ah.ClientID, ah.UID, ah.Token)
}
func (ah AuthHeader) HttpHeader() http.Header {
	hdr := make(http.Header)
	hdr["uid"] = []string{ah.UID}
	hdr["access-token"] = []string{ah.Token}
	hdr["client"] = []string{ah.ClientID}
	return hdr
}

func SetupWSConnection(driverID int) *websocket.Conn {
	ah := FetchValidDriverRequestHeader(db.CurrentDB(), driverID)
	conn, _, err := websocket.DefaultDialer.Dial(
		fmt.Sprintf("ws://localhost:4343/api/v3/drivers/%d/location", driverID),
		ah.HttpHeader())
	if err != nil {
		panic(err)
	}
	return conn
}

func FetchValidDriverRequestHeader(conn *sqlx.DB, userID int) AuthHeader {
	ident := identity.FetchIdentityByID(conn, userID)
	var clientID string
	var uid string
	var token string
	for k, v := range ident.ParsedTokens {
		expiry := time.Unix(v.Expiry, -1)
		if time.Now().Before(expiry) {
			fmt.Println("found an appropriate token")
			clientID = k
			uid = ident.Uid
			tkn, err := bcrypt.GenerateFromPassword([]byte(v.Token), bcrypt.DefaultCost)
			if err != nil {
				panic(err)
			}
			token = string(tkn)
			break
		}
	}
	return AuthHeader{ClientID: clientID, UID: uid, Token: token}
}

package web

import (
	//"github.com/MOOVE-Network/location_service/identity"
	"github.com/gorilla/websocket"
	//log "github.com/sirupsen/logrus"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//var hub := NewHub()
func LocationSocket(w http.ResponseWriter, r *http.Request) {
	//ident := r.Context().Value("identity").(*identity.Identity)
	//conn, err := upgrader.Upgrade(w, r, nil)
	//if err != nil {
	//log.Error(err)
	//}
	//log.Info(conn)
}

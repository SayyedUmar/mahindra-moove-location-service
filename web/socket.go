package web

import (
	"fmt"
	"github.com/MOOVE-Network/location_service/identity"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var hub = NewHub()

func LocationSocket(w http.ResponseWriter, r *http.Request) {
	ident := r.Context().Value("identity").(*identity.Identity)
	log.Debugf("headers %v", r.Header)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	client := NewClient(hub, conn, ident.Id)
	hub.Register <- client
	for {
		select {
		case message := <-client.Receive:
			log.Infof("from server %s \n", message)
			client.Send <- []byte(fmt.Sprintf("Received from server %s", message))
		}
	}
}

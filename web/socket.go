package web

import (
	"github.com/MOOVE-Network/location_service/identity"
	"github.com/MOOVE-Network/location_service/socketstore"
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
	go readMessages(client)
}
func readMessages(client *Client) {
	for {
		message := <-client.Receive
		switch socketstore.MessageType(message) {
		case "LOCATION":
			locationUpdate, err := socketstore.LocationUpdateFromJSON(message)
			if err != nil {
				log.Warnf("Unable to decode location update message %s", string(message))
			}
			client.hub.Send(string(locationUpdate.TripID), message)
			client.hub.Send(string(client.ID), message)
		case "HEARTBEAT":
			log.Infof("received heartbeat %s", message)
		default:
			log.Warnf("Unknown message type detected %s", string(message))
		}
		log.Infof("got message %s", message)
	}
}

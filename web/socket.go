package web

import (
	"github.com/MOOVE-Network/location_service/db"
	"github.com/MOOVE-Network/location_service/identity"
	"github.com/MOOVE-Network/location_service/socketstore"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(*http.Request) bool { return true },
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
			log.Infof("got location %s ", message)
			locationUpdate, err := socketstore.LocationUpdateFromJSON(message)
			if err != nil {
				log.Warnf("Unable to decode location update message %s", string(message))
			}
			client.hub.Send(string(locationUpdate.TripID), message)
			client.hub.Send(string(client.ID), message)
		case "HEARTBEAT":
			hb := &db.HeartBeat{UserID: client.ID, UpdatedAt: time.Now()}
			heartBeats[client.ID] = hb
		default:
			log.Warnf("Unknown message type detected %s", string(message))
		}
	}
}

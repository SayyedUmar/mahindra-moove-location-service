package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/MOOVE-Network/location_service/db"
	"github.com/MOOVE-Network/location_service/identity"
	"github.com/MOOVE-Network/location_service/socketstore"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
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
				continue
			}
			tlMutex.Lock()
			tripLocations = append(tripLocations, locationUpdate.ToTripLocation)
			tlMutex.Unlock()
			client.hub.Send(strconv.Itoa(locationUpdate.TripID), message)
			client.hub.Send(strconv.Itoa(client.ID), message)
		case "HEARTBEAT":
			hb := &db.HeartBeat{UserID: client.ID, UpdatedAt: time.Now()}
			hbMutex.Lock()
			heartBeats[client.ID] = hb
			hbMutex.Unlock()
		case "SUBSCRIBE":
			subscription, err := socketstore.SubscribeEventFromJSON(message)
			if err != nil {
				log.Warnf("Unable to decode subscription %s", string(message))
			}
			hub.Subscribe(subscription.Topic, client)
		default:
			log.Warnf("Unknown message type detected %s", string(message))
		}
	}
}

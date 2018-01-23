package web

import (
	"bytes"
	"encoding/json"
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

func acknowledge(wsMsg socketstore.WsMessage, sendChan chan<- []byte) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	err := enc.Encode(wsMsg)
	if err != nil {
		log.Error("Error acknowledging ", err)
	} else {
		sendChan <- buf.Bytes()
	}
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func parseLocation(msg []byte) (Location, error) {
	var loc Location
	err := json.Unmarshal(msg, &loc)
	if err != nil {
		return Location{}, err
	}
	return loc, nil
}
func readMessages(client *Client) {
	for {
		message := <-client.Receive
		wsMsg, _ := socketstore.WsMessageFromJSON(message)
		switch wsMsg.EventType {
		case "LOCATION":
			log.Infof("got location %s ", message)
			go acknowledge(wsMsg, client.Send)
			locationUpdate, err := socketstore.LocationUpdateFromJSON(message)
			if err != nil {
				log.Warnf("Unable to decode location update message %s", string(message))
				continue
			}
			tlMutex.Lock()
			tripLocations = append(tripLocations, locationUpdate.ToTripLocation())
			tlMutex.Unlock()
			client.hub.Send(strconv.Itoa(locationUpdate.TripID), message)
			client.hub.Send(strconv.Itoa(client.ID), message)
		case "HEARTBEAT":
			loc, err := parseLocation(message)
			var hb *db.HeartBeat
			if err != nil {
				hb = &db.HeartBeat{UserID: client.ID, UpdatedAt: time.Now()}
			} else {
				hb = &db.HeartBeat{UserID: client.ID, UpdatedAt: time.Now(), Lat: loc.Lat, Lng: loc.Lng}
			}
			hbMutex.Lock()
			heartBeats[client.ID] = hb
			hbMutex.Unlock()
		case "SUBSCRIBE":
			go acknowledge(wsMsg, client.Send)
			subscription, err := socketstore.SubscribeEventFromJSON(message)
			if err != nil {
				log.Warnf("Unable to decode subscription %s", string(message))
			}
			hub.Subscribe(subscription.Topic, client)
		case "GEOFENCE":
			go acknowledge(wsMsg, client.Send)
			log.Infof("Got GEOFENCE event %s\n", message)
		default:
			log.Warnf("Unknown message type detected %s", string(message))
		}
	}
}

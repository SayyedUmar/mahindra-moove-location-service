package web

import (
	"bytes"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer 2KB
	maxMessageSize = 2048
)

type subscription struct {
	topic  string
	client *Client
}
type message struct {
	topic   string
	message []byte
}

// Hub serves a hub to hold connections to multiple clients
type Hub struct {
	topics      map[string][]*Client
	Clients     map[*Client]bool
	Register    chan *Client
	Unregister  chan *Client
	topicChan   chan subscription
	messageChan chan message
}

// NewHub creates a new hub
func NewHub() *Hub {
	return &Hub{
		topics:      make(map[string][]*Client),
		Clients:     make(map[*Client]bool),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		topicChan:   make(chan subscription),
		messageChan: make(chan message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case subs := <-h.topicChan:
			h.topics[subs.topic] = append(h.topics[subs.topic], subs.client)
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				close(client.Send)
				delete(h.Clients, client)
			}
		case message := <-h.messageChan:
			if message.topic == "*" {
				for k, _ := range h.Clients {
					k.Send <- message.message
				}
			}
			if clients, ok := h.topics[message.topic]; ok {
				for _, client := range clients {
					client.Send <- message.message
				}
			}
		}
	}
}

func (h *Hub) Subscribe(topic string, c *Client) {
	h.topicChan <- subscription{topic: topic, client: c}
}

func (h *Hub) Send(topic string, data []byte) {
	h.messageChan <- message{topic: topic, message: data}
}

type Client struct {
	hub     *Hub
	ID      int
	Send    chan []byte
	Receive chan []byte
	conn    *websocket.Conn
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func NewClient(hub *Hub, conn *websocket.Conn, id int) *Client {
	client := &Client{
		hub:     hub,
		conn:    conn,
		ID:      id,
		Send:    make(chan []byte, 256),
		Receive: make(chan []byte),
	}
	go client.readPump()
	go client.writePump()
	return client
}

func (c *Client) readPump() {
	defer func() {
		c.hub.Unregister <- c
		err := c.conn.Close()
		if err != nil {
			log.Warn("Error closing websocket client connection")
		}
	}()
	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
		log.Error(err)
		return err
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Error(err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Errorf("Error reading from socket: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.Receive <- message
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		err := c.conn.Close()
		if err != nil {
			log.Warnf("Error closing websocket connection : %v\n", err)
		}
	}()
	for {
		select {
		case message, ok := <-c.Send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// the hub closed the channel
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Warn("unable to get websocket writer")
				return
			}
			_, err = w.Write(message)
			if err != nil {
				log.Warn("unable to write to socket", err)
			}
			if err := w.Close(); err != nil {
				log.Error("unable to close writer", err)
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Warn("write error", err)
				return
			}
		}
	}
}

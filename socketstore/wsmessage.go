package socketstore

import (
	"bytes"
	"encoding/json"
)

type WsMessage struct {
	ID        string `json:"id"`
	EventType string `json:"eventType"`
}

func WsMessageFromJSON(msg []byte) (WsMessage, error) {
	var wsMsg WsMessage
	buffer := bytes.NewBuffer(msg)
	decoder := json.NewDecoder(buffer)
	err := decoder.Decode(&wsMsg)
	if err != nil {
		return WsMessage{}, err
	}
	return wsMsg, nil
}

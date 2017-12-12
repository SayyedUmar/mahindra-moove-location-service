package socketstore

import (
	"bytes"
	"encoding/json"
)

type WsMessage struct {
	EventType string `json:"eventType"`
}

func MessageType(msg []byte) string {
	var wsMsg WsMessage
	buffer := bytes.NewBuffer(msg)
	decoder := json.NewDecoder(buffer)
	err := decoder.Decode(&wsMsg)
	if err != nil {
		return ""
	}
	return wsMsg.EventType
}

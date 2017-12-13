package socketstore

import (
	"bytes"
	"encoding/json"
)

type SubscribeEvent struct {
	Topic string `json:"topic"`
}

func SubscribeEventFromJSON(msg []byte) (*SubscribeEvent, error) {
	var subscription SubscribeEvent
	buffer := bytes.NewBuffer(msg)
	decoder := json.NewDecoder(buffer)
	err := decoder.Decode(&subscription)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

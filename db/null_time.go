package db

import (
	"encoding/json"
	"time"
)

type NullTime struct {
	Valid bool
	Value time.Time
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return json.Marshal(nt.Value)
	}
	return json.Marshal(nil)
}

func (nt NullTime) UnmarshalJSON(data []byte) error {
	var val *time.Time
	if err := json.Unmarshal(data, val); err != nil {
		return err
	}
	if val != nil {
		nt.Valid = true
		nt.Value = *val
	} else {
		nt.Valid = false
	}
	return nil
}

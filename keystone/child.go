package keystone

import (
	"encoding/json"
)

type Child struct {
	Type    string
	ID      string
	Data    []byte
	written bool
}

func (e *Entity) AddChild(childType string, data []byte) {
	e.Children = append(e.Children, Child{
		written: false,
		Type:    childType,
		Data:    data,
	})
}

func (e *Entity) AddChildJson(childType string, data interface{}) error {
	if dataBytes, err := json.Marshal(data); err == nil {
		e.AddChild(childType, dataBytes)
		return nil
	} else {
		return err
	}
}

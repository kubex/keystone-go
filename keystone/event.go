package keystone

import (
	"time"
)

type Event struct {
	written bool
	ID      string            `json:"-"`
	Type    string            `json:"e"`
	Time    time.Time         `json:"-"`
	Data    map[string]string `json:"d"`
	Actor   *Actor            `json:"a"`
}

func (e *Entity) AddEvent(eventType string, properties map[string]string) {
	e.Events = append(e.Events, Event{
		written: false,
		Type:    eventType,
		Time:    time.Now(),
		Data:    properties,
	})
}

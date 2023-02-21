package keystone

import (
	"github.com/kubex/definitions-go/app"
	"time"
)

type Event struct {
	written bool
	ID      string        `json:"-"`
	Type    app.ScopedKey `json:"e"`
	Time    time.Time     `json:"-"`
	Data    []Property    `json:"d"`
	Actor   *Actor        `json:"a"`
}

func (e *Entity) AddEvent(eventType string, properties ...Property) {
	e.Events = append(e.Events, Event{
		written: false,
		Type:    app.NewScopedKey(eventType, defaultSetGlobalAppID),
		Time:    time.Now(),
		Data:    properties,
	})
}

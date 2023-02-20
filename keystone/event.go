package keystone

import (
	"github.com/kubex/definitions-go/app"
	"time"
)

type Event struct {
	ID    string        `json:"-"`
	Type  app.ScopedKey `json:"e"`
	Time  time.Time     `json:"-"`
	Data  []Property    `json:"d"`
	Actor *Actor        `json:"a"`
}

func (e *Entity) AddEvent(eventType string, properties ...Property) {
	e.Events = append(e.Events, Event{
		Type: app.NewScopedKey(eventType, defaultSetGlobalAppID),
		Time: time.Now(),
		Data: properties,
	})
}

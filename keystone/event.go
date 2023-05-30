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

type EntityEventProvider interface {
	ClearEvents() error
	GetEvents() ([]Event, error)
}

type EntityEvents struct {
	Events []Event
}

func (e *EntityEvents) ClearEvents() error {
	e.Events = []Event{}
	return nil
}

func (e *EntityEvents) GetEvents() ([]Event, error) {
	return e.Events, nil
}

func (e *EntityEvents) AddEvent(eventType string, properties map[string]string) {
	e.Events = append(e.Events, Event{
		written: false,
		Type:    eventType,
		Time:    time.Now(),
		Data:    properties,
	})
}

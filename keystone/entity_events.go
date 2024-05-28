package keystone

import (
	"github.com/kubex/keystone-go/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// EntityEventProvider is an interface for entities that can have events
type EntityEventProvider interface {
	ClearKeystoneEvents() error
	GetKeystoneEvents() []*proto.EntityEvent
}

// EntityEvents is a struct that implements EntityEventProvider
type EntityEvents struct {
	ksEntityEvents []*proto.EntityEvent
}

// ClearKeystoneEvents clears the events
func (e *EntityEvents) ClearKeystoneEvents() error {
	e.ksEntityEvents = []*proto.EntityEvent{}
	return nil
}

// GetKeystoneEvents returns the events
func (e *EntityEvents) GetKeystoneEvents() []*proto.EntityEvent {
	return e.ksEntityEvents
}

// AddKeystoneEvent adds an event
func (e *EntityEvents) AddKeystoneEvent(eventType string, properties map[string]string) {
	e.ksEntityEvents = append(e.ksEntityEvents, &proto.EntityEvent{
		Type: &proto.Key{Key: eventType},
		Time: timestamppb.New(time.Now()),
		Data: properties,
	})
}

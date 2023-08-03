package keystone

import (
	"github.com/kubex/keystone-go/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type EntityEventProvider interface {
	ClearKeystoneEvents() error
	GetKeystoneEvents() ([]*proto.EntityEvent, error)
}

type EntityEvents struct {
	ksEntityEvents []*proto.EntityEvent
}

func (e *EntityEvents) ClearKeystoneEvents() error {
	e.ksEntityEvents = []*proto.EntityEvent{}
	return nil
}

func (e *EntityEvents) GetKeystoneEvents() ([]*proto.EntityEvent, error) {
	return e.ksEntityEvents, nil
}

func (e *EntityEvents) AddKeystoneEvent(eventType string, properties map[string]string) {
	e.AddKeystoneEventWithDedupe("", eventType, properties)
}

func (e *EntityEvents) AddKeystoneEventWithDedupe(dedupeKey, eventType string, properties map[string]string) {
	e.ksEntityEvents = append(e.ksEntityEvents, &proto.EntityEvent{
		Type: &proto.Key{Key: eventType},
		Tid:  dedupeKey,
		Time: timestamppb.New(time.Now()),
		Data: properties,
	})
}

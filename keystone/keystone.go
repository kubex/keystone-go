package keystone

import (
	"github.com/kubex/definitions-go/k4"
)

type Entity struct {
	WorkspaceID      string
	ID               k4.ID
	Properties       map[string]Property
	DeleteProperties map[string]Property
	LogEntries       []LogEntry
	Events           []Event
}

func NewEntity(workspaceID string) *Entity {
	return &Entity{
		WorkspaceID: workspaceID,
	}
}

func ExistingEntity(workspaceID, entityID string) *Entity {
	return &Entity{
		WorkspaceID: workspaceID,
		ID:          k4.IDFromString(entityID),
	}
}

func (e *Entity) Mutate(prop ...Property) {
	for _, p := range prop {
		e.Properties[p.Name.String()] = p
	}
}

func (e *Entity) Delete(prop Property) {
	e.DeleteProperties[prop.Name.String()] = prop
	delete(e.Properties, prop.Name.String())
}

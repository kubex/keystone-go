package keystone

import (
	"github.com/kubex/definitions-go/app"
	"github.com/kubex/definitions-go/k4"
)

type Entity struct {
	WorkspaceID      string
	ID               k4.ID
	Schema           app.ScopedKey
	Properties       map[string]Property
	DeleteProperties map[string]Property
	LogEntries       []LogEntry
	Events           []Event
	Children         []Child
}

func newEntity(workspaceID string) *Entity {
	return &Entity{
		WorkspaceID:      workspaceID,
		Properties:       make(map[string]Property),
		DeleteProperties: make(map[string]Property),
	}
}

func NewEntity(workspaceID, schema string) *Entity {
	e := newEntity(workspaceID)
	e.Schema = app.NewScopedKey(schema, defaultSetGlobalAppID)
	return e
}

func ExistingEntity(workspaceID, entityID string) *Entity {
	e := newEntity(workspaceID)
	e.ID = k4.IDFromString(entityID)
	return e
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

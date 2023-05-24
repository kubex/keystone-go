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
	Relationships    []Relationship
}

func (w Workspace) newEntity() *Entity {
	return &Entity{
		WorkspaceID:      w.workspaceID,
		Properties:       make(map[string]Property),
		DeleteProperties: make(map[string]Property),
	}
}

func (e *Entity) Mutate(prop ...Property) {
	for _, p := range prop {
		e.Properties[p.Name] = p
	}
}

func (e *Entity) Delete(prop Property) {
	e.DeleteProperties[prop.Name] = prop
	delete(e.Properties, prop.Name)
}

type Workspace struct {
	workspaceID string
}

func NewWorkspace(workspaceID string) Workspace {
	return Workspace{workspaceID: workspaceID}
}

func (w Workspace) ID() string { return w.workspaceID }

func (w Workspace) NewEntity(schema string) *Entity {
	e := w.newEntity()
	e.Schema = app.NewScopedKey(schema, defaultSetGlobalAppID)
	return e
}

func (w Workspace) ExistingEntity(entityID string) *Entity {
	e := w.newEntity()
	e.ID = k4.IDFromString(entityID)
	return e
}

package keystone

import "strings"

type BaseEntity struct {
	EntityEvents
	EntityLabels
	EntityLinks
	EntityLogger
	EntityRelationships

	_entityID string
}

func (e *BaseEntity) GetKeystoneID() string {
	return e._entityID
}

func (e *BaseEntity) SetKeystoneID(id string) {
	e._entityID = id
}

type testEntity struct {
	BaseEntity
}

type Entity interface {
	GetKeystoneID() string
	SetKeystoneID(id string)
}

type ChildEntity interface {
	GetKeystoneParentID() string
}

type BaseChildEntity struct {
	BaseEntity
}

func (e *BaseChildEntity) GetKeystoneParentID() string {
	split := strings.Split(e._entityID, "-")
	return split[0]
}

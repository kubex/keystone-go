package keystone

import (
	"github.com/kubex/keystone-go/proto"
	"strings"
)

func RemoteEntity(entityID string) *BaseEntity {
	return &BaseEntity{
		_entityID: entityID,
	}
}

type BaseEntity struct {
	EntityEvents
	EntityLabels
	EntityLogger
	EntityRelationships
	EntitySensors
	EntityLock
	EntityDetails
	_lastLoad *proto.EntityResponse
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
	GetKeystoneChildID() string
	SetKeystoneParentID(id string)
	SetKeystoneChildID(id string)
}

type BaseChildEntity struct {
	BaseEntity
	_parentID string
	_childID  string
}

func (e *BaseChildEntity) SetKeystoneID(id string) {
	e._entityID = id
	split := strings.Split(e._entityID, "-")
	e._parentID = split[0]
	if len(split) > 1 {
		e._childID = split[1]
	}
}

func (e *BaseChildEntity) SetKeystoneParentID(id string) {
	if strings.Contains(id, "-") {
		e.SetKeystoneID(id)
	} else {
		e._parentID = id
	}

	if e._entityID == "" {
		e._entityID = e._parentID
	}
}

func (e *BaseChildEntity) SetKeystoneChildID(id string) {
	e._childID = id
}

func (e *BaseChildEntity) GetKeystoneParentID() string {
	if e._parentID == "" {
		split := strings.Split(e._entityID, "-")
		e._parentID = split[0]
	}
	return e._parentID
}

func (e *BaseChildEntity) GetKeystoneChildID() string {
	if e._childID == "" {
		split := strings.Split(e._entityID, "-")
		if len(split) < 2 {
			return ""
		}

		e._childID = split[1]
	}
	return e._childID
}

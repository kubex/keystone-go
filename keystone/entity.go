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
}

type BaseChildEntity struct {
	BaseEntity
}

func (e *BaseChildEntity) GetKeystoneParentID() string {
	split := strings.Split(e._entityID, "-")
	return split[0]
}

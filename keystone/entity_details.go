package keystone

import (
	"github.com/kubex/keystone-go/proto"
	"time"
)

type EntityDetail interface {
	SetEntityDetail(entity *proto.Entity)
}

type EntityDetails struct {
	ksCreated     time.Time
	ksStateChange time.Time
	ksState       proto.EntityState
	ksLastUpdate  time.Time
}

func (e *EntityDetails) DateCreated() time.Time           { return e.ksCreated }
func (e *EntityDetails) LastUpdated() time.Time           { return e.ksLastUpdate }
func (e *EntityDetails) KeystoneState() proto.EntityState { return e.ksState }

func (e *EntityDetails) SetEntityDetail(entity *proto.Entity) {
	if entity == nil {
		return
	}

	e.ksCreated = entity.GetCreated().AsTime()
	e.ksStateChange = entity.GetStateChange().AsTime()
	e.ksState = entity.GetState()
	e.ksLastUpdate = entity.GetLastUpdate().AsTime()
}

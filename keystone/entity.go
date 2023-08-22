package keystone

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

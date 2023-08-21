package keystone

type BaseEntity struct {
	EntityEvents
	EntityLabels
	EntityLinks
	EntityLogger
	EntityRelationships

	entityID string
}

func (e *BaseEntity) GetKeystoneID() string {
	return e.entityID
}

func (e *BaseEntity) SetKeystoneID(id string) {
	e.entityID = id
}

type testEntity struct {
	BaseEntity
}

type Entity interface {
	GetKeystoneID() string
	SetKeystoneID(id string)
}

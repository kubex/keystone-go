package keystone

type Entity struct {
	EntityEvents
	EntityLabels
	EntityLinks
	EntityLogger
	EntityRelationships

	entityID string
}

func (e *Entity) GetKeystoneID() string {
	return e.entityID
}

package keystone

import (
	"time"
)

type Relationship struct {
	Relationship  string
	DestinationID string
	Since         time.Time
	Data          map[string]string
	written       bool
}

func (e *Entity) Edge(rel Relationship) {
	e.Relationships = append(e.Relationships, rel)
}

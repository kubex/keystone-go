package keystone

import (
	"github.com/kubex/definitions-go/app"
	"time"
)

type Relationship struct {
	Relationship  app.ScopedKey
	DestinationID string
	Since         time.Time
	Data          []Property
	written       bool
}

func (e *Entity) Edge(rel Relationship) {
	e.Relationships = append(e.Relationships, rel)
}

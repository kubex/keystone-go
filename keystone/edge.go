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
}

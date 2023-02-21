package keystone

import (
	"github.com/kubex/definitions-go/app"
)

type Child struct {
	Type    app.ScopedKey
	ID      string
	Data    []byte
	written bool
}

func (e *Entity) AddChild(childType string, data []byte) {
	e.Children = append(e.Children, Child{
		written: false,
		Type:    app.NewScopedKey(childType, defaultSetGlobalAppID),
		Data:    data,
	})
}

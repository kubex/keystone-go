package keystone

import "github.com/kubex/definitions-go/app"

type Child struct {
	Type    app.ScopedKey
	ID      string
	Data    []byte
	written bool
}

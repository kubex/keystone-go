package keystone

import "github.com/kubex/keystone-go/proto"

type EntityDefinition interface {
	GetKeystoneDefinition() TypeDefinition
}

type TypeDefinition struct {
	Type        string // Unique Type Name e.g. library/user
	Name        string // Friendly name of the entity e.g. Library User
	Description string // Description of the entity
	Singular    string // Name for a single one of these entities e.g. User
	Plural      string // Name for a collection of these entities e.g. Users
	Options     []proto.Schema_Option
}

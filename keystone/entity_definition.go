package keystone

import "github.com/kubex/keystone-go/proto"

// EntityDefinition is an interface that defines the keystone entity
type EntityDefinition interface {
	GetKeystoneDefinition() TypeDefinition
}

// TypeDefinition is a definition of a keystone type
type TypeDefinition struct {
	Type           string // Unique Type Name e.g. library/user
	Name           string // Friendly name of the entity e.g. Library User
	Description    string // Description of the entity
	Singular       string // Name for a single one of these entities e.g. User
	Plural         string // Name for a collection of these entities e.g. Users
	Options        []proto.Schema_Option
	ActiveDataSets []*proto.ADS
}

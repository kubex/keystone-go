package keystone

import (
	"time"

	"github.com/kubex/keystone-go/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// EntityRelationshipProvider is an interface for entities that can have relationships
type EntityRelationshipProvider interface {
	ClearKeystoneRelationships() error
	GetKeystoneRelationships() []*proto.EntityRelationship
	SetKeystoneRelationships(links []*proto.EntityRelationship)
}

// EntityRelationships is a struct that implements EntityRelationshipProvider
type EntityRelationships struct {
	ksEntityRelationships []*proto.EntityRelationship
}

// ClearKeystoneRelationships clears the relationships
func (e *EntityRelationships) ClearKeystoneRelationships() error {
	e.ksEntityRelationships = []*proto.EntityRelationship{}
	return nil
}

// GetKeystoneRelationships returns the relationships
func (e *EntityRelationships) GetKeystoneRelationships() []*proto.EntityRelationship {
	return e.ksEntityRelationships
}

// SetKeystoneRelationships sets the relationships
func (e *EntityRelationships) SetKeystoneRelationships(links []*proto.EntityRelationship) {
	e.ksEntityRelationships = links
}

// AddKeystoneRelationship adds a relationship
func (e *EntityRelationships) AddKeystoneRelationship(source, target string, meta map[string]string, since time.Time) {
	e.ksEntityRelationships = append(e.ksEntityRelationships, &proto.EntityRelationship{
		Relationship: &proto.Key{Key: source},
		TargetId:     target,
		Data:         meta,
		Since:        timestamppb.New(since),
	})
}

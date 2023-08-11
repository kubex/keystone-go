package keystone

import (
	"time"

	"github.com/kubex/keystone-go/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EntityRelationshipProvider interface {
	ClearKeystoneRelationships() error
	GetKeystoneRelationships() []*proto.EntityRelationship
}

type EntityRelationships struct {
	ksEntityRelationships []*proto.EntityRelationship
}

func (e *EntityRelationships) ClearKeystoneRelationships() error {
	e.ksEntityRelationships = []*proto.EntityRelationship{}
	return nil
}

func (e *EntityRelationships) GetKeystoneRelationships() []*proto.EntityRelationship {
	return e.ksEntityRelationships
}

func (e *EntityRelationships) AddKeystoneRelationship(source, target string, meta map[string]string, since time.Time) {
	e.ksEntityRelationships = append(e.ksEntityRelationships, &proto.EntityRelationship{
		Relationship: &proto.Key{Key: source},
		TargetId:     target,
		Data:         meta,
		Since:        timestamppb.New(since),
	})
}

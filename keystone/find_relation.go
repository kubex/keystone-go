package keystone

import "github.com/kubex/keystone-go/proto"

type relationOf struct {
	source  string
	relType *proto.Key
}

func (f relationOf) Apply(config *filterRequest) {
	config.RelationOf = &proto.RelationOf{
		SourceId:     f.source,
		Relationship: f.relType,
	}
}

func RelationOf(entityID, relationshipType, relVendor, relApp string) FindOption {
	return relationOf{
		source: entityID,
		relType: &proto.Key{
			Source: &proto.VendorApp{
				VendorId: relVendor,
				AppId:    relApp,
			},
			Key: relationshipType,
		},
	}
}

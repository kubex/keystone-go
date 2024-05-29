package keystone

import "github.com/kubex/keystone-go/proto"

type relationOf struct {
	source      string
	destination string
	relType     *proto.Key
}

func (f relationOf) Apply(config *filterRequest) {
	config.RelationOf = &proto.RelationOf{
		SourceId:      f.source,
		DestinationId: f.destination,
		Relationship:  f.relType,
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

func RelationTo(entityID, relationshipType, relVendor, relApp string) FindOption {
	return relationOf{
		destination: entityID,
		relType: &proto.Key{
			Source: &proto.VendorApp{
				VendorId: relVendor,
				AppId:    relApp,
			},
			Key: relationshipType,
		},
	}
}

func RelationOfSibling(entityID, relationshipType string) FindOption {
	return RelationOf(entityID, relationshipType, "", "")
}

func RelationToSibling(entityID, relationshipType string) FindOption {
	return RelationTo(entityID, relationshipType, "", "")
}

package keystone

import (
	"github.com/kubex/keystone-go/proto"
)

type Retriever interface {
	BaseRequest() *proto.EntityRequest
}

type RetrieveByEntityID struct {
	EntityID string
}

func (l RetrieveByEntityID) BaseRequest() *proto.EntityRequest {
	return &proto.EntityRequest{
		EntityId: l.EntityID,
	}
}

type RetrieveByUnique struct {
	UniqueID string
	Property string
}

func (l RetrieveByUnique) BaseRequest() *proto.EntityRequest {
	return &proto.EntityRequest{
		UniqueId: &proto.IDLookup{
			SchemaId: "SCHEMAID",
			Property: l.Property,
			UniqueId: l.UniqueID,
		},
	}
}

type RetrieveOption interface {
	Apply(config *proto.EntityRequest)
}

func WithProperties(properties ...string) RetrieveOption {
	return propertyLoader{properties: properties}
}

func WithDecryptedProperties(properties ...string) RetrieveOption {
	return propertyLoader{properties: properties, decrypt: true}
}

func WithLinks(links []string) RetrieveOption {
	return linksLoader{Links: links}
}

func WithAuthorization(auth *proto.Authorization) RetrieveOption {
	return authorizationLoader{auth: auth}
}

func WithSummary() RetrieveOption {
	return summaryLoader{summary: true}
}

func WithDatum() RetrieveOption {
	return datumLoader{datum: true}
}

func WithLabels() RetrieveOption {
	return labelLoader{labels: true}
}

func WithChildren(childType string, ids ...string) RetrieveOption {
	return childrenLoader{childType: childType, ids: ids}
}

type propertyLoader struct {
	properties []string
	decrypt    bool
}

func (l propertyLoader) Apply(config *proto.EntityRequest) {
	if config.Properties == nil {
		config.Properties = make([]*proto.PropertyRequest, 0)
	}

	config.Properties = append(config.Properties, &proto.PropertyRequest{Properties: l.properties, Decrypt: l.decrypt})
}

type linksLoader struct{ Links []string }

func (l linksLoader) Apply(config *proto.EntityRequest) {
	if config.LinkByType == nil {
		config.LinkByType = make([]*proto.Key, 0)
	}
	for _, link := range l.Links {
		config.LinkByType = append(config.LinkByType, &proto.Key{Key: link})
	}
}

type authorizationLoader struct{ auth *proto.Authorization }

func (l authorizationLoader) Apply(config *proto.EntityRequest) { config.Authorization = l.auth }

type summaryLoader struct{ summary bool }

func (l summaryLoader) Apply(config *proto.EntityRequest) { config.Summary = l.summary }

type datumLoader struct{ datum bool }

func (l datumLoader) Apply(config *proto.EntityRequest) { config.Datum = l.datum }

type labelLoader struct{ labels bool }

func (l labelLoader) Apply(config *proto.EntityRequest) { config.Labels = l.labels }

type childrenLoader struct {
	childType string
	ids       []string
}

func (l childrenLoader) Apply(config *proto.EntityRequest) {
	if config.Children == nil {
		config.Children = make([]*proto.ChildRequest, 0)
	}

	config.Children = append(config.Children, &proto.ChildRequest{
		Type: &proto.Key{Key: l.childType},
		Cid:  l.ids,
	})
}

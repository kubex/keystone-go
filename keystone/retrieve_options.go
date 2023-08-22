package keystone

import (
	"github.com/kubex/keystone-go/proto"
)

// Retriever is an interface that defines a retriever
type Retriever interface {
	BaseRequest() *proto.EntityRequest
}

// RetrieveByEntityID is a retriever that retrieves an entity by its ID
type RetrieveByEntityID struct {
	EntityID string
}

// BaseRequest returns the base RetrieveByEntityID request
func (l RetrieveByEntityID) BaseRequest() *proto.EntityRequest {
	return &proto.EntityRequest{
		EntityId: l.EntityID,
		View:     &proto.EntityView{},
	}
}

// RetrieveByUnique is a retriever that retrieves an entity by its unique ID
type RetrieveByUnique struct {
	UniqueID string
	Property string
}

// BaseRequest returns the base RetrieveByUnique request
func (l RetrieveByUnique) BaseRequest() *proto.EntityRequest {
	return &proto.EntityRequest{
		UniqueId: &proto.IDLookup{
			SchemaId: "SCHEMAID",
			Property: l.Property,
			UniqueId: l.UniqueID,
		},
	}
}

// RetrieveOption is an interface for options to be applied to an entity request
type RetrieveOption interface {
	Apply(config *proto.EntityRequest)
}

// WithProperties is a retrieve option that loads properties
func WithView(name string) RetrieveOption {
	return viewName{name: name}
}

// WithProperties is a retrieve option that loads properties
func WithProperties(properties ...string) RetrieveOption {
	return propertyLoader{properties: properties}
}

// WithDecryptedProperties is a retrieve option that loads decrypted properties
func WithDecryptedProperties(properties ...string) RetrieveOption {
	return propertyLoader{properties: properties, decrypt: true}
}

// WithProperty is a retrieve option that loads properties
func WithProperty(decrypt bool, properties ...string) RetrieveOption {
	return propertyLoader{properties: properties, decrypt: decrypt}
}

// WithLinks is a retrieve option that loads links
func WithLinks(links ...string) RetrieveOption {
	return linksLoader{Links: links}
}

// WithRelationships is a retrieve option that loads relationships
func WithRelationships(keys ...string) RetrieveOption {
	return relationshipsLoader{keys: keys}
}

// WithSummary is a retrieve option that loads summaries
func WithSummary() RetrieveOption {
	return summaryLoader{summary: true}
}

// WithDatum is a retrieve option that loads datum
func WithDatum() RetrieveOption {
	return datumLoader{datum: true}
}

// WithLabels is a retrieve option that loads labels
func WithLabels() RetrieveOption {
	return labelLoader{labels: true}
}

// WithChildren is a retrieve option that loads children
func WithChildren(childType string, ids ...string) RetrieveOption {
	return childrenLoader{childType: childType, ids: ids}
}

type propertyLoader struct {
	properties []string
	decrypt    bool
}

func (l propertyLoader) Apply(config *proto.EntityRequest) {
	if config.View.Properties == nil {
		config.View.Properties = make([]*proto.PropertyRequest, 0)
	}

	config.View.Properties = append(config.View.Properties, &proto.PropertyRequest{Properties: l.properties, Decrypt: l.decrypt})
}

type viewName struct {
	name string
}

func (l viewName) Apply(config *proto.EntityRequest) {
	config.View.Name = l.name
}

type linksLoader struct{ Links []string }

func (l linksLoader) Apply(config *proto.EntityRequest) {
	if config.View.LinkByType == nil {
		config.View.LinkByType = make([]*proto.Key, 0)
	}
	for _, link := range l.Links {
		config.View.LinkByType = append(config.View.LinkByType, &proto.Key{Key: link})
	}
}

type relationshipsLoader struct{ keys []string }

func (l relationshipsLoader) Apply(config *proto.EntityRequest) {
	if config.View.RelationshipByType == nil {
		config.View.RelationshipByType = make([]*proto.Key, 0)
	}
	for _, key := range l.keys {
		config.View.RelationshipByType = append(config.View.RelationshipByType, &proto.Key{Key: key})
	}
}

type summaryLoader struct{ summary bool }

func (l summaryLoader) Apply(config *proto.EntityRequest) { config.View.Summary = l.summary }

type datumLoader struct{ datum bool }

func (l datumLoader) Apply(config *proto.EntityRequest) { config.View.Datum = l.datum }

type labelLoader struct{ labels bool }

func (l labelLoader) Apply(config *proto.EntityRequest) { config.View.Labels = l.labels }

type childrenLoader struct {
	childType string
	ids       []string
}

func (l childrenLoader) Apply(config *proto.EntityRequest) {
	if config.View.Children == nil {
		config.View.Children = make([]*proto.ChildRequest, 0)
	}

	config.View.Children = append(config.View.Children, &proto.ChildRequest{
		Type: &proto.Key{Key: l.childType},
		Cid:  l.ids,
	})
}

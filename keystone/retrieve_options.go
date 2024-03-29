package keystone

import (
	"github.com/kubex/keystone-go/proto"
)

// RetrieveBy is an interface that defines a retriever
type RetrieveBy interface {
	BaseRequest() *proto.EntityRequest
}

// byEntityID is a retriever that retrieves an entity by its ID
type byEntityID struct {
	EntityID string
}

func ByEntityID(entityID string) RetrieveBy {
	return byEntityID{EntityID: entityID}
}

// BaseRequest returns the base byEntityID request
func (l byEntityID) BaseRequest() *proto.EntityRequest {
	return &proto.EntityRequest{
		View:     &proto.EntityView{},
		EntityId: l.EntityID,
	}
}

// byUniqueProperty is a retriever that retrieves an entity by its unique ID
type byUniqueProperty struct {
	UniqueID string
	Property string
}

func ByUniqueProperty(uniqueID, property string) RetrieveBy {
	return byUniqueProperty{UniqueID: uniqueID, Property: property}
}

// BaseRequest returns the base byUniqueProperty request
func (l byUniqueProperty) BaseRequest() *proto.EntityRequest {
	return &proto.EntityRequest{
		View: &proto.EntityView{},
		UniqueId: &proto.IDLookup{
			SchemaId: "SCHEMAID", //TODO: THIS!
			Property: l.Property,
			UniqueId: l.UniqueID,
		},
	}
}

// RetrieveOption is an interface for options to be applied to an entity request
type RetrieveOption interface {
	Apply(config *proto.EntityView)
}

type retrieveOptions []RetrieveOption

func (o retrieveOptions) Apply(config *proto.EntityView) {
	for _, opt := range o {
		opt.Apply(config)
	}
}

func RetrieveOptions(opts ...RetrieveOption) RetrieveOption {
	return retrieveOptions(opts)
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

func (l propertyLoader) Apply(config *proto.EntityView) {
	if config.Properties == nil {
		config.Properties = make([]*proto.PropertyRequest, 0)
	}

	config.Properties = append(config.Properties, &proto.PropertyRequest{Properties: l.properties, Decrypt: l.decrypt})
}

type linksLoader struct{ Links []string }

func (l linksLoader) Apply(config *proto.EntityView) {
	if config.LinkByType == nil {
		config.LinkByType = make([]*proto.Key, 0)
	}
	for _, link := range l.Links {
		config.LinkByType = append(config.LinkByType, &proto.Key{Key: link})
	}
}

type relationshipsLoader struct{ keys []string }

func (l relationshipsLoader) Apply(config *proto.EntityView) {
	if config.RelationshipByType == nil {
		config.RelationshipByType = make([]*proto.Key, 0)
	}
	for _, key := range l.keys {
		config.RelationshipByType = append(config.RelationshipByType, &proto.Key{Key: key})
	}
}

type childrenLoader struct {
	childType string
	ids       []string
}

func (l childrenLoader) Apply(config *proto.EntityView) {
	if config.Children == nil {
		config.Children = make([]*proto.ChildRequest, 0)
	}

	config.Children = append(config.Children, &proto.ChildRequest{
		Type: &proto.Key{Key: l.childType},
		Cid:  l.ids,
	})
}

type viewName struct{ name string }

func (l viewName) Apply(config *proto.EntityView) { config.Name = l.name }
func WithView(name string) RetrieveOption {
	return viewName{name: name}
}

type summaryLoader struct{ summary bool }

func (l summaryLoader) Apply(config *proto.EntityView) { config.Summary = l.summary }

// WithSummary is a retrieve option that loads summaries
func WithSummary() RetrieveOption {
	return summaryLoader{summary: true}
}

type datumLoader struct{ datum bool }

func (l datumLoader) Apply(config *proto.EntityView) { config.Datum = l.datum }

// WithDatum is a retrieve option that loads datum
func WithDatum() RetrieveOption {
	return datumLoader{datum: true}
}

type labelLoader struct{ labels bool }

func (l labelLoader) Apply(config *proto.EntityView) { config.Labels = l.labels }

type relationshipCount struct{ count bool }

func (l relationshipCount) Apply(config *proto.EntityView) { config.RelationshipCount = l.count }
func WithRelationshipCount() RetrieveOption                { return relationshipCount{count: true} }

type childCount struct{ count bool }

func (l childCount) Apply(config *proto.EntityView) { config.ChildCount = l.count }
func WithChildCount() RetrieveOption                { return relationshipCount{count: true} }

type linkCount struct{ count bool }

func (l linkCount) Apply(config *proto.EntityView) { config.LinkCount = l.count }
func WithLinkCount() RetrieveOption                { return relationshipCount{count: true} }

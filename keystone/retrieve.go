package keystone

import (
	"context"
	"errors"

	"github.com/kubex/keystone-go/proto"
)

// Actor is a struct that represents an actor
type Actor struct {
	connection   *Connection
	workspaceID  string
	mutator      *proto.Mutator
	loadedEntity *proto.EntityResponse
}

func (a *Actor) VendorID() string {
	return a.connection.appID.GetVendorId()
}

func (a *Actor) AppID() string {
	return a.connection.appID.GetAppId()
}

func (a *Actor) VendorApp() *proto.VendorApp {
	return &a.connection.appID
}

func (a *Actor) WorkspaceID() string {
	return a.workspaceID
}

func (a *Actor) authorization() *proto.Authorization {
	return &proto.Authorization{
		Source:      &a.connection.appID,
		Token:       a.connection.token,
		WorkspaceId: a.workspaceID,
	}
}

// SetClient sets the client name for the actor
func (a *Actor) SetClient(client string) {
	if a.mutator != nil {
		a.mutator.Client = client
	}
}

// Get retrieves an entity by the given retrieveBy, storing the result in dst
func (a *Actor) Get(ctx context.Context, retrieveBy Retriever, dst interface{}, opts ...RetrieveOption) error {
	entityRequest := retrieveBy.BaseRequest()
	entityRequest.Authorization = a.authorization()
	for _, opt := range opts {
		opt.Apply(entityRequest)
	}

	_, loadByUnique := retrieveBy.(RetrieveByUnique)
	_, genericResult := dst.(GenericResult)
	if loadByUnique && genericResult {
		return errors.New("invalid retrieveBy and dst combination")
	}

	view := entityRequest.View

	// set source
	for _, p := range view.Properties {
		p.Source = a.authorization().GetSource()
	}
	for _, l := range view.LinkByType {
		l.Source = a.authorization().GetSource()
	}
	for _, r := range view.RelationshipByType {
		r.Source = a.authorization().GetSource()
	}

	if _, ok := retrieveBy.(RetrieveByUnique); ok {
		schema, registered := a.connection.registerType(dst)
		if !registered {
			// wait for the type to be registered with the keystone server
			a.connection.SyncSchema().Wait()
		}

		schemaID := schema.Id
		if schemaID == "" {
			schemaID = schema.Type
		}
		entityRequest.UniqueId.SchemaId = schemaID
	}

	resp, err := a.connection.Retrieve(ctx, entityRequest)
	if err != nil {
		return err
	}
	a.loadedEntity = resp
	if gr, ok := dst.(GenericResult); ok {
		return UnmarshalGeneric(resp, gr)
	}

	return Unmarshal(resp, dst)
}

// Find returns a list of entities matching the given entityType and retrieveProperties
func (a *Actor) Find(ctx context.Context, entityType string, retrieveProperties []string, options ...FindOption) ([]*proto.EntityResponse, error) {
	findRequest := &proto.FindRequest{
		Authorization: a.authorization(),
		Schema:        &proto.Key{Key: entityType, Source: a.authorization().Source},
	}

	fReq := &filterRequest{Properties: []*proto.PropertyRequest{{Properties: retrieveProperties}}}

	for _, opt := range options {
		opt.Apply(fReq)
	}

	findRequest.Properties = fReq.Properties
	findRequest.Filters = fReq.Filters
	findRequest.Labels = fReq.Labels
	findRequest.RelationOf = fReq.RelationOf

	resp, err := a.connection.Find(ctx, findRequest)
	if err != nil {
		return nil, err
	}
	return resp.Entities, nil
}

// List returns a list of entities within an active set
func (a *Actor) List(ctx context.Context, entityType, activeSetName string, retrieveProperties []string, options ...FindOption) ([]*proto.EntityResponse, error) {
	listRequest := &proto.ADSListRequest{
		Authorization: a.authorization(),
		Schema:        &proto.Key{Key: entityType, Source: a.authorization().Source},
		AdsName:       activeSetName,
		Properties:    retrieveProperties,
	}

	fReq := &filterRequest{}
	for _, opt := range options {
		opt.Apply(fReq)
	}

	listRequest.Filters = fReq.Filters
	listRequest.Limit = fReq.Limit
	listRequest.Offset = fReq.Offset
	listRequest.SortProperty = fReq.SortProperty
	listRequest.SortDirection = fReq.SortDirection

	resp, err := a.connection.ADSList(ctx, listRequest)
	if err != nil {
		return nil, err
	}
	return resp.Entities, nil
}

type filterRequest struct {
	Properties    []*proto.PropertyRequest
	Filters       []*proto.PropertyFilter
	Labels        []*proto.EntityLabel
	RelationOf    *proto.RelationOf
	Limit         int32
	Offset        int32
	SortProperty  string
	SortDirection string
}

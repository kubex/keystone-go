package keystone

import (
	"context"
	"errors"

	"github.com/kubex/keystone-go/proto"
)

// Actor is a struct that represents an actor
type Actor struct {
	connection  *Connection
	workspaceID string
	mutator     *proto.Mutator
}

func (a *Actor) ReplaceConnection(c *Connection) { a.connection = c }

func (a *Actor) UserAgent() string { return a.mutator.GetUserAgent() }
func (a *Actor) RemoteIp() string  { return a.mutator.GetRemoteIp() }
func (a *Actor) UserId() string    { return a.mutator.GetUserId() }
func (a *Actor) Client() string    { return a.mutator.GetClient() }

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
	if a == nil || a.connection == nil {
		return nil
	}
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
func (a *Actor) Get(ctx context.Context, retrieveBy RetrieveBy, dst interface{}, retrieve RetrieveOption) error {
	entityRequest := retrieveBy.BaseRequest()
	entityRequest.Authorization = a.authorization()
	if retrieve != nil {
		retrieve.Apply(entityRequest.View)
	}

	_, loadByUnique := retrieveBy.(byUniqueProperty)
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

	schema, registered := a.connection.registerType(dst)
	if !registered {
		// wait for the type to be registered with the keystone server
		a.connection.SyncSchema().Wait()
	}
	entityRequest.Schema = &proto.Key{Key: schema.GetType(), Source: a.authorization().Source}

	if _, ok := retrieveBy.(byUniqueProperty); ok {
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
	if be, ok := dst.(BaseEntity); ok {
		be._lastLoad = resp
	}
	if gr, ok := dst.(GenericResult); ok {
		return UnmarshalGeneric(resp, gr)
	}

	return Unmarshal(resp, dst)
}

// Find returns a list of entities matching the given entityType and retrieveProperties
func (a *Actor) Find(ctx context.Context, entityType string, retrieve RetrieveOption, options ...FindOption) ([]*proto.EntityResponse, error) {
	findRequest := &proto.FindRequest{
		Authorization: a.authorization(),
		Schema:        &proto.Key{Key: entityType, Source: a.authorization().Source},
		View:          &proto.EntityView{},
	}

	if retrieve != nil {
		retrieve.Apply(findRequest.View)
	}

	fReq := &filterRequest{Properties: []*proto.PropertyRequest{}}

	for _, opt := range options {
		opt.Apply(fReq)
	}

	findRequest.PropertyFilters = fReq.Filters
	findRequest.LabelFilters = fReq.Labels
	findRequest.RelationOf = fReq.RelationOf
	findRequest.ParentEntityId = fReq.ParentEntityID

	resp, err := a.connection.Find(ctx, findRequest)
	if err != nil {
		return nil, err
	}
	return resp.Entities, nil
}

// List returns a list of entities within an active set
func (a *Actor) List(ctx context.Context, entityType, activeSetName string, retrieveProperties []string, options ...FindOption) ([]*proto.EntityResponse, error) {
	listRequest := &proto.ActiveSetListRequest{
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

	resp, err := a.connection.ActiveSetList(ctx, listRequest)
	if err != nil {
		return nil, err
	}
	return resp.Entities, nil
}

type filterRequest struct {
	Properties     []*proto.PropertyRequest
	Filters        []*proto.PropertyFilter
	Labels         []*proto.EntityLabel
	RelationOf     *proto.RelationOf
	ParentEntityID string
	Limit          int32
	Offset         int32
	SortProperty   string
	SortDirection  string
}

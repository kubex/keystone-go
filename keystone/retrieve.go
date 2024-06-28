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
	user        *proto.User
}

func (a *Actor) ReplaceConnection(c *Connection) { a.connection = c }

func (a *Actor) UserAgent() string { return a.user.GetUserAgent() }
func (a *Actor) RemoteIp() string  { return a.user.GetRemoteIp() }
func (a *Actor) UserId() string    { return a.user.GetUserId() }
func (a *Actor) Client() string    { return a.user.GetClient() }

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

func (a *Actor) Authorization() *proto.Authorization {
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
	if a.user != nil {
		a.user.Client = client
	}
}

func (a *Actor) GetByID(ctx context.Context, entityID string, dst interface{}, retrieve ...RetrieveOption) error {
	return a.Get(ctx, ByEntityID(Type(dst), entityID), dst, retrieve...)
}

func (a *Actor) GetByUniqueProperty(ctx context.Context, uniqueId, propertyName string, dst interface{}, retrieve ...RetrieveOption) error {
	return a.Get(ctx, ByUniqueProperty(Type(dst), uniqueId, propertyName), dst, retrieve...)
}

// Get retrieves an entity by the given retrieveBy, storing the result in dst
func (a *Actor) Get(ctx context.Context, retrieveBy RetrieveBy, dst interface{}, retrieve ...RetrieveOption) error {
	entityRequest := retrieveBy.BaseRequest()
	entityRequest.Authorization = a.Authorization()
	for _, rOpt := range retrieve {
		rOpt.Apply(entityRequest.View)
	}

	_, loadByUnique := retrieveBy.(byUniqueProperty)
	_, genericResult := dst.(GenericResult)
	if loadByUnique && genericResult {
		return errors.New("invalid retrieveBy and dst combination")
	}

	view := entityRequest.View

	// set source
	for _, p := range view.Properties {
		p.Source = a.Authorization().GetSource()
	}

	for _, r := range view.RelationshipByType {
		r.Source = a.Authorization().GetSource()
	}

	schema, registered := a.connection.registerType(dst)
	if !registered {
		// wait for the type to be registered with the keystone server
		a.connection.SyncSchema().Wait()
	}
	entityRequest.Schema = &proto.Key{Key: schema.GetType(), Source: a.Authorization().Source}

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
		Authorization: a.Authorization(),
		Schema:        &proto.Key{Key: entityType, Source: a.Authorization().Source},
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
	listRequest := &proto.ListRequest{
		Authorization: a.Authorization(),
		Schema:        &proto.Key{Key: entityType, Source: a.Authorization().Source},
		FromView:      activeSetName,
		Properties:    retrieveProperties,
	}

	fReq := &filterRequest{}
	for _, opt := range options {
		opt.Apply(fReq)
	}

	listRequest.Filters = fReq.Filters
	listRequest.Page = &proto.PageRequest{
		PerPage:    fReq.PerPage,
		PageNumber: fReq.PageNumber,
	}

	listRequest.Sort = []*proto.PropertySort{{
		Property:   fReq.SortProperty,
		Descending: fReq.SortDescending,
	},
	}

	resp, err := a.connection.List(ctx, listRequest)
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
	PerPage        int32
	PageNumber     int32
	SortProperty   string
	SortDescending bool
}

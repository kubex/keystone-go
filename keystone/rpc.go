package keystone

import (
	"context"
	"github.com/kubex/definitions-go/app"
	"github.com/kubex/definitions-go/k4"
	"github.com/kubex/keystone-go/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

type Connection struct {
	client proto.KeystoneClient
	actor  Actor
	appID  app.GlobalAppID
	token  string
}

func (c *Connection) ProtoClient() proto.KeystoneClient {
	return c.client
}

// IsAppSchema checks if the schema is for the current app
func (c *Connection) IsAppSchema(schema *proto.Key) bool {
	return c.appID.VendorID == schema.GetVendorId() && c.appID.AppID == schema.GetAppId()
}

func (c *Connection) authorization() *proto.Authorization {
	return &proto.Authorization{
		VendorId: c.appID.VendorID,
		AppId:    c.appID.AppID,
		Token:    c.token,
	}
}

func (c *Connection) MakeKey(name string) *proto.Key {
	return &proto.Key{
		VendorId: c.appID.VendorID,
		AppId:    c.appID.AppID,
		Key:      name,
	}
}

type FilterOperator string

const FilterEqual FilterOperator = "eq"
const FilterNotEqual FilterOperator = "neq"
const FilterGreaterThan FilterOperator = "gt"
const FilterGreaterThanEqual FilterOperator = "gte"
const FilterLessThan FilterOperator = "lt"
const FilterLessThanEqual FilterOperator = "lte"
const FilterLike FilterOperator = "like"
const FilterIn FilterOperator = "in"

func (c *Connection) NewFilter(property string, operator FilterOperator, value *proto.Value) *proto.FilterProperty {
	return &proto.FilterProperty{
		Property: c.MakeKey(property),
		Operator: string(operator),
		Value:    value,
	}
}

type Actor struct {
	UserAgent string
	RemoteIP  string
	UserID    string
	Client    string
}

func NewActor(remoteIP, userID, userAgent string) Actor {
	return Actor{
		UserAgent: userAgent,
		RemoteIP:  remoteIP,
		UserID:    userID,
		Client:    "keystone-go",
	}
}

func NewConnection(client proto.KeystoneClient, appID app.GlobalAppID, accessToken string, actor Actor) *Connection {
	return &Connection{client: client, actor: actor, appID: appID, token: accessToken}
}

func (c *Connection) Apply(ctx context.Context, entity *Entity) (*proto.MutateResponse, error) {
	applyMutation := &proto.Mutation{
		Mutator: &proto.Mutator{
			UserAgent: c.actor.UserAgent,
			RemoteIp:  c.actor.RemoteIP,
			UserId:    c.actor.UserID,
			Client:    c.actor.Client,
		},
	}

	for _, prop := range entity.Properties {

		if !prop.updated {
			continue
		}
		prop.updated = true
		applyMutation.Properties = append(applyMutation.Properties, prop.toProto())
	}

	for _, prop := range entity.DeleteProperties {
		applyMutation.PropertyDeletes = append(applyMutation.PropertyDeletes, prop.Name)
	}

	for _, logEntry := range entity.LogEntries {
		if logEntry.Actor == nil {
			logEntry.Actor = &c.actor
		}

		applyMutation.Logs = append(applyMutation.Logs, &proto.Log{
			Level:     logEntry.Level.toProto(),
			Message:   logEntry.Message,
			Reference: logEntry.Reference,
			TraceId:   logEntry.TraceID,
			Time:      timestamppb.New(logEntry.Time),
		})
		logEntry.written = true
	}

	for _, event := range entity.Events {
		applyMutation.Events = append(applyMutation.Events, &proto.Event{
			Type: event.Type,
			Time: timestamppb.New(event.Time),
			Data: event.Data,
		})
		event.written = true
	}

	for _, child := range entity.Children {
		protoChild := &proto.Child{
			Type: child.Type,
			Data: child.Data,
		}
		if child.ID != "" {
			protoChild.Id = child.ID
		}
		applyMutation.Children = append(applyMutation.Children, protoChild)
		child.written = true
	}

	for _, rel := range entity.Relationships {
		if !rel.written {
			protoRel := &proto.Relationship{
				Relationship:  rel.Relationship,
				DestinationId: rel.DestinationID,
				Since:         timestamppb.New(rel.Since),
				Data:          rel.Data,
				BiDirectional: false,
			}
			applyMutation.Relationships = append(applyMutation.Relationships, protoRel)
			rel.written = true
			log.Println("writing relationship")
		} else {
			log.Println("skipping relationship write")
		}
	}

	mutate := &proto.MutateRequest{
		Authorization: c.authorization(),
		WorkspaceId:   entity.WorkspaceID,
		EntityId:      entity.ID.String(),
		Schema:        toKey(entity.Schema),
		Mutations:     []*proto.Mutation{applyMutation},
	}

	if entity.ID.String() == "" {
		mutate.EntityId = entity.ID.String()
	}

	mutateResp, err := c.client.Mutate(ctx, mutate)
	if err != nil {
		return nil, err
	}

	if mutateResp.Success {
		entity.ID = k4.IDFromString(mutateResp.EntityId)
	}

	//TODO: on success, update property updates to true, instead of on read

	return mutateResp, err
}

func (c *Connection) Retrieve(ctx context.Context, workspaceID, entityId string, retrieveProperties []string) (*proto.EntityResponse, error) {
	retrieveReq := &proto.RetrieveRequest{
		WorkspaceId: workspaceID,
		EntityId:    entityId,
		Properties:  []*proto.PropertyRequest{},
	}

	for _, prop := range retrieveProperties {
		retrieveReq.Properties = append(retrieveReq.Properties, &proto.PropertyRequest{
			Property: c.MakeKey(prop),
		})
	}

	return c.client.Retrieve(ctx, retrieveReq)
}

func (c *Connection) Lookup(ctx context.Context, workspaceID, idLookup string) ([]*proto.EntityResponse, error) {
	located, err := c.client.Lookup(ctx, &proto.LookupRequest{
		WorkspaceId: workspaceID,
		LookupId:    idLookup,
	})

	return located.GetEntities(), err
}

func (c *Connection) Find(ctx context.Context, workspaceID, entityType string, retrieveProperties []string, options ...Option) ([]*proto.EntityResponse, error) {

	findReq := &proto.FindRequest{
		WorkspaceId: workspaceID,
		Schema:      c.MakeKey(entityType),
		Properties:  []*proto.PropertyRequest{},
	}

	for _, option := range options {
		for _, filter := range option.Filters {
			findReq.Filters = append(findReq.Filters, filter)
		}
	}

	for _, prop := range retrieveProperties {
		findReq.Properties = append(findReq.Properties, &proto.PropertyRequest{
			Property: c.MakeKey(prop),
		})
	}

	located, err := c.client.Find(ctx, findReq)
	return located.GetEntities(), err
}

type Option struct {
	Filters []*proto.FilterProperty
}

func NewFilterOption(filters ...*proto.FilterProperty) Option {
	return Option{
		Filters: filters,
	}
}

func TextValue(value string) *proto.Value {
	return &proto.Value{
		Type: proto.ValueType_VALUE_TEXT,
		Text: value,
	}
}

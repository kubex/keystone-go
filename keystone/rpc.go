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

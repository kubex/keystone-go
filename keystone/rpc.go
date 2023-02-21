package keystone

import (
	"context"
	"github.com/kubex/definitions-go/k4"
	"github.com/kubex/keystone-go/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Connection struct {
	client proto.KeystoneClient
	actor  Actor
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

func NewConnection(client proto.KeystoneClient, actor Actor) *Connection {
	return &Connection{client: client, actor: actor}
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

		protoProp := &proto.Property{
			Key: toKey(prop.Name),
		}

		switch prop.Type {
		case PropertyTypeText:
			protoProp.Text = prop.Text
		case PropertyTypeInt:
			protoProp.Int = prop.Int
		case PropertyTypeBool:
			protoProp.Bool = prop.Bool
		case PropertyTypeFloat:
			protoProp.Float = float32(prop.Float)
		case PropertyTypeTime:
			protoProp.Time = timestamppb.New(*prop.Time)
		}

		protoProp.Type = prop.Type.toProto()

		protoProp.Classification = prop.Classification.toProto()

		if prop.Secret != "" {
			protoProp.SecureText = prop.Secret
		}

		applyMutation.Properties = append(applyMutation.Properties, protoProp)
	}

	for _, prop := range entity.DeleteProperties {
		applyMutation.PropertyDeletes = append(applyMutation.PropertyDeletes, toKey(prop.Name))
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
		var props []*proto.Property
		for _, prop := range event.Data {
			props = append(props, prop.toProto())
		}

		applyMutation.Events = append(applyMutation.Events, &proto.Event{
			Type: toKey(event.Type),
			Time: timestamppb.New(event.Time),
			Data: props,
		})
		event.written = true
	}

	for _, child := range entity.Children {
		protoChild := &proto.Child{
			Type: toKey(child.Type),
			Data: child.Data,
		}
		if child.ID != "" {
			protoChild.Id = child.ID
		}
		applyMutation.Children = append(applyMutation.Children, protoChild)
		child.written = true
	}

	mutate := &proto.MutateRequest{
		WorkspaceId: entity.WorkspaceID,
		Schema:      toKey(entity.Schema),
		Mutations:   []*proto.Mutation{applyMutation},
	}

	if entity.ID.Full() == "" {
		mutate.EntityId = entity.ID.Full()
	}

	mutateResp, err := c.client.Mutate(ctx, mutate)

	if mutateResp.Success {
		entity.ID = k4.IDFromString(mutateResp.EntityId)
	}

	//TODO: on success, update property updates to true, instead of on read

	return mutateResp, err
}

package keystone

import (
	"context"
	"github.com/kubex/keystone-go/proto"
	"log"
)

func (a *Actor) GetByID(ctx context.Context, entityID string, dst interface{}) error {

	resp, err := a.connection.ProtoClient().Retrieve(ctx, &proto.RetrieveRequest{
		Authorization: a.authorization(),
		EntityId:      entityID,
		Properties: []*proto.PropertyRequest{
			{Keys: []string{"address~"}},
		},
	})
	if err != nil {
		return err
	}

	return entityResponseToDst(resp, dst)
}

func (a *Actor) GetByUnique(ctx context.Context, key, value string, dst interface{}) error {

	schema, _ := a.connection.registerType(dst)
	schemaID := schema.Id
	if schemaID == "" {
		schemaID = schema.Type
	}

	resp, err := a.connection.ProtoClient().Retrieve(ctx, &proto.RetrieveRequest{
		Authorization: a.authorization(),
		UniqueId: &proto.IDLookup{
			Field:    key,
			UniqueId: value,
			SchemaId: schemaID,
		},
		Properties: []*proto.PropertyRequest{
			{Keys: []string{"address~"}},
			{Keys: []string{"name", "email"}},
		},
	})

	if err != nil {
		return err
	}

	return entityResponseToDst(resp, dst)
}

func entityResponseToDst(resp *proto.EntityResponse, dst interface{}) error {

	log.Println(resp.GetProperties())

	return nil
}

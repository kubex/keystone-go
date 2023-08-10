package keystone

import (
	"context"
	"errors"

	"github.com/kubex/keystone-go/proto"
)

type GenericResult map[string]interface{}

type Actor struct {
	connection   *Connection
	workspaceID  string
	mutator      *proto.Mutator
	loadedEntity *proto.EntityResponse
}

func (a *Actor) authorization() *proto.Authorization {
	return &proto.Authorization{
		Source:      &a.connection.appID,
		Token:       a.connection.token,
		WorkspaceId: a.workspaceID,
	}
}

func (a *Actor) SetClient(client string) {
	a.mutator.Client = client
}

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

	resp, err := a.connection.ProtoClient().Retrieve(ctx, entityRequest)
	if err != nil {
		return err
	}
	a.loadedEntity = resp
	if gr, ok := dst.(GenericResult); ok {
		return UnmarshalGeneric(resp, gr)
	}

	return Unmarshal(resp, dst)
}

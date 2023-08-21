package keystone

import (
	"context"
	"errors"
	"reflect"

	"github.com/kubex/keystone-go/proto"
)

// Mutate is a function that can mutate an entity
func (a *Actor) Mutate(ctx context.Context, src interface{}, comment string) error {
	if reflect.TypeOf(src).Kind() != reflect.Pointer {
		return errors.New("mutate requires a pointer to a struct")
	}

	//log.Println("Processing Mutate request")
	schema, registered := a.connection.registerType(src)
	if !registered {
		// wait for the type to be registered with the keystone server
		a.connection.SyncSchema().Wait()
	}
	//log.Println("Marshalling entity", src)

	encoder := &PropertyEncoder{}
	mutation := encoder.Marshal(src)
	mutation.Mutator = a.mutator
	entityID := encoder.EntityID
	mutation.Comment = comment
	if rawEntity, ok := src.(Entity); ok && entityID == "" {
		entityID = rawEntity.GetKeystoneID()
	}

	if entityWithLabels, ok := src.(EntityLabelProvider); ok {
		mutation.Labels = entityWithLabels.GetKeystoneLabels()
	}

	if entityWithLinks, ok := src.(EntityLinkProvider); ok {
		mutation.Links = entityWithLinks.GetKeystoneLinks()
	}

	if entityWithRelationships, ok := src.(EntityRelationshipProvider); ok {
		mutation.Relationships = entityWithRelationships.GetKeystoneRelationships()
	}

	if entityWithEvents, ok := src.(EntityEventProvider); ok {
		mutation.Events = entityWithEvents.GetKeystoneEvents()
	}

	if a.loadedEntity != nil {
		mutation.Properties = a.getChangedProperties(a.loadedEntity, &proto.EntityResponse{Properties: mutation.Properties})
	}
	m := &proto.MutateRequest{
		Authorization: &proto.Authorization{WorkspaceId: a.workspaceID, Source: &a.connection.appID},
		EntityId:      entityID,
		Schema:        &proto.Key{Key: schema.Type, Source: schema.Source}, // TODO: Should probably provide the schema ID if we have it - and verify against the type / source
		Mutation:      mutation,
	}

	mResp, err := a.connection.Mutate(ctx, m)

	if err == nil && mResp.Success {
		if rawEntity, ok := src.(Entity); ok && entityID == "" {
			rawEntity.SetKeystoneID(mResp.GetEntityId())
		}
	}

	return err
}

func (a *Actor) getChangedProperties(existing, newValues *proto.EntityResponse) []*proto.EntityProperty {
	exMap := makeEntityPropertyMap(existing)
	newMap := makeEntityPropertyMap(newValues)

	var result []*proto.EntityProperty
	for k, v := range newMap {
		if _, ok := exMap[k]; !ok {
			result = append(result, v)
			continue
		}
		if newMap[k].Property.Key == exMap[k].Property.Key &&
			newMap[k].Value.Text == exMap[k].Value.Text &&
			newMap[k].Value.SecureText == exMap[k].Value.SecureText &&
			newMap[k].Value.Int == exMap[k].Value.Int &&
			newMap[k].Value.Float == exMap[k].Value.Float &&
			newMap[k].Value.Bool == exMap[k].Value.Bool &&
			reflect.DeepEqual(newMap[k].Value.Map, exMap[k].Value.Map) &&
			reflect.DeepEqual(newMap[k].Value.Set, exMap[k].Value.Set) &&
			reflect.DeepEqual(newMap[k].Value.Time, exMap[k].Value.Time) {
			continue
		}
	}
	return result
}

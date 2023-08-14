package keystone

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/kubex/keystone-go/proto"
)

func (a *Actor) Marshal(src interface{}, comment string) error {
	if reflect.TypeOf(src).Kind() != reflect.Pointer {
		return errors.New("marshal requires a pointer to a struct")
	}

	//log.Println("Processing Marshal request")
	schema, registered := a.connection.registerType(src)
	if !registered {
		// wait for the type to be registered with the keystone server
		a.connection.SyncSchema().Wait()
	}
	//log.Println("Marshalling entity", src)

	mutation := &proto.Mutation{
		Mutator:    a.mutator,
		Comment:    comment,
		Properties: []*proto.EntityProperty{},
		Children:   []*proto.EntityChild{},
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

	extractor := &PropertyExtractor{}
	if err := extractor.Extract(src); err != nil {
		fmt.Println("Error extracting properties", err)
		return err
	}

	mutation.Properties = extractor.Properties()
	mutation.Children = extractor.Children()
	if a.loadedEntity != nil {
		mutation.Properties = a.getChangedProperties(a.loadedEntity, mutation.Properties)
	}
	m := &proto.MutateRequest{
		Authorization: &proto.Authorization{WorkspaceId: a.workspaceID, Source: &a.connection.appID},
		EntityId:      extractor.EntityID,
		Schema:        &proto.Key{Key: schema.Type, Source: schema.Source}, // TODO: Should probably provide the schema ID if we have it - and verify against the type / source
		Mutation:      mutation,
	}

	_, err := a.connection.ProtoClient().Mutate(context.Background(), m)
	return err
}

func (a *Actor) getChangedProperties(existing *proto.EntityResponse, newValues []*proto.EntityProperty) []*proto.EntityProperty {
	exMap := makeEntityPropertyMap(existing.GetProperties())
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

func supportedType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.String, reflect.Int32, reflect.Int64, reflect.Int, reflect.Bool, reflect.Float32, reflect.Float64, reflect.Map:
		return true
	}

	switch t {
	case typeOfSecretString, typeOfAmount, typeOfTime, typeOfStringSlice:
		return true
	}

	return false
}

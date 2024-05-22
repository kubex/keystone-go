package keystone

import (
	"context"
	"github.com/kubex/keystone-go/proto"
	"reflect"
)

func (a *Actor) SetDynamicProperties(ctx context.Context, entityID string, setProperties []*proto.EntityProperty, removeProperties []string, comment string) error {

	for _, prop := range setProperties {
		prop.Property.Source = a.VendorApp()
	}

	mutation := &proto.Mutation{
		DynamicProperties:       setProperties,
		RemoveDynamicProperties: removeProperties,
		Comment:                 comment,
	}
	mutation.Mutator = a.mutator

	m := &proto.MutateRequest{
		Authorization: &proto.Authorization{WorkspaceId: a.workspaceID, Source: &a.connection.appID},
		EntityId:      entityID,
		Mutation:      mutation,
	}

	return mutateToError(a.connection.Mutate(ctx, m))
}

func NewProperties(props map[string]interface{}) []*proto.EntityProperty {
	properties := make([]*proto.EntityProperty, 0, len(props))
	for key, value := range props {
		prop := NewProperty(key, value)
		if prop != nil {
			properties = append(properties, prop)
		}
	}
	return properties
}

func NewProperty(key string, value interface{}) *proto.EntityProperty {
	v := reflect.ValueOf(value)
	prop, isEmpty := entityPropertyFromField(v, v.Type(), fieldOptions{name: key})
	if isEmpty {
		return nil
	}
	return prop
}

func NewAppendProperty(key string, value interface{}) *proto.EntityProperty {
	prop := NewProperty(key, value)
	prop.ModifyType = proto.EntityProperty_Append
	return prop
}

func NewReduceProperty(key string, value interface{}) *proto.EntityProperty {
	prop := NewProperty(key, value)
	prop.ModifyType = proto.EntityProperty_Reduce
	return prop
}

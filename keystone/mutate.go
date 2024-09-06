package keystone

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/kubex/keystone-go/proto"
)

func (a *Actor) RemoteMutate(ctx context.Context, src interface{}, comment string) error {
	mutation := &proto.Mutation{}
	entityID := ""
	if rawEntity, ok := src.(Entity); ok {
		entityID = rawEntity.GetKeystoneID()
	}

	if entityID == "" {
		return errors.New("entityID is required for remote mutations")
	}

	if entityWithSensor, ok := src.(EntitySensorProvider); ok {
		mutation.Measurements = entityWithSensor.GetKeystoneSensorMeasurements()
	}
	if entityWithEvents, ok := src.(EntityEventProvider); ok {
		mutation.Events = entityWithEvents.GetKeystoneEvents()
	}
	if entityWithLogs, ok := src.(EntityLogProvider); ok {
		mutation.Logs = entityWithLogs.GetKeystoneLogs()
	}

	m := &proto.MutateRequest{
		Authorization: a.Authorization(),
		EntityId:      entityID,
		Mutation:      mutation,
	}

	return mutateToError(a.connection.Mutate(ctx, m))
}

type MutateOption interface {
	apply(*proto.MutateRequest)
}

// OnConflictUseID should set the unique properties that can be used to identify an existing identity
func OnConflictUseID(property ...string) MutateOption {
	return onConflictUseID{Property: property}
}

type onConflictUseID struct {
	Property []string
}

func (m onConflictUseID) apply(mutate *proto.MutateRequest) {
	mutate.ConflictUniquePropertyAcquire = m.Property
}

// MutateProperties Only mutate the specified properties
func MutateProperties(property ...string) MutateOption {
	return mutateProperties{Property: property}
}

type mutateProperties struct {
	Property []string
}

func (m mutateProperties) apply(mutate *proto.MutateRequest) {
	var keepProps []*proto.EntityProperty

	for _, prop := range mutate.Mutation.Properties {
		for _, p := range m.Property {
			if prop.Property == p {
				keepProps = append(keepProps, prop)
				break
			}
		}
	}
	mutate.Mutation.Properties = keepProps
}

// Mutate is a function that can mutate an entity
func (a *Actor) Mutate(ctx context.Context, src interface{}, comment string, options ...MutateOption) error {
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
	mutation.Mutator = a.user
	entityID := ""
	mutation.Comment = comment
	if rawEntity, ok := src.(Entity); ok {
		entityID = rawEntity.GetKeystoneID()
	}

	if entityWithLabels, ok := src.(EntityLabelProvider); ok {
		mutation.Labels = entityWithLabels.GetKeystoneLabels()
	}

	if entityWithSensor, ok := src.(EntitySensorProvider); ok {
		mutation.Measurements = entityWithSensor.GetKeystoneSensorMeasurements()
	}

	if entityWithRelationships, ok := src.(EntityRelationshipProvider); ok {
		mutation.Relationships = entityWithRelationships.GetKeystoneRelationships()
	}

	if entityWithEvents, ok := src.(EntityEventProvider); ok {
		mutation.Events = entityWithEvents.GetKeystoneEvents()
	}

	if entityWithLogs, ok := src.(EntityLogProvider); ok {
		mutation.Logs = entityWithLogs.GetKeystoneLogs()
	}

	if base, ok := src.(BaseEntity); ok && base._lastLoad != nil {
		mutation.Properties = a.getChangedProperties(base._lastLoad, &proto.EntityResponse{Properties: mutation.Properties})
	} else if entityID != "" {
		mutation.Properties = a.getChangedProperties(nil, &proto.EntityResponse{Properties: mutation.Properties})
	}

	m := &proto.MutateRequest{
		Authorization: a.Authorization(),
		EntityId:      entityID,
		Schema:        &proto.Key{Key: schema.Type, Source: schema.Source}, // TODO: Should probably provide the schema ID if we have it - and verify against the type / source
		Mutation:      mutation,
	}

	for _, option := range options {
		option.apply(m)
	}

	mResp, err := a.connection.Mutate(ctx, m)

	if err == nil && mResp.Success {
		if rawEntity, ok := src.(Entity); ok && entityID == "" {
			rawEntity.SetKeystoneID(mResp.GetEntityId())
		}
	}

	return mutateToError(mResp, err)
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
		if newMap[k].Property == exMap[k].Property &&
			newMap[k].Value.Text == exMap[k].Value.Text &&
			newMap[k].Value.SecureText == exMap[k].Value.SecureText &&
			newMap[k].Value.Int == exMap[k].Value.Int &&
			newMap[k].Value.Float == exMap[k].Value.Float &&
			newMap[k].Value.Bool == exMap[k].Value.Bool &&
			reflect.DeepEqual(newMap[k].Value.Array.KeyValue, exMap[k].Value.Array.KeyValue) &&
			reflect.DeepEqual(newMap[k].Value.Array.Strings, exMap[k].Value.Array.Strings) &&
			reflect.DeepEqual(newMap[k].Value.Array.Ints, exMap[k].Value.Array.Ints) &&
			reflect.DeepEqual(newMap[k].Value.ArrayAppend.KeyValue, exMap[k].Value.ArrayAppend.KeyValue) &&
			reflect.DeepEqual(newMap[k].Value.ArrayAppend.Strings, exMap[k].Value.ArrayAppend.Strings) &&
			reflect.DeepEqual(newMap[k].Value.ArrayAppend.Ints, exMap[k].Value.ArrayAppend.Ints) &&
			reflect.DeepEqual(newMap[k].Value.ArrayReduce.KeyValue, exMap[k].Value.ArrayReduce.KeyValue) &&
			reflect.DeepEqual(newMap[k].Value.ArrayReduce.Strings, exMap[k].Value.ArrayReduce.Strings) &&
			reflect.DeepEqual(newMap[k].Value.ArrayReduce.Ints, exMap[k].Value.ArrayReduce.Ints) &&
			reflect.DeepEqual(newMap[k].Value.Time, exMap[k].Value.Time) {
			continue
		}
	}
	return result
}

func mutateToError(resp *proto.MutateResponse, err error) error {
	if err != nil {
		return err
	}

	if resp == nil {
		return errors.New("nil response")
	}

	if resp.ErrorCode > 0 || resp.ErrorMessage != "" {
		return fmt.Errorf("error %d: %s", resp.ErrorCode, resp.ErrorMessage)
	}
	return nil
}

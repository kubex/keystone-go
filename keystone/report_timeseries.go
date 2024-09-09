package keystone

import (
	"context"
	"errors"
	"github.com/kubex/keystone-go/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
)

// ReportTimeSeries writes point in time data
func (a *Actor) ReportTimeSeries(ctx context.Context, src interface{}) error {
	if reflect.TypeOf(src).Kind() != reflect.Pointer {
		return errors.New("mutate requires a pointer to a struct")
	}

	schema, registered := a.connection.registerType(src)
	if !registered {
		// wait for the type to be registered with the keystone server
		a.connection.SyncSchema().Wait()
	}

	var inputTime *timestamppb.Timestamp

	encoder := &PropertyEncoder{}
	mutation := encoder.Marshal(src)
	mutation.Mutator = a.user
	entityID := ""
	if rawEntity, ok := src.(Entity); ok {
		entityID = rawEntity.GetKeystoneID()
	}

	if tsEntity, ok := src.(TSEntity); ok {
		inputTime = timestamppb.New(tsEntity.GetTimeSeriesInputTime())
	} else {
		return errors.New("you must pass a TimeSeriesEntity as the source")
	}

	if entityWithLabels, ok := src.(EntityLabelProvider); ok {
		mutation.Labels = entityWithLabels.GetKeystoneLabels()
	}

	/*
		if entityWithSensor, ok := src.(EntitySensorProvider); ok {
			mutation.Measurements = entityWithSensor.GetKeystoneSensorMeasurements()
		}

		if entityWithRelationships, ok := src.(EntityRelationshipProvider); ok {
			mutation.Relationships = entityWithRelationships.GetKeystoneRelationships()
		}*/

	mutation.Properties = a.getChangedProperties(nil, &proto.EntityResponse{Properties: mutation.Properties})

	m := &proto.ReportTimeSeriesRequest{
		Authorization: a.Authorization(),
		EntityId:      entityID,
		Schema:        &proto.Key{Key: schema.Type, Source: schema.Source},
		Mutation:      mutation,
		Timestamp:     inputTime,
	}

	mResp, err := a.connection.ReportTimeSeries(ctx, m)

	if err == nil && mResp.Success {
		if rawEntity, ok := src.(Entity); ok && entityID == "" {
			rawEntity.SetKeystoneID(mResp.GetEntityId())
		}
	}

	return mutateToError(mResp, err)
}

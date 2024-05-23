package keystone

import (
	"github.com/kubex/keystone-go/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// EntitySensorProvider is an interface for entities that can provide Sensors
type EntitySensorProvider interface {
	ClearKeystoneSensorMeasurements() error
	GetKeystoneSensorMeasurements() []*proto.EntitySensorMeasurement
}

// EntitySensors is a struct that implements EntitySensorProvider
type EntitySensors struct {
	ksEntitySensorsMeasurements []*proto.EntitySensorMeasurement
}

// ClearKeystoneSensorMeasurements clears the Sensors
func (e *EntitySensors) ClearKeystoneSensorMeasurements() error {
	e.ksEntitySensorsMeasurements = []*proto.EntitySensorMeasurement{}
	return nil
}

// GetKeystoneSensorMeasurements returns the Sensors measurements
func (e *EntitySensors) GetKeystoneSensorMeasurements() []*proto.EntitySensorMeasurement {
	return e.ksEntitySensorsMeasurements
}

// AddKeystoneSensorMeasurement adds a Sensor measurement
func (e *EntitySensors) AddKeystoneSensorMeasurement(sensor string, value float32) {
	e.ksEntitySensorsMeasurements = append(e.ksEntitySensorsMeasurements, &proto.EntitySensorMeasurement{
		Sensor: sensor,
		Value:  value,
		At:     timestamppb.Now(),
	})
}

// AddKeystoneSensorMeasurementWithData adds a Sensor measurement
func (e *EntitySensors) AddKeystoneSensorMeasurementWithData(sensor string, value float32, data map[string]string) {
	e.ksEntitySensorsMeasurements = append(e.ksEntitySensorsMeasurements, &proto.EntitySensorMeasurement{
		Sensor: sensor,
		Value:  value,
		At:     timestamppb.Now(),
		Data:   data,
	})
}

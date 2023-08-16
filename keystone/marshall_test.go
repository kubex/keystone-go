package keystone

import (
	"reflect"
	"testing"
)

func TestMarshal(t *testing.T) {
	addr := &Address{
		Line1: "123 Fake St.",
		Line2: "Apt. 2",
		City:  "Springfield",
	}

	if err := getTestActor(&MockConnector{}).Marshal(addr, "test marshal"); err != nil {
		t.Error(err)
	}
}

func TestGetChangedProperties(t *testing.T) {
	diff := getTestActor(nil).getChangedProperties(testAddressResponse, testAddressResponseWithPostCode)
	if len(diff) != 1 {
		t.Errorf("expected 1 changed property, got %d", len(diff))
	}
	if diff[0].Property.Key != "post_code" {
		t.Errorf("expected post_code to be changed, got %s", diff[0])
	}
}

func TestSupportedType(t *testing.T) {
	types := []reflect.Type{
		reflect.TypeOf(""),
		reflect.TypeOf(0),
		reflect.TypeOf(0.0),
		reflect.TypeOf(true),
		typeOfAmount,
		typeOfSecretString,
		typeOfTime,
		typeOfStringSlice,
	}

	for _, typ := range types {
		if !supportedType(typ) {
			t.Errorf("type %v should be supported", typ)
		}
	}

	types = []reflect.Type{
		reflect.TypeOf(struct{}{}),
		reflect.TypeOf(reflect.PointerTo(reflect.TypeOf(""))),
	}

	for _, typ := range types {
		if supportedType(typ) {
			t.Errorf("type %v should not be supported", typ)
		}
	}
}

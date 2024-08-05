// Package keystone is a database abstraction layer
package keystone

import (
	"encoding/json"
	"reflect"
	"time"
)

// GenericResult is a map that can be used to retrieve a generic result
type GenericResult map[string]interface{}

// NestedChild is an interface that defines a child entity
type NestedChild interface {
	ChildID() string
	SetChildID(id string)
}

type NestedChildAggregateValue interface {
	AggregateValue() int64
	SetAggregateValue(val int64)
}

type NestedChildData interface {
	KeystoneData() map[string][]byte
	HydrateKeystoneData(data map[string][]byte)
}

type BaseNestedChild struct {
	_childID string
}

func getChildData(from any) map[string][]byte {
	childData := make(map[string][]byte)

	val := reflect.ValueOf(from)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		if !typ.Field(i).IsExported() {
			continue
		}
		fieldName := typ.Field(i).Name
		childData[fieldName], _ = json.Marshal(val.Field(i).Interface())
	}
	return childData
}

func hydrateChildData(data map[string][]byte, onto any) any {
	val := reflect.ValueOf(onto)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		if !typ.Field(i).IsExported() {
			continue
		}
		fieldName := typ.Field(i).Name
		if val.Field(i).CanSet() {
			_ = json.Unmarshal(data[fieldName], val.Field(i).Addr().Interface())
		}
	}

	return val.Interface()
}

func (e *BaseNestedChild) SetChildID(id string) {
	e._childID = id
}

func (e *BaseNestedChild) ChildID() string {
	return e._childID
}

var (
	typeOfTime         = reflect.TypeOf(time.Time{})
	typeOfSecretString = reflect.TypeOf(SecretString{})
	typeOfAmount       = reflect.TypeOf(Amount{})
	typeOfStringSlice  = reflect.TypeOf([]string{})
)

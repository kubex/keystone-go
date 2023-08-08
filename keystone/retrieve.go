package keystone

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/kubex/keystone-go/proto"
)

type GenericResult map[string]interface{}

func (a *Actor) GetByID(ctx context.Context, entityID string, dst interface{}) error {
	resp, err := a.connection.ProtoClient().Retrieve(ctx, &proto.EntityRequest{
		Authorization: a.authorization(),
		EntityId:      entityID,
		Properties: []*proto.PropertyRequest{
			{Properties: []string{"address~"}},
			{Properties: []string{"name", "email", "city", "state", "country", "postcode", "amount_paid"}},
		},
	})
	if err != nil {
		return err
	}

	if gr, ok := dst.(GenericResult); ok {
		return UnmarshalGeneric(resp, gr)
	}

	return Unmarshal(resp, dst)
}

func (a *Actor) GetByUnique(ctx context.Context, key, value string, dst interface{}) error {
	schema, registered := a.connection.registerType(dst)
	if !registered {
		// wait for the type to be registered with the keystone server
		a.connection.SyncSchema().Wait()
	}

	schemaID := schema.Id
	if schemaID == "" {
		schemaID = schema.Type
	}

	resp, err := a.connection.ProtoClient().Retrieve(ctx, &proto.EntityRequest{
		Authorization: a.authorization(),
		UniqueId: &proto.IDLookup{
			Property: key,
			UniqueId: value,
			SchemaId: schemaID,
		},
		Properties: []*proto.PropertyRequest{
			{Properties: []string{"address~"}},
			{Properties: []string{"name", "email", "city", "state", "country", "postcode", "amount_paid", "lead_date"}},
		},
	})

	if err != nil {
		return err
	}

	return Unmarshal(resp, dst)
}

func makeEntityPropertyMap(resp *proto.EntityResponse) map[string]*proto.EntityProperty {
	//log.Println(resp.GetProperties())
	entityPropertyMap := map[string]*proto.EntityProperty{}
	for _, p := range resp.GetProperties() {
		entityPropertyMap[p.Property.Key] = p
	}
	return entityPropertyMap
}

func Unmarshal(resp *proto.EntityResponse, dst interface{}) error {
	entityPropertyMap := makeEntityPropertyMap(resp)
	return entityResponseToDst(entityPropertyMap, dst, "")
}

func entityResponseToDst(entityPropertyMap map[string]*proto.EntityProperty, dst interface{}, prefix string) error {
	dstVal := reflect.ValueOf(dst)
	fmt.Println("entityResponseToDst", dstVal, dstVal.Type(), prefix)
	for dstVal.Kind() == reflect.Pointer || dstVal.Kind() == reflect.Interface {
		dstVal = dstVal.Elem()
	}
	for i := 0; i < dstVal.NumField(); i++ {
		field := dstVal.Type().Field(i)
		fieldValue := dstVal.Field(i)
		fieldOpt := getFieldOptions(field)
		fieldOpt.name = prefix + fieldOpt.name
		if supportedType(field.Type) {
			setFieldValue(field, fieldValue, fieldOpt, entityPropertyMap)
		} else if field.Type.Kind() == reflect.Struct {
			entityResponseToDst(entityPropertyMap, fieldValue.Addr().Interface(), fieldOpt.name+".")
			continue
		}
	}

	return nil
}

func UnmarshalGeneric(resp *proto.EntityResponse, dst GenericResult) error {
	entityPropertyMap := makeEntityPropertyMap(resp)
	for _, p := range entityPropertyMap {
		if p.Value.GetText() != "" {
			dst[p.Property.Key] = p.Value.GetText()
		}
		if p.Value.GetInt() != 0 {
			dst[p.Property.Key] = p.Value.GetInt()
		}
		if p.Value.GetBool() {
			dst[p.Property.Key] = p.Value.GetBool()
		}
		if p.Value.GetFloat() != 0 {
			dst[p.Property.Key] = p.Value.GetFloat()
		}
		if p.Value.GetSecureText() != "" {
			dst[p.Property.Key] = p.Value.GetSecureText()
		}
		if len(p.Value.GetSet()) > 0 {
			dst[p.Property.Key] = p.Value.GetSet()
		}
		if len(p.Value.GetMap()) > 0 {
			dst[p.Property.Key] = p.Value.GetMap()
		}
		if p.Value.GetTime() != nil {
			dst[p.Property.Key] = time.Unix(p.Value.GetTime().Seconds, int64(p.Value.GetTime().Nanos))
		}
	}
	return nil
}

func setFieldValue(field reflect.StructField, fieldValue reflect.Value, fieldOpt fieldOptions, entityPropertyMap map[string]*proto.EntityProperty) {
	storedProperty, ok := entityPropertyMap[fieldOpt.name]
	if !ok {
		//fmt.Println("no value", fieldOpt.name)
		return
	}
	fmt.Println("setting field", fieldOpt.name, storedProperty.Value, field.Name)
	switch field.Type.Kind() {
	case reflect.String:
		fieldValue.SetString(storedProperty.Value.Text)
	case reflect.Int32, reflect.Int64, reflect.Int:
		fieldValue.SetInt(storedProperty.Value.Int)
	case reflect.Bool:
		fieldValue.SetBool(storedProperty.Value.Bool)
	case reflect.Float32, reflect.Float64:
		fieldValue.SetFloat(storedProperty.Value.Float)
	case reflect.Map:
		for k, v := range storedProperty.Value.Map {
			fieldValue.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
		}
	case reflect.Slice:
		for _, v := range storedProperty.Value.Set {
			fieldValue.Set(reflect.Append(fieldValue, reflect.ValueOf(v)))
		}
	}

	switch field.Type {
	case typeOfSecretString:
		if _, ok := fieldValue.Interface().(SecretString); ok {
			fieldValue.Set(reflect.ValueOf(SecretString{
				Masked:   storedProperty.Value.SecureText,
				Original: storedProperty.Value.Text,
			}))
		}
	case typeOfAmount:
		if _, ok := fieldValue.Interface().(Amount); ok {
			fieldValue.Set(reflect.ValueOf(Amount{
				Currency: storedProperty.Value.Text,
				Units:    storedProperty.Value.Int,
			}))
		}
	case typeOfTime:
		if _, ok := fieldValue.Interface().(time.Time); ok {
			fieldValue.Set(reflect.ValueOf(time.Unix(storedProperty.Value.Int, 0)))
		}
	}
}

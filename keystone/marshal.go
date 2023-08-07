package keystone

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/kubex/keystone-go/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (a *Actor) Marshal(src interface{}, comment string) {
	//log.Println("Processing Marshal request")
	schema, registered := a.connection.registerType(src)
	if !registered {
		// wait for the type to be registered with the keystone server
		a.connection.SyncSchema().Wait()
	}
	//log.Println("Marshalling entity", src)

	v := reflect.ValueOf(src)
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	t := v.Type()

	mutation := &proto.Mutation{
		Mutator:    a.mutator,
		Comment:    comment,
		Properties: []*proto.EntityProperty{},
	}

	eid := "" // try to get from src

	properties := fieldsToProperties(v, t, "")
	for _, p := range properties {
		if p.Property.Key[0] == '_' {
			if p.Property.Key == "_entity_id" && p.Value.Text != "" {
				eid = p.Value.Text
			}
		} else {
			mutation.Properties = append(mutation.Properties, p)
		}
	}

	m := &proto.MutateRequest{
		Authorization: &proto.Authorization{WorkspaceId: a.workspaceID, Source: &a.connection.appID},
		EntityId:      eid,
		Schema:        &proto.Key{Key: schema.Type, Source: schema.Source}, // TODO: Should probably provide the schema ID if we have it - and verify against the type / source
		Mutation:      mutation,
	}

	_, err := a.connection.ProtoClient().Mutate(context.Background(), m)
	if err != nil {
		log.Println(err)
	} else {
		//log.Println(res)
	}

}

func fieldsToProperties(value reflect.Value, t reflect.Type, prefix string) []*proto.EntityProperty {
	var returnProperties []*proto.EntityProperty

	// Iterate all the fields in the struct
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := value.Field(i)
		if field.Anonymous {
			at := field.Type
			if at.Kind() == reflect.Pointer {
				at = at.Elem()
				fieldValue = fieldValue.Elem()
			}

			if !field.IsExported() && at.Kind() != reflect.Struct {
				// ignore embedded unexported fields
				//fmt.Println("skipping unexported anonymous ", field.Name)
				continue
			}

			returnProperties = append(returnProperties, fieldsToProperties(fieldValue, at, prefix)...)
			continue

		} else if !field.IsExported() {
			//fmt.Println("skipping unexported field ", field.Name)
			continue
		}

		fOpt := getFieldOptions(field)
		if fOpt.name == "" {
			continue
		}
		fOpt.name = prefix + fOpt.name

		if !supportedType(field.Type) {
			if field.Type.Kind() == reflect.Struct {
				returnProperties = append(returnProperties, fieldsToProperties(fieldValue, field.Type, fOpt.name+".")...)
			} else {
				//fmt.Println("skipping unsupported type ", field.Type.Kind())
			}
			continue
		}

		protoProp := &proto.EntityProperty{}
		protoProp.Property = &proto.Key{Key: fOpt.name}
		var isEmpty bool
		protoProp.Value, isEmpty = propertyFromField(fieldValue, field)

		if fOpt.omitempty && (protoProp.Value == nil || isEmpty) {
			continue
		}

		returnProperties = append(returnProperties, protoProp)
	}

	return returnProperties
}

func supportedType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.String, reflect.Int32, reflect.Int64, reflect.Int, reflect.Bool, reflect.Float32, reflect.Float64, reflect.Map, reflect.Slice:
		return true
	}

	switch t {
	case typeOfSecretString, typeOfAmount, typeOfTime:
		return true
	}

	return false
}

func propertyFromField(val reflect.Value, fieldType reflect.StructField) (*proto.Value, bool) {
	prop := &proto.Value{}

	switch fieldType.Type.Kind() {
	case reflect.String:
		prop.Text = val.String()
		return prop, prop.Text == ""
	case reflect.Int32, reflect.Int64, reflect.Int:
		prop.Int = val.Int()
		return prop, prop.Int == 0
	case reflect.Bool:
		prop.Bool = val.Bool()
		return prop, !prop.Bool
	case reflect.Float32, reflect.Float64:
		prop.Float = val.Float()
		return prop, prop.Float == 0
	case reflect.Map:
		prop.Map = map[string][]byte{}
		iter := val.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()
			prop.Map[k.String()] = []byte(v.String())
		}
		return prop, len(prop.Map) == 0
	case reflect.Slice:
		if set, ok := val.Slice(0, val.Len()).Interface().([]string); ok {
			prop.Set = set
		} else {
			fmt.Println("only slice of string is supported")
		}
		return prop, len(prop.Set) == 0
	}

	switch fieldType.Type {
	case typeOfSecretString:
		if iVal, ok := val.Interface().(SecretString); ok {
			prop.Text = iVal.Masked
			prop.SecureText = iVal.Original
			return prop, prop.Text == "" && prop.SecureText == ""
		}
	case typeOfAmount:
		if iVal, ok := val.Interface().(Amount); ok {
			prop.Text = iVal.Currency
			prop.Int = iVal.Units
			return prop, prop.Text == "" && prop.Int == 0
		}
	case typeOfTime:
		if iVal, ok := val.Interface().(time.Time); ok {
			prop.Time = timestamppb.New(iVal)
			return prop, prop.Time == nil
		}
	}

	return nil, true
}

package keystone

import (
	"github.com/kubex/keystone-go/proto"
	"log"
	"reflect"
	"strings"
	"time"
)

var (
	typeOfTime         = reflect.TypeOf(time.Time{})
	typeOfSecretString = reflect.TypeOf(SecretString{})
)

const MarshalFieldId = "_entity_id"
const MarshalState = "_state"
const MarshalStateChange = "_state_change"
const MarshalCreated = "_created"
const MarshalSchema = "_schema"
const MarshalSchemaFull = "_schema_full"
const MarshalSchemaVendor = "_schema_vendor"
const MarshalSchemaApp = "_schema_app"

func getFieldName(f reflect.StructField) string {
	tag := f.Tag.Get("keystone")
	if tag == "" {
		return snakeCase(f.Name)
	} else if tag == "-" {
		return ""
	}
	return strings.ToLower(tag)
}

func Unmarshal(response *proto.EntityResponse, dst interface{}) error {
	v := reflect.ValueOf(dst)
	t := v.Type().Elem()

	propNameMap := make(map[string]*proto.Property, 0)
	for _, p := range response.Properties {
		propNameMap[p.GetName()] = p
	}

	// Iterate all the fields in the struct
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		log.Println(field.Name, field.Anonymous, field.IsExported())
		log.Println(field.Type.String(), field.Type.Kind().String())
		fName := getFieldName(field)

		switch fName {
		case MarshalFieldId:
			v.Elem().Field(i).SetString(response.GetEntityId())
			continue
		case MarshalState:
			v.Elem().Field(i).SetString(response.GetState().String())
			continue
		case MarshalStateChange:
			if v.Elem().Field(i).Type() == typeOfTime {
				v.Elem().Field(i).Set(reflect.ValueOf(response.StateChange))
			} else {
				log.Println("WARNING: Field _state_change is not a time.Time")
			}
			continue
		case MarshalCreated:
			if v.Elem().Field(i).Type() == typeOfTime {
				v.Elem().Field(i).Set(reflect.ValueOf(response.Created))
			} else {
				log.Println("WARNING: Field _created is not a time.Time")
			}
			continue
		case MarshalSchema:
			v.Elem().Field(i).SetString(response.GetSchema().GetKey())
			continue
		case MarshalSchemaFull:
			v.Elem().Field(i).SetString(response.GetSchema().String())
			continue
		case MarshalSchemaVendor:
			v.Elem().Field(i).SetString(response.GetSchema().GetVendorId())
			continue
		case MarshalSchemaApp:
			v.Elem().Field(i).SetString(response.GetSchema().GetAppId())
			continue
		}

		prop, ok := propNameMap[fName]
		if !ok {
			continue
		}

		switch field.Type.Kind() {
		case reflect.String:
			if prop.GetSecureText() != "" {
				v.Elem().Field(i).SetString(prop.GetSecureText())
			} else {
				v.Elem().Field(i).SetString(prop.GetText())
			}
			continue
		case reflect.Bool:
			v.Elem().Field(i).SetBool(prop.GetBool())
			continue
		case reflect.Int32, reflect.Int64:
			v.Elem().Field(i).SetInt(prop.GetInt())
			continue
		case reflect.Float32, reflect.Float64:
			v.Elem().Field(i).SetFloat(float64(prop.GetFloat()))
			continue
		}

		switch field.Type {
		case typeOfSecretString:
			v.Elem().Field(i).Set(reflect.ValueOf(SecretString{
				Masked:   prop.GetText(),
				Original: prop.GetSecureText(),
			}))
			continue
		case typeOfTime:
			v.Elem().Field(i).Set(reflect.ValueOf(prop.GetTime().AsTime()))
			continue
		}

		log.Println("Failed Prop ", field.Type.String())
	}

	return nil
}

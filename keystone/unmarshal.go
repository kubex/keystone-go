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
	typeOfAmount       = reflect.TypeOf(Amount{})
)

type fieldOptions struct {
	name      string
	indexed   bool
	lookup    bool
	omitempty bool
}

const MarshalFieldId = "_entity_id"
const MarshalState = "_state"
const MarshalStateChange = "_state_change"
const MarshalCreated = "_created"
const MarshalSchema = "_schema"
const MarshalSchemaFull = "_schema_full"
const MarshalSchemaVendor = "_schema_vendor"
const MarshalSchemaApp = "_schema_app"

func getFieldOptions(f reflect.StructField) fieldOptions {
	tag := f.Tag.Get("keystone")
	opt := fieldOptions{}

	tagParts := strings.Split(tag, ",")
	for i, part := range tagParts {
		if i == 0 {
			if part == "" {
				opt.name = snakeCase(f.Name)
			} else if part == "-" {
				return opt
			} else {
				opt.name = strings.ToLower(part)
			}
			continue
		}
		switch part {
		case "indexed":
			opt.indexed = true
		case "lookup":
			opt.lookup = true
		case "omitempty":
			opt.omitempty = true
		}
	}
	return opt
}

func Unmarshal(response *proto.EntityResponse, dst interface{}) error {
	v := reflect.ValueOf(dst)
	t := v.Type().Elem()

	appRoot := response.GetSchema().GetVendorId() + "/" + response.GetSchema().GetAppId() + "/"

	propNameMap := make(map[string]*proto.Property, 0)
	for _, p := range response.Properties {
		propNameMap[strings.TrimPrefix(p.GetName(), appRoot)] = p
	}

	// Iterate all the fields in the struct
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fOpt := getFieldOptions(field)

		switch fOpt.name {
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
		case "":
			continue
		}

		prop, ok := propNameMap[fOpt.name]
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
		case typeOfAmount:
			v.Elem().Field(i).Set(reflect.ValueOf(Amount{
				Currency: prop.GetText(),
				Units:    prop.GetInt(),
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

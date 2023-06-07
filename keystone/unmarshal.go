package keystone

import (
	"fmt"
	"github.com/kubex/keystone-go/proto"
	"reflect"
)

import (
	"errors"
	"log"
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
		part = strings.TrimSpace(part)
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
		case "indexed", "index":
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

		if !field.IsExported() {
			fmt.Println("skipping unexported field ", field.Name)
			continue
		}

		if err := hydrateByTagName(fOpt.name, v, response, i); err == nil {
			continue
		}

		prop, ok := propNameMap[fOpt.name]
		if !ok {
			//TODO: This code breaks known structs, e.g. time.Time
			/*if field.Type.Kind() == reflect.Struct {
				n := handleStruct(field, propNameMap)
				v.Elem().Field(i).Set(n.Elem())
			}*/
			continue
		}

		if err := hydrateFromProperty(v, field, prop, i); err == nil {
			continue
		} else {
			log.Println(err.Error())
		}
	}

	return nil
}

func hydrateFromProperty(v reflect.Value, field reflect.StructField, prop *proto.Property, i int) error {
	switch field.Type.Kind() {
	case reflect.String:
		if prop.GetValue().GetSecureText() != "" {
			v.Elem().Field(i).SetString(prop.GetValue().GetSecureText())
		} else {
			v.Elem().Field(i).SetString(prop.GetValue().GetText())
		}
		return nil
	case reflect.Bool:
		v.Elem().Field(i).SetBool(prop.GetValue().GetBool())
		return nil
	case reflect.Int32, reflect.Int64:
		v.Elem().Field(i).SetInt(prop.GetValue().GetInt())
		return nil
	case reflect.Float32, reflect.Float64:
		v.Elem().Field(i).SetFloat(float64(prop.GetValue().GetFloat()))
		return nil
	}

	switch field.Type {
	case typeOfSecretString:
		v.Elem().Field(i).Set(reflect.ValueOf(SecretString{
			Masked:   prop.GetValue().GetText(),
			Original: prop.GetValue().GetSecureText(),
		}))
		return nil
	case typeOfAmount:
		v.Elem().Field(i).Set(reflect.ValueOf(Amount{
			Currency: prop.GetValue().GetText(),
			Units:    prop.GetValue().GetInt(),
		}))
		return nil
	case typeOfTime:
		v.Elem().Field(i).Set(reflect.ValueOf(prop.GetValue().GetTime().AsTime()))
		return nil
	}

	return fmt.Errorf("failed Prop %s", field.Type.String())
}

func hydrateByTagName(name string, v reflect.Value, response *proto.EntityResponse, i int) error {
	switch name {
	case MarshalFieldId:
		v.Elem().Field(i).SetString(response.GetEntityId())
		return nil
	case MarshalState:
		v.Elem().Field(i).SetString(response.GetState().String())
		return nil
	case MarshalStateChange:
		if v.Elem().Field(i).Type() == typeOfTime {
			v.Elem().Field(i).Set(reflect.ValueOf(response.StateChange))
		} else {
			log.Println("WARNING: Field _state_change is not a time.Time")
		}
		return nil
	case MarshalCreated:
		if v.Elem().Field(i).Type() == typeOfTime {
			v.Elem().Field(i).Set(reflect.ValueOf(response.Created.AsTime()))
		} else {
			log.Println("WARNING: Field _created is not a time.Time")
		}
		return nil
	case MarshalSchema:
		v.Elem().Field(i).SetString(response.GetSchema().GetKey())
		return nil
	case MarshalSchemaFull:
		v.Elem().Field(i).SetString(response.GetSchema().String())
		return nil
	case MarshalSchemaVendor:
		v.Elem().Field(i).SetString(response.GetSchema().GetVendorId())
		return nil
	case MarshalSchemaApp:
		v.Elem().Field(i).SetString(response.GetSchema().GetAppId())
		return nil
	case "":
		return nil
	}
	return errors.New("unsupported schema field")
}

func handleStruct(field reflect.StructField, propNameMap map[string]*proto.Property) reflect.Value {
	e := field.Type
	n := reflect.New(e)
	for x := 0; x < e.NumField(); x++ {
		nestedField := e.Field(x)
		nestedFieldOpts := getFieldOptions(nestedField)

		if !nestedField.IsExported() {
			fmt.Println("skipping unexported nested field ", field.Name)
			continue
		}

		prop, ok := propNameMap[nestedFieldOpts.name]
		if ok {
			//field = nestedField
			_ = hydrateFromProperty(n, nestedField, prop, x)
		}
	}
	return n
}

package keystone

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/kubex/keystone-go/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// PropertyEncoder extracts properties and children from an entity
type PropertyEncoder struct {
	properties []*proto.EntityProperty
	children   []*proto.EntityChild
}

// Marshal extracts properties and children from an entity
func (p *PropertyEncoder) Marshal(entity interface{}) *proto.Mutation {
	entityValue := reflect.ValueOf(entity)
	for entityValue.Kind() == reflect.Ptr || entityValue.Kind() == reflect.Interface {
		entityValue = entityValue.Elem()
	}

	p.fieldsToProperties(entityValue, entityValue.Type(), "")
	return &proto.Mutation{Properties: p.getProperties(), Children: p.children}
}

// getProperties returns the extracted properties
func (p *PropertyEncoder) getProperties() []*proto.EntityProperty {
	var properties []*proto.EntityProperty
	for _, prop := range p.properties {
		if prop.Property.Key[0] != '_' {
			properties = append(properties, prop)
		}
	}

	return properties
}

func (p *PropertyEncoder) fieldsToProperties(value reflect.Value, t reflect.Type, prefix string) {
	for i := 0; i < t.NumField(); i++ {
		field, fieldValue := t.Field(i), value.Field(i)
		if field.Anonymous {
			if field.Type.Kind() == reflect.Pointer {
				field.Type = field.Type.Elem()
				fieldValue = fieldValue.Elem()
			}
			if field.IsExported() || field.Type.Kind() == reflect.Struct {
				p.fieldsToProperties(fieldValue, field.Type, prefix)
			}
		}

		if !field.IsExported() {
			continue //fmt.Println("skipping unexported field ", field.Name)
		}

		fOpt := getFieldOptions(field, prefix)
		if fOpt.name == "" {
			continue
		}

		if supportedType(field.Type) {
			if protoProp, isEmpty := p.entityPropertyFromField(fieldValue, field, fOpt); !isEmpty {
				p.properties = append(p.properties, protoProp)
			}
		} else {
			if field.Type.Kind() == reflect.Pointer {
				fieldValue = fieldValue.Elem()
				if fieldValue.IsValid() {
					field.Type = fieldValue.Type()
				}
			}
			if field.Type.Kind() == reflect.Slice && fieldValue.Len() > 0 {
				if _, child := fieldValue.Index(0).Interface().(EntityChild); child {
					for i := 0; i < fieldValue.Len(); i++ {
						childData, err := json.Marshal(fieldValue.Index(i).Interface())
						if err == nil {
							p.children = append(p.children, &proto.EntityChild{
								Type: &proto.Key{Key: snakeCase(fOpt.name)},
								Data: childData,
							})
						}
					}
				}
			}

			if field.Type.Kind() == reflect.Struct || field.Type.Kind() == reflect.Pointer {
				p.fieldsToProperties(fieldValue, field.Type, fOpt.name+".")
			} else {
				//fmt.Println("skipping unsupported type ", field.Type.Kind())
			}
		}
	}
}

func (p *PropertyEncoder) entityPropertyFromField(fieldValue reflect.Value, fieldType reflect.StructField, fOpt fieldOptions) (*proto.EntityProperty, bool) {
	prop := &proto.EntityProperty{Property: &proto.Key{Key: fOpt.name}, Value: &proto.Value{}}

	switch fieldType.Type.Kind() {
	case reflect.String:
		prop.Value.Text = fieldValue.String()
		return prop, prop.Value.Text == ""
	case reflect.Int32, reflect.Int64, reflect.Int:
		prop.Value.Int = fieldValue.Int()
		return prop, prop.Value.Int == 0
	case reflect.Bool:
		prop.Value.Bool = fieldValue.Bool()
		return prop, !prop.Value.Bool
	case reflect.Float32, reflect.Float64:
		prop.Value.Float = fieldValue.Float()
		return prop, prop.Value.Float == 0
	case reflect.Map:
		prop.Value.Map = map[string][]byte{}
		iter := fieldValue.MapRange()
		for iter.Next() {
			prop.Value.Map[iter.Key().String()] = []byte(iter.Value().String())
		}
		return prop, len(prop.Value.Map) == 0
	case reflect.Slice:
		if set, ok := fieldValue.Interface().([]string); ok {
			prop.Value.Set = set
		} else {
			fmt.Println("only []string is supported")
		}
		return prop, len(prop.Value.Set) == 0
	}

	switch fieldType.Type {
	case typeOfSecretString:
		if iVal, ok := fieldValue.Interface().(SecretString); ok {
			prop.Value.Text = iVal.Masked
			prop.Value.SecureText = iVal.Original
			return prop, prop.Value.Text == "" && prop.Value.SecureText == ""
		}
	case typeOfAmount:
		if iVal, ok := fieldValue.Interface().(Amount); ok {
			prop.Value.Text = iVal.Currency
			prop.Value.Int = iVal.Units
			return prop, prop.Value.Text == "" && prop.Value.Int == 0
		}
	case typeOfTime:
		if iVal, ok := fieldValue.Interface().(time.Time); ok {
			prop.Value.Time = timestamppb.New(iVal)
			return prop, prop.Value.Time == nil
		}
	}

	return prop, true
}

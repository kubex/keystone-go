package keystone

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/kubex/keystone-go/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// EntityIDKey is the key for the entity ID property
const EntityIDKey = "_entity_id"

// PropertyExtractor extracts properties and children from an entity
type PropertyExtractor struct {
	EntityID   string
	properties []*proto.EntityProperty
	children   []*proto.EntityChild
}

// Extract extracts properties and children from an entity
func (p *PropertyExtractor) Extract(entity interface{}) error {
	entityValue := reflect.ValueOf(entity)
	for entityValue.Kind() == reflect.Ptr || entityValue.Kind() == reflect.Interface {
		entityValue = entityValue.Elem()
	}

	p.fieldsToProperties(entityValue, entityValue.Type(), "")
	return nil
}

// Properties returns the extracted properties
func (p *PropertyExtractor) Properties() []*proto.EntityProperty {
	var properties []*proto.EntityProperty
	for _, prop := range p.properties {
		if prop.Property.Key[0] == '_' {
			if prop.Property.Key == EntityIDKey && prop.Value.Text != "" {
				p.EntityID = prop.Value.Text
			}
		} else {
			properties = append(properties, prop)
		}
	}

	return properties
}

// Children returns the extracted children
func (p *PropertyExtractor) Children() []*proto.EntityChild {
	return p.children
}

func (p *PropertyExtractor) fieldsToProperties(value reflect.Value, t reflect.Type, prefix string) {
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

			p.fieldsToProperties(fieldValue, at, prefix)
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
						if err != nil {
							continue
						}
						p.children = append(p.children, &proto.EntityChild{
							Type: &proto.Key{Key: snakeCase(fOpt.name)},
							Data: childData,
						})
					}
				}
			}

			if field.Type.Kind() == reflect.Struct || field.Type.Kind() == reflect.Pointer {
				p.fieldsToProperties(fieldValue, field.Type, fOpt.name+".")
			} else {
				fmt.Println("skipping unsupported type ", field.Type.Kind())
			}
			continue
		}

		protoProp := &proto.EntityProperty{}
		protoProp.Property = &proto.Key{Key: fOpt.name}
		var isEmpty bool
		protoProp.Value, isEmpty = p.propertyValueFromField(fieldValue, field)

		if fOpt.omitempty && (protoProp.Value == nil || isEmpty) {
			continue
		}

		p.properties = append(p.properties, protoProp)
	}
}

func (p *PropertyExtractor) propertyValueFromField(fieldValue reflect.Value, fieldType reflect.StructField) (*proto.Value, bool) {
	prop := &proto.Value{}

	switch fieldType.Type.Kind() {
	case reflect.String:
		prop.Text = fieldValue.String()
		return prop, prop.Text == ""
	case reflect.Int32, reflect.Int64, reflect.Int:
		prop.Int = fieldValue.Int()
		return prop, prop.Int == 0
	case reflect.Bool:
		prop.Bool = fieldValue.Bool()
		return prop, !prop.Bool
	case reflect.Float32, reflect.Float64:
		prop.Float = fieldValue.Float()
		return prop, prop.Float == 0
	case reflect.Map:
		prop.Map = map[string][]byte{}
		iter := fieldValue.MapRange()
		for iter.Next() {
			prop.Map[iter.Key().String()] = []byte(iter.Value().String())
		}
		return prop, len(prop.Map) == 0
	case reflect.Slice:
		if set, ok := fieldValue.Interface().([]string); ok {
			prop.Set = set
		} else {
			fmt.Println("only []string is supported")
		}
		return prop, len(prop.Set) == 0
	}

	switch fieldType.Type {
	case typeOfSecretString:
		if iVal, ok := fieldValue.Interface().(SecretString); ok {
			prop.Text = iVal.Masked
			prop.SecureText = iVal.Original
			return prop, prop.Text == "" && prop.SecureText == ""
		}
	case typeOfAmount:
		if iVal, ok := fieldValue.Interface().(Amount); ok {
			prop.Text = iVal.Currency
			prop.Int = iVal.Units
			return prop, prop.Text == "" && prop.Int == 0
		}
	case typeOfTime:
		if iVal, ok := fieldValue.Interface().(time.Time); ok {
			prop.Time = timestamppb.New(iVal)
			return prop, prop.Time == nil
		}
	}

	return nil, true
}

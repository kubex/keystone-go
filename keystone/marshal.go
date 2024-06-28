package keystone

import (
	"fmt"
	"reflect"
	"time"

	"github.com/kubex/keystone-go/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// PropertyEncoder extracts properties and Children from an entity
type PropertyEncoder struct {
	properties []*proto.EntityProperty
	children   []*proto.EntityChild
}

// Marshal extracts properties and Children from an entity
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
		if prop.Property[0] != '_' {
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
			if protoProp, isEmpty := entityPropertyFromField(fieldValue, field.Type, fOpt); !isEmpty {
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
				firstChild := fieldValue.Index(0).Interface()
				if _, ok := firstChild.(NestedChild); ok {
					for i := 0; i < fieldValue.Len(); i++ {
						ch := fieldValue.Index(i).Interface()
						child := ch.(NestedChild)
						ech := &proto.EntityChild{
							Type: &proto.Key{Key: snakeCase(fOpt.name)},
						}

						if childData, ok := ch.(NestedChildAggregateValue); ok {
							ech.Value = childData.AggregateValue()
						}

						if childData, ok := child.(NestedChildData); ok {
							ech.Data = childData.KeystoneData()
						} else {
							ech.Data = getChildData(ch)
						}

						p.children = append(p.children, ech)
					}
				} else {
					fmt.Println("skipping unsupported slice type ", field.Type.Kind(), fieldValue.Index(0).Interface().(NestedChild))
				}
			}

			if field.Type.Kind() == reflect.Struct || field.Type.Kind() == reflect.Pointer {
				p.fieldsToProperties(fieldValue, field.Type, fOpt.name+".")
			} else {
				fmt.Println("skipping unsupported type ", field.Type.Kind())
			}
		}
	}
}

func entityPropertyFromField(fieldValue reflect.Value, fieldType reflect.Type, fOpt fieldOptions) (*proto.EntityProperty, bool) {
	prop := &proto.EntityProperty{Property: fOpt.name, Value: &proto.Value{}}
	switch fieldType.Kind() {
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
		prop.Value.Map = map[string]string{}
		iter := fieldValue.MapRange()
		for iter.Next() {
			var mv string
			if _, ok := iter.Value().Interface().(string); ok {
				mv = iter.Value().String()
				/*} else if iter.Value().Int() != 0 {
				mv = strconv.Itoa(int(iter.Value().Int()))
				*/
			} else {
				fmt.Println("only map[string]string is supported (" + fOpt.name + ")")
			}
			prop.Value.Map[iter.Key().String()] = mv
		}
		return prop, len(prop.Value.Map) == 0
	case reflect.Slice:
		if set, ok := fieldValue.Interface().([]string); ok {
			prop.Value.Set = set
		} else {
			fmt.Println("only []string is supported (" + fOpt.name + ")")
		}
		return prop, len(prop.Value.Set) == 0
	}

	switch fieldType {
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
			if !iVal.IsZero() {
				prop.Value.Time = timestamppb.New(iVal)
				return prop, prop.Value.Time == nil
			}
		}
	}

	return prop, true
}

package keystone

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/kubex/keystone-go/proto"
)

func UnmarshalAppend(dstPtr any, resp ...*proto.EntityResponse) error {
	dstT := reflect.TypeOf(dstPtr)
	if dstT.Kind() != reflect.Pointer || dstT.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("dst must be a slice pointer")
	}

	valuePtr := reflect.ValueOf(dstPtr)
	sliceLen := valuePtr.Elem().Len()
	appendLen := len(resp)
	dst := valuePtr.Elem()

	// Pointer > Slice > Slice Element
	elementType := dstT.Elem().Elem()
	pointer := elementType.Kind() == reflect.Ptr
	finalSlice := reflect.MakeSlice(reflect.SliceOf(elementType), sliceLen+appendLen, sliceLen+appendLen)

	for x := 0; x < sliceLen; x++ {
		finalSlice.Index(x).Set(dst.Index(x))
	}

	if pointer {
		elementType = elementType.Elem()
	}

	for i, r := range resp {
		dstEle := reflect.New(elementType)
		ifa := dstEle.Interface()
		if err := Unmarshal(r, ifa); err != nil {
			return err
		}
		val := reflect.ValueOf(ifa)
		if pointer {
			val = reflect.ValueOf(&ifa)
		}
		if val.Kind() == reflect.Ptr {
			val = reflect.ValueOf(val.Elem().Interface())
		}
		finalSlice.Index(i + sliceLen).Set(val)
	}

	dst.Set(finalSlice)

	return nil
}

func Unmarshal(resp *proto.EntityResponse, dst interface{}) error {
	//TODO: Support _count_relation_::type
	//TODO: Support _count_relation_vendor:app:type
	entityPropertyMap := makeEntityPropertyMap(resp)
	if entityWithLinks, ok := dst.(EntityLinkProvider); ok {
		entityWithLinks.SetKeystoneLinks(resp.GetLinks())
	}
	if entityWithRelationships, ok := dst.(EntityRelationshipProvider); ok {
		entityWithRelationships.SetKeystoneRelationships(resp.GetRelationships())
	}
	err := entityResponseToDst(entityPropertyMap, resp.Children, dst, "")

	if baseEntity, ok := dst.(Entity); ok {
		baseEntity.SetKeystoneID(resp.GetEntity().GetEntityId())
	}

	return err
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

func makeEntityPropertyMap(resp *proto.EntityResponse) map[string]*proto.EntityProperty {
	entityPropertyMap := map[string]*proto.EntityProperty{}
	for _, p := range resp.GetProperties() {
		entityPropertyMap[p.Property.Key] = p
	}
	return entityPropertyMap
}

func entityResponseToDst(entityPropertyMap map[string]*proto.EntityProperty, children []*proto.EntityChild, dst interface{}, prefix string) error {
	dstVal := reflect.ValueOf(dst)
	//fmt.Println("entityResponseToDst", dstVal, dstVal.Type(), prefix)
	for dstVal.Kind() == reflect.Pointer || dstVal.Kind() == reflect.Interface {
		dstVal = dstVal.Elem()
	}
	for i := 0; i < dstVal.NumField(); i++ {
		field := dstVal.Type().Field(i)
		fieldValue := dstVal.Field(i)
		fieldOpt := getFieldOptions(field, prefix)
		fieldOpt.name = prefix + fieldOpt.name
		if supportedType(field.Type) {
			setFieldValue(field, fieldValue, fieldOpt, entityPropertyMap)
		} else {
			// hydrate children
			if field.Type.Kind() == reflect.Slice && len(children) > 0 {
				for _, child := range children {
					if child.Type.Key == fieldOpt.name {
						el := reflect.New(field.Type.Elem())
						if err := json.Unmarshal(child.Data, el.Interface()); err != nil {
							continue
						}
						fieldValue.Set(reflect.Append(fieldValue, el.Elem()))
					}
				}
				continue
			}

			if field.Type.Kind() == reflect.Struct {
				if err := entityResponseToDst(entityPropertyMap, children, fieldValue.Addr().Interface(), fieldOpt.name+"."); err != nil {
					return err
				}
				continue
			}
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
	//fmt.Println("setting field", fieldOpt.name, storedProperty.Value, field.Name)
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
	}

	switch field.Type {
	case typeOfStringSlice:
		for _, v := range storedProperty.Value.Set {
			fieldValue.Set(reflect.Append(fieldValue, reflect.ValueOf(v)))
		}
	case typeOfSecretString:
		if _, ok := fieldValue.Interface().(SecretString); ok {
			fieldValue.Set(reflect.ValueOf(SecretString{
				Masked:   storedProperty.Value.Text,
				Original: storedProperty.Value.SecureText,
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
			if storedProperty.Value.Time != nil {
				t := time.Unix(storedProperty.Value.GetTime().Seconds, int64(storedProperty.Value.GetTime().Nanos))
				fieldValue.Set(reflect.ValueOf(t))
			}
		}
	}
}

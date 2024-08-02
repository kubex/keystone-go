package keystone

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/kubex/keystone-go/proto"
)

func UnmarshalAppend(dstPtr any, resp ...*proto.EntityResponse) error {
	dstT := reflect.TypeOf(dstPtr)
	if dstT.Kind() != reflect.Pointer || dstT.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("dst must be a slice pointer")
	}

	sort.Sort(proto.EntityResponseIDSort(resp))

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
	entityPropertyMap := makeEntityPropertyMap(resp)

	if resp.GetEntity() != nil {
		e := resp.GetEntity()
		split := strings.Split(e.GetEntityId(), "-")
		entityPropertyMap["_entity_id"] = &proto.EntityProperty{Value: valueFromAny(split[0])}
		if len(split) == 2 {
			entityPropertyMap["_child_id"] = &proto.EntityProperty{Value: valueFromAny(split[1])}
		}

		entityPropertyMap["_schema_id"] = &proto.EntityProperty{Value: valueFromAny(e.GetSchemaId())}
		entityPropertyMap["_created"] = &proto.EntityProperty{Value: valueFromAny(e.GetCreated())}
		entityPropertyMap["_state_change"] = &proto.EntityProperty{Value: valueFromAny(e.GetStateChange())}
		entityPropertyMap["_state"] = &proto.EntityProperty{Value: valueFromAny(e.GetState())}
		entityPropertyMap["_last_update"] = &proto.EntityProperty{Value: valueFromAny(e.GetLastUpdate())}
	}

	var countReplace = map[string]int64{}

	if entityDetail, ok := dst.(EntityDetail); ok {
		entityDetail.SetEntityDetail(resp.GetEntity())
	}

	if resp.GetRelationshipCounts() != nil {
		for _, v := range resp.GetRelationshipCounts() {
			t := v.GetType()
			cnt := int64(v.GetCount())
			if t.GetKey() == "" {
				countReplace["_count_relation"] = cnt
			} else {
				countReplace[fmt.Sprintf("_count_relation:%s:%s:%s", t.GetSource().GetVendorId(), t.GetSource().GetAppId(), t.GetKey())] = cnt
				countReplace[fmt.Sprintf("_count_relation:%s:%s", t.GetSource().GetAppId(), t.GetKey())] = cnt
				countReplace[fmt.Sprintf("_count_relation:%s", t.GetKey())] = cnt
			}
		}
	}

	if resp.GetDescendantCounts() != nil {
		for _, v := range resp.GetDescendantCounts() {
			t := v.GetType()
			cnt := int64(v.GetCount())
			if t.GetKey() == "" {
				countReplace["_count_descendant"] = cnt
			} else {
				countReplace[fmt.Sprintf("_count_descendant:%s:%s:%s", t.GetSource().GetVendorId(), t.GetSource().GetAppId(), t.GetKey())] = cnt
				countReplace[fmt.Sprintf("_count_descendant:%s:%s", t.GetSource().GetAppId(), t.GetKey())] = cnt
				countReplace[fmt.Sprintf("_count_descendant:%s", t.GetKey())] = cnt
			}
		}
	}

	for variant, cnt := range countReplace {
		entityPropertyMap[variant] = &proto.EntityProperty{Property: variant, Value: &proto.Value{Int: cnt}}
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
			dst[p.Property] = p.Value.GetText()
		}
		if p.Value.GetInt() != 0 {
			dst[p.Property] = p.Value.GetInt()
		}
		if p.Value.GetBool() {
			dst[p.Property] = p.Value.GetBool()
		}
		if p.Value.GetFloat() != 0 {
			dst[p.Property] = p.Value.GetFloat()
		}
		if p.Value.GetSecureText() != "" {
			dst[p.Property] = p.Value.GetSecureText()
		}
		if len(p.Value.GetSet()) > 0 {
			dst[p.Property] = p.Value.GetSet()
		}
		if len(p.Value.GetMap()) > 0 {
			dst[p.Property] = p.Value.GetMap()
		}
		if p.Value.GetTime() != nil {
			dst[p.Property] = time.Unix(p.Value.GetTime().Seconds, int64(p.Value.GetTime().Nanos))
		}
	}
	return nil
}

func makeEntityPropertyMap(resp *proto.EntityResponse) map[string]*proto.EntityProperty {
	entityPropertyMap := map[string]*proto.EntityProperty{}
	for _, p := range resp.GetProperties() {
		entityPropertyMap[p.Property] = p
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
		} else if field.IsExported() {
			// hydrate Children
			if field.Type.Kind() == reflect.Slice && len(children) > 0 {
				childSlice := reflect.MakeSlice(field.Type, 0, 0)
				for _, child := range children {
					if child.Type.Key == fieldOpt.name {
						el := reflect.New(field.Type.Elem())
						if el.Type().Kind() == reflect.Ptr {
							el = reflect.New(el.Elem().Type().Elem())
						}

						ch := el.Interface()

						if childData, ok := ch.(NestedChild); ok {
							childData.SetChildID(child.Cid)
						}

						if childData, ok := ch.(NestedChildAggregateValue); ok {
							childData.SetAggregateValue(child.Value)
						}

						if childData, ok := ch.(NestedChildData); ok {
							childData.HydrateKeystoneData(child.Data)
						} else {
							hydrateChildData(child.Data, ch)
						}
						childSlice = reflect.Append(childSlice, el)
					}
				}
				if fieldValue.CanSet() {
					fieldValue.Set(childSlice)
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

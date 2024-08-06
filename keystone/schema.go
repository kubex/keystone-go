package keystone

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/kubex/keystone-go/proto"
)

type schemaDef struct {
	schema     *proto.Schema
	definition TypeDefinition
}

func baseType(input interface{}) reflect.Type {
	v := reflect.ValueOf(input)
	t := v.Type()

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func Type(input interface{}) string {
	return getType(baseType(input))
}

func getType(p reflect.Type) string {
	return strings.ReplaceAll(snakeCase(p.Name()), "_", "-")
}

func typeToName(ksType string) string {
	return cases.Title(language.English).String(strings.ReplaceAll(ksType, "-", " "))
}

func typeToSchema(input interface{}) schemaDef {
	t := baseType(input)
	ksType := getType(t)
	returnSchema := &proto.Schema{
		Name: typeToName(ksType),
		Type: ksType,
	}

	if _, ok := input.(ChildEntity); ok {
		returnSchema.IsChild = true
	}

	retDef := schemaDef{schema: returnSchema, definition: TypeDefinition{}}

	if definer, ok := input.(EntityDefinition); ok {
		def := definer.GetKeystoneDefinition()
		retDef.definition = def

		if def.Name != "" {
			returnSchema.Name = def.Name
		}
		if def.Type != "" {
			returnSchema.Type = def.Type
		}

		returnSchema.Description = def.Description
		returnSchema.Singular = def.Singular
		returnSchema.Plural = def.Plural

		returnSchema.Options = append(returnSchema.Options, def.Options...)
	}

	returnSchema.Properties = getProperties(t, "")

	return retDef
}

func getProperties(t reflect.Type, prefix string) []*proto.Property {
	var returnFields []*proto.Property

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t == reflect.TypeOf(GenericResult{}) {
		return returnFields
	}

	// Iterate all the fields in the struct
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			t := field.Type
			if t.Kind() == reflect.Pointer {
				t = t.Elem()
			}

			if !field.IsExported() && t.Kind() != reflect.Struct {
				// ignore embedded unexported fields
				fmt.Println("skipping unexported anonymous ", field.Name)
				continue
			}

			returnFields = append(returnFields, getProperties(t, prefix)...)
			continue

		} else if !field.IsExported() {
			//fmt.Println("skipping unexported field ", field.Name)
			continue
		}

		fOpt := getFieldOptions(field, prefix)
		if fOpt.name == "" {
			continue
		}

		// not supported assumed a nested struct field
		if !supportedType(field.Type) {
			if field.Type.Kind() == reflect.Pointer {
				field.Type = field.Type.Elem()
			}

			if field.Type.Kind() == reflect.Struct {
				returnFields = append(returnFields, getProperties(field.Type, fOpt.name+".")...)
			} else {
				log.Println("skipping unsupported field", field.Name, field.Type.Kind())
			}
			continue
		}

		protoField := &proto.Property{}
		protoField.DataType, protoField.ExtendedType = getFieldType(field)
		fOpt.applyTo(protoField)

		returnFields = append(returnFields, protoField)
	}

	return returnFields
}

func supportedType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.String, reflect.Int32, reflect.Int64, reflect.Int, reflect.Bool, reflect.Float32, reflect.Float64, reflect.Map:
		return true
	}

	if t.Kind() == reflect.Slice {
		switch t.Elem().Kind() {
		case reflect.Uint8, reflect.String, reflect.Int, reflect.Int32, reflect.Int64:
			return true
		}
	}

	switch t {
	case typeOfSecretString, typeOfVerifyString,
		typeOfAmount, typeOfTime,
		typeOfStringSet, typeOfIntSet:
		return true
	}

	return false
}

func appendOption(protoField *proto.Property, option proto.Property_Option, when bool) {
	if when {
		protoField.Options = append(protoField.Options, option)
	}
}

func getFieldType(fieldType reflect.StructField) (proto.Property_Type, proto.Property_ExtendedType) {
	extendedType := proto.Property_ExtendedNone
	switch fieldType.Type.Kind() {
	case reflect.String:
		return proto.Property_Text, extendedType
	case reflect.Int32, reflect.Int64, reflect.Int:
		return proto.Property_Number, extendedType
	case reflect.Bool:
		return proto.Property_Boolean, extendedType
	case reflect.Float32, reflect.Float64:
		return proto.Property_Float, extendedType
	case reflect.Map:
		return proto.Property_KeyValue, extendedType
	case reflect.Slice:
		switch fieldType.Type.Elem().Kind() {
		case reflect.String:
			return proto.Property_Strings, extendedType
		case reflect.Int, reflect.Int64, reflect.Int32:
			return proto.Property_Ints, extendedType
		case reflect.Uint8:
			return proto.Property_Bytes, extendedType
		}
		return proto.Property_Strings, extendedType
	}

	switch fieldType.Type {
	case typeOfSecretString:
		return proto.Property_SecureText, extendedType
	case typeOfAmount:
		return proto.Property_Amount, extendedType
	case typeOfTime:
		return proto.Property_Time, extendedType
	case typeOfVerifyString:
		return proto.Property_VerifyText, extendedType
	case typeOfIntSet:
		return proto.Property_IntSet, extendedType
	case typeOfStringSet:
		return proto.Property_StringSet, extendedType
	}

	return proto.Property_Bytes, extendedType
}

func getFieldOptions(f reflect.StructField, prefix string) fieldOptions {
	tag := f.Tag.Get("keystone")
	opt := fieldOptions{}

	tagParts := strings.Split(tag, ",")
	for i, part := range tagParts {
		part = strings.TrimSpace(part)
		if i == 0 {
			if part == "" {
				opt.name = prefix + snakeCase(f.Name)
			} else if part == "-" {
				return opt
			} else {
				opt.name = prefix + strings.ToLower(part)
			}
			continue
		}
		switch part {
		case "omitempty":
			opt.omitempty = true

		case "unique":
			opt.unique = true
		case "indexed", "query":
			opt.indexed = true
		case "searchable", "search":
			opt.searchable = true
		case "immutable":
			opt.immutable = true
		case "required", "req":
			opt.required = true
		case "lookup":
			opt.reverseLookup = true
		case "verify":
			opt.verifyOnly = true

		case "pii", "personal", "gdpr":
			opt.personalData = true
		case "user":
			opt.userInputData = true
		}
	}
	return opt
}

type fieldOptions struct {
	name string

	// marshal
	omitempty bool

	// options
	unique        bool
	indexed       bool
	searchable    bool
	immutable     bool
	required      bool
	reverseLookup bool
	verifyOnly    bool

	// Data classification
	personalData  bool
	userInputData bool
}

func (fOpt fieldOptions) applyTo(protoField *proto.Property) {
	protoField.Name = fOpt.name
	if fOpt.personalData {
		protoField.ExtendedType = proto.Property_Personal
		protoField.DataType = proto.Property_SecureText
	} else if fOpt.userInputData {
		protoField.ExtendedType = proto.Property_UserInput
	} else if fOpt.verifyOnly {
		protoField.DataType = proto.Property_VerifyText
	}

	appendOption(protoField, proto.Property_Unique, fOpt.unique)
	appendOption(protoField, proto.Property_Indexed, fOpt.indexed)
	appendOption(protoField, proto.Property_Immutable, fOpt.immutable)
	appendOption(protoField, proto.Property_Required, fOpt.required)
	appendOption(protoField, proto.Property_ReverseLookup, fOpt.reverseLookup)
	appendOption(protoField, proto.Property_Searchable, fOpt.searchable)
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func snakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

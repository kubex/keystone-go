package keystone

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/kubex/keystone-go/proto"
)

func baseType(input interface{}) reflect.Type {
	v := reflect.ValueOf(input)
	t := v.Type()

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func Type(input interface{}) string { return ksType(baseType(input)) }
func ksType(p reflect.Type) string  { return strings.ReplaceAll(snakeCase(p.Name()), "_", " ") }

func typeToSchema(input interface{}) *proto.Schema {

	t := baseType(input)
	returnSchema := &proto.Schema{
		Name: ksType(t),
		Type: t.Name(),
	}

	if definer, ok := input.(EntityDefinition); ok {
		def := definer.GetKeystoneDefinition()
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

	return returnSchema
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
			fmt.Println("skipping unexported field ", field.Name)
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
				log.Println("skipping unsupported field ", field.Name, field.Type.Kind())
			}
			continue
		}

		protoField := &proto.Property{}
		protoField.DataType, protoField.Classification = getFieldType(field)
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

	switch t {
	case typeOfSecretString, typeOfAmount, typeOfTime, typeOfStringSlice:
		return true
	}

	return false
}

func appendOption(protoField *proto.Property, option proto.Property_Option, when bool) {
	if when {
		protoField.Options = append(protoField.Options, option)
	}
}

func getFieldType(fieldType reflect.StructField) (proto.Property_Type, proto.Property_Classification) {
	defaultClassification := proto.Property_Anonymous
	switch fieldType.Type.Kind() {
	case reflect.String:
		return proto.Property_Text, defaultClassification
	case reflect.Int32, reflect.Int64, reflect.Int:
		return proto.Property_Number, defaultClassification
	case reflect.Bool:
		return proto.Property_Boolean, defaultClassification
	case reflect.Float32, reflect.Float64:
		return proto.Property_Float, defaultClassification
	case reflect.Map:
		return proto.Property_Map, defaultClassification
	case reflect.Slice:
		return proto.Property_Set, defaultClassification
	}

	switch fieldType.Type {
	case typeOfSecretString:
		return proto.Property_Text, proto.Property_Secure
	case typeOfAmount:
		return proto.Property_Amount, defaultClassification
	case typeOfTime:
		return proto.Property_Time, defaultClassification
	}

	return proto.Property_Text, defaultClassification
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
		case "indexed":
			opt.indexed = true
		case "immutable":
			opt.immutable = true
		case "required":
			opt.required = true
		case "lookup":
			opt.reverseLookup = true

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
	immutable     bool
	required      bool
	reverseLookup bool

	// data classification
	personalData  bool
	userInputData bool
}

func (fOpt fieldOptions) applyTo(protoField *proto.Property) {
	protoField.Name = fOpt.name
	if fOpt.personalData {
		protoField.Classification = proto.Property_Personal
	} else if fOpt.userInputData {
		protoField.Classification = proto.Property_UserInput
	}

	appendOption(protoField, proto.Property_Unique, fOpt.unique)
	appendOption(protoField, proto.Property_Indexed, fOpt.indexed)
	appendOption(protoField, proto.Property_Immutable, fOpt.immutable)
	appendOption(protoField, proto.Property_Required, fOpt.required)
	appendOption(protoField, proto.Property_ReverseLookup, fOpt.reverseLookup)
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func snakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

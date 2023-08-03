package keystone

import (
	"fmt"
	"github.com/kubex/keystone-go/proto"
	"reflect"
	"regexp"
	"strings"
	"time"
)

func typeToSchema(input interface{}) *proto.Schema {
	v := reflect.ValueOf(input)
	t := v.Type()

	name := strings.ReplaceAll(snakeCase(t.Name()), "_", " ")
	returnSchema := &proto.Schema{
		Name: name,
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

	returnSchema.Fields = getFields(t, "")

	return returnSchema
}

func getFields(t reflect.Type, prefix string) []*proto.Field {
	var returnFields []*proto.Field

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

			returnFields = append(returnFields, getFields(t, prefix)...)
			continue

		} else if !field.IsExported() {
			fmt.Println("skipping unexported field ", field.Name)
			continue
		}

		fOpt := getFieldOptions(field)
		if fOpt.name == "" {
			continue
		}
		fOpt.name = prefix + fOpt.name

		protoField := &proto.Field{}
		var supported bool
		protoField.DataType, protoField.Classification, supported = getFieldType(field)

		// not supported assumed a nested struct field
		if !supported {
			returnFields = append(returnFields, getFields(field.Type, fOpt.name+".")...)
			continue
		}
		fOpt.applyTo(protoField)

		returnFields = append(returnFields, protoField)
	}

	return returnFields
}

func appendOption(protoField *proto.Field, option proto.Field_Option, when bool) {
	if when {
		protoField.Options = append(protoField.Options, option)
	}
}

func getFieldType(fieldType reflect.StructField) (proto.Field_Type, proto.Field_Classification, bool) {
	defaultClassification := proto.Field_Anonymous
	switch fieldType.Type.Kind() {
	case reflect.String:
		return proto.Field_Text, defaultClassification, true
	case reflect.Int32, reflect.Int64:
		return proto.Field_Number, defaultClassification, true
	case reflect.Bool:
		return proto.Field_Boolean, defaultClassification, true
	case reflect.Float32, reflect.Float64:
		return proto.Field_Float, defaultClassification, true
	}

	switch fieldType.Type {
	case typeOfSecretString:
		return proto.Field_Text, proto.Field_Secure, true
	case typeOfAmount:
		return proto.Field_Amount, defaultClassification, true
	case typeOfTime:
		return proto.Field_Time, defaultClassification, true
	}

	return proto.Field_Text, defaultClassification, fieldType.Type.Kind() != reflect.Struct
}

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

var (
	typeOfTime         = reflect.TypeOf(time.Time{})
	typeOfSecretString = reflect.TypeOf(SecretString{})
	typeOfAmount       = reflect.TypeOf(Amount{})
)

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

func (fOpt fieldOptions) applyTo(protoField *proto.Field) {
	protoField.Name = fOpt.name
	if fOpt.personalData {
		protoField.Classification = proto.Field_Personal
	} else if fOpt.userInputData {
		protoField.Classification = proto.Field_UserInput
	}

	appendOption(protoField, proto.Field_Unique, fOpt.unique)
	appendOption(protoField, proto.Field_Indexed, fOpt.indexed)
	appendOption(protoField, proto.Field_Immutable, fOpt.immutable)
	appendOption(protoField, proto.Field_Required, fOpt.required)
	appendOption(protoField, proto.Field_ReverseLookup, fOpt.reverseLookup)
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func snakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

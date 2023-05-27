package keystone

import (
	"github.com/kubex/keystone-go/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Property struct {
	Name           string
	Type           PropertyType
	Classification Classification
	Text           string
	Secret         string
	Int            int64
	Bool           bool
	Float          float64
	Time           *time.Time
	updated        bool
}

func (p *Property) toProto() *proto.Property {
	var useTime *timestamppb.Timestamp
	if p.Time != nil {
		useTime = timestamppb.New(*p.Time)
	}
	return &proto.Property{
		Name:           p.Name,
		Type:           p.Type.toProto(),
		Classification: p.Classification.toProto(),
		Text:           p.Text,
		Int:            p.Int,
		Bool:           p.Bool,
		Float:          float32(p.Float),
		Time:           useTime,
		SecureText:     p.Secret,
	}
}

func Secret(name, secureData, preview string) Property {
	return Property{updated: true, Name: name, Secret: secureData, Text: preview, Classification: ClassificationSecure, Type: PropertyTypeText}
}

func Personal(name, sensitiveData, preview string) Property {
	return Property{updated: true, Name: name, Secret: sensitiveData, Text: preview, Classification: ClassificationPersonal, Type: PropertyTypeText}
}

func ID(name, uniqueID string) Property {
	return Property{updated: true, Name: name, Text: uniqueID, Type: PropertyTypeText, Classification: ClassificationID}
}

func Text(name, input string) Property {
	return Property{updated: true, Name: name, Text: input, Type: PropertyTypeText}
}

func Time(name string, input time.Time) Property {
	return Property{updated: true, Name: name, Time: &input, Type: PropertyTypeTime}
}

func Int(name string, input int64) Property {
	return Property{updated: true, Name: name, Int: input, Type: PropertyTypeInt}
}

func Bool(name string, input bool) Property {
	return Property{updated: true, Name: name, Bool: input, Type: PropertyTypeBool}
}

func Float(name string, input float64) Property {
	return Property{updated: true, Name: name, Float: input, Type: PropertyTypeFloat}
}

func (p *Property) Updated() {
	p.updated = true
}

func (p *Property) Value() interface{} {
	switch p.Type {
	case PropertyTypeText:
		return p.Text
	case PropertyTypeInt:
		return p.Int
	case PropertyTypeBool:
		return p.Bool
	case PropertyTypeFloat:
		return p.Float
	case PropertyTypeTime:
		return p.Time
	}
	return nil
}

func (p *Property) AsPersonal()  { p.Classification = ClassificationPersonal }
func (p *Property) AsUserInput() { p.Classification = ClassificationUserInput }
func (p *Property) AsSecure()    { p.Classification = ClassificationSecure }
func (p *Property) AsIndexed()   { p.Classification = ClassificationIndexed }
func (p *Property) AsID()        { p.Classification = ClassificationID }
func (p *Property) AsAnonymous() { p.Classification = ClassificationAnonymous }

func (p *Property) GetText() string {
	if p != nil {
		return p.Text
	}
	return ""
}

func (p *Property) GetSecureText() string {
	if p != nil {
		return p.Secret
	}
	return ""
}

func (p *Property) GetInt() int64 {
	if p != nil {
		return p.Int
	}
	return 0
}

func (p *Property) GetBool() bool {
	if p != nil {
		return p.Bool
	}
	return false
}

func (p *Property) GetFloat() float64 {
	if p != nil {
		return p.Float
	}
	return 0.0
}

func (p *Property) GetTime() *time.Time {
	if p != nil {
		return p.Time
	}
	return nil
}

func PersonalData(p Property) Property  { p.AsPersonal(); return p }
func UserInput(p Property) Property     { p.AsUserInput(); return p }
func SecureData(p Property) Property    { p.AsSecure(); return p }
func Indexed(p Property) Property       { p.AsIndexed(); return p }
func Identifier(p Property) Property    { p.AsID(); return p }
func AnonymousData(p Property) Property { p.AsAnonymous(); return p }

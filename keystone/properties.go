package keystone

import (
	"github.com/kubex/definitions-go/app"
	"github.com/kubex/keystone-go/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Property struct {
	Name           app.ScopedKey
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
	return &proto.Property{
		Key:            toKey(p.Name),
		Type:           p.Type.toProto(),
		Classification: p.Classification.toProto(),
		Text:           p.Text,
		Int:            p.Int,
		Bool:           p.Bool,
		Float:          float32(p.Float),
		Time:           timestamppb.New(*p.Time),
		SecureText:     p.Secret,
	}
}

func Encrypted(name, decrypted, preview string) Property {
	return Property{updated: true, Name: app.NewScopedKey(name, defaultSetGlobalAppID), Secret: decrypted, Text: preview, Classification: ClassificationSecure, Type: PropertyTypeText}
}

func Text(name, input string) Property {
	return Property{updated: true, Name: app.NewScopedKey(name, defaultSetGlobalAppID), Text: input, Type: PropertyTypeText}
}

func Time(name string, input time.Time) Property {
	return Property{updated: true, Name: app.NewScopedKey(name, defaultSetGlobalAppID), Time: &input, Type: PropertyTypeTime}
}

func Int(name string, input int64) Property {
	return Property{updated: true, Name: app.NewScopedKey(name, defaultSetGlobalAppID), Int: input, Type: PropertyTypeInt}
}

func Bool(name string, input bool) Property {
	return Property{updated: true, Name: app.NewScopedKey(name, defaultSetGlobalAppID), Bool: input, Type: PropertyTypeBool}
}

func Float(name string, input float64) Property {
	return Property{updated: true, Name: app.NewScopedKey(name, defaultSetGlobalAppID), Float: input, Type: PropertyTypeFloat}
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

func PersonalData(p Property) Property  { p.AsPersonal(); return p }
func UserInput(p Property) Property     { p.AsUserInput(); return p }
func SecureData(p Property) Property    { p.AsSecure(); return p }
func IndexedData(p Property) Property   { p.AsIndexed(); return p }
func Identifier(p Property) Property    { p.AsID(); return p }
func AnonymousData(p Property) Property { p.AsAnonymous(); return p }

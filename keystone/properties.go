package keystone

import (
	"github.com/kubex/definitions-go/app"
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

func (v *Property) Updated() {
	v.updated = true
}

func (v *Property) Value() interface{} {
	switch v.Type {
	case PropertyTypeText:
		return v.Text
	case PropertyTypeInt:
		return v.Int
	case PropertyTypeBool:
		return v.Bool
	case PropertyTypeFloat:
		return v.Float
	case PropertyTypeTime:
		return v.Time
	}
	return nil
}

package keystone

import (
	"github.com/kubex/definitions-go/app"
	"github.com/kubex/keystone-go/proto"
)

type Classification string

const (
	ClassificationAnonymous Classification = "a"
	ClassificationPersonal  Classification = "p"
	ClassificationUserInput Classification = "u" // Unknown free-type input by a user
	ClassificationIndexed   Classification = "q" //Data that should be indexed for queries
	ClassificationID        Classification = "i" // Index in 1-X relationship
	ClassificationSecure    Classification = "s" // Not to be indexed, should be encrypted
)

func (c Classification) toProto() proto.DataClassification {
	switch c {
	case ClassificationAnonymous:
		return proto.DataClassification_CLASSIFICATION_ANONYMOUS
	case ClassificationPersonal:
		return proto.DataClassification_CLASSIFICATION_PERSONAL
	case ClassificationUserInput:
		return proto.DataClassification_CLASSIFICATION_USER_INPUT
	case ClassificationIndexed:
		return proto.DataClassification_CLASSIFICATION_INDEXED
	case ClassificationID:
		return proto.DataClassification_CLASSIFICATION_ID
	case ClassificationSecure:
		return proto.DataClassification_CLASSIFICATION_SECURE
	}
	return proto.DataClassification_CLASSIFICATION_ANONYMOUS
}

func GetClassification(classification proto.DataClassification) Classification {
	switch classification {
	case proto.DataClassification_CLASSIFICATION_ANONYMOUS:
		return ClassificationAnonymous
	case proto.DataClassification_CLASSIFICATION_PERSONAL:
		return ClassificationPersonal
	case proto.DataClassification_CLASSIFICATION_USER_INPUT:
		return ClassificationUserInput
	case proto.DataClassification_CLASSIFICATION_INDEXED:
		return ClassificationIndexed
	case proto.DataClassification_CLASSIFICATION_ID:
		return ClassificationID
	case proto.DataClassification_CLASSIFICATION_SECURE:
		return ClassificationSecure
	}
	return ClassificationAnonymous
}

type PropertyType string

const (
	PropertyTypeText  = "s"
	PropertyTypeInt   = "i"
	PropertyTypeBool  = "b"
	PropertyTypeFloat = "f"
	PropertyTypeTime  = "t"
)

func (p PropertyType) toProto() proto.ValueType {
	switch p {
	case PropertyTypeText:
		return proto.ValueType_VALUE_TEXT
	case PropertyTypeInt:
		return proto.ValueType_VALUE_NUMBER
	case PropertyTypeBool:
		return proto.ValueType_VALUE_BOOLEAN
	case PropertyTypeFloat:
		return proto.ValueType_VALUE_FLOAT
	case PropertyTypeTime:
		return proto.ValueType_VALUE_TIME
	}
	return proto.ValueType_VALUE_TEXT
}

func GetPropertyType(valueType proto.ValueType) PropertyType {
	switch valueType {
	case proto.ValueType_VALUE_TEXT:
		return PropertyTypeText
	case proto.ValueType_VALUE_NUMBER:
		return PropertyTypeInt
	case proto.ValueType_VALUE_BOOLEAN:
		return PropertyTypeBool
	case proto.ValueType_VALUE_FLOAT:
		return PropertyTypeFloat
	case proto.ValueType_VALUE_TIME:
		return PropertyTypeTime
	}
	return PropertyTypeText
}

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelNotice
	LogLevelWarn
	LogLevelError
	LogLevelCritical
	LogLevelAlert
	LogLevelFatal
)

func (l LogLevel) toProto() proto.LogLevel {
	switch l {
	case LogLevelDebug:
		return proto.LogLevel_LevelDebug
	case LogLevelInfo:
		return proto.LogLevel_LevelInfo
	case LogLevelNotice:
		return proto.LogLevel_LevelNotice
	case LogLevelWarn:
		return proto.LogLevel_LevelWarn
	case LogLevelError:
		return proto.LogLevel_LevelError
	case LogLevelCritical:
		return proto.LogLevel_LevelCritical
	case LogLevelAlert:
		return proto.LogLevel_LevelAlert
	case LogLevelFatal:
		return proto.LogLevel_LevelFatal
	}
	return proto.LogLevel_LevelInfo
}

func GetLogLevel(logLevel proto.LogLevel) LogLevel {
	switch logLevel {
	case proto.LogLevel_LevelDebug:
		return LogLevelDebug
	case proto.LogLevel_LevelInfo:
		return LogLevelInfo
	case proto.LogLevel_LevelNotice:
		return LogLevelNotice
	case proto.LogLevel_LevelWarn:
		return LogLevelWarn
	case proto.LogLevel_LevelError:
		return LogLevelError
	case proto.LogLevel_LevelCritical:
		return LogLevelCritical
	case proto.LogLevel_LevelAlert:
		return LogLevelAlert
	case proto.LogLevel_LevelFatal:
		return LogLevelFatal
	}
	return LogLevelInfo
}

func GetScopedKey(key *proto.Key) app.ScopedKey {
	gaid := app.NewID(key.GetVendorId(), key.GetAppId())
	return app.NewScopedKey(key.GetKey(), &gaid)
}

func toKey(scopedKey app.ScopedKey) *proto.Key {
	return &proto.Key{
		VendorId: scopedKey.VendorID,
		AppId:    scopedKey.AppID,
		Key:      scopedKey.Key,
	}
}

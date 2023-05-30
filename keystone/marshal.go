package keystone

import (
	"errors"
	"reflect"
	"time"
)

func Marshal(src interface{}, dst *Entity) error {
	if src == nil {
		return errors.New("Cannot marshall nil")
	}
	if dst.Properties == nil {
		dst.Properties = make(map[string]Property)
	}

	if logs, ok := src.(EntityLogProvider); ok {
		logEntries, logErr := logs.GetLogs()
		if logErr != nil {
			return logErr
		}
		for _, logEntry := range logEntries {
			dst.LogEntries = append(dst.LogEntries, logEntry)
		}
		_ = logs.ClearLogs()
	}

	if events, ok := src.(EntityEventProvider); ok {
		eventEntries, eventErr := events.GetEvents()
		if eventErr != nil {
			return eventErr
		}
		for _, eventEntry := range eventEntries {
			dst.Events = append(dst.Events, eventEntry)
		}
		_ = events.ClearEvents()
	}

	v := reflect.ValueOf(src)
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		fName := getFieldName(field)
		if fName == "" {
			continue
		}

		switch field.Type.Kind() {
		case reflect.String:
			dst.Properties[fName] = Text(fName, v.Field(i).String())
			continue
		case reflect.Bool:
			dst.Properties[fName] = Bool(fName, v.Field(i).Bool())
			continue
		case reflect.Int32, reflect.Int64:
			dst.Properties[fName] = Int(fName, v.Field(i).Int())
			continue
		case reflect.Float32, reflect.Float64:
			dst.Properties[fName] = Float(fName, v.Field(i).Float())
			continue
		}

		switch field.Type {
		case typeOfSecretString:
			if val, ok := v.Field(i).Interface().(SecretString); ok {
				dst.Properties[fName] = Secret(fName, val.Original, val.Masked)
			}
			continue
		case typeOfTime:
			if val, ok := v.Field(i).Interface().(time.Time); ok {
				dst.Properties[fName] = Time(fName, val)
			}
			continue
		}

	}

	return nil
}

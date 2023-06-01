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
		fOpt := getFieldOptions(field)
		if fOpt.name == "" {
			continue
		}
		fName := fOpt.name

		switch field.Type.Kind() {
		case reflect.String:
			dst.addProperty(fName, prepareProperty(Text(fName, v.Field(i).String()), fOpt))
			continue
		case reflect.Bool:
			dst.addProperty(fName, prepareProperty(Bool(fName, v.Field(i).Bool()), fOpt))
			continue
		case reflect.Int32, reflect.Int64:
			dst.addProperty(fName, prepareProperty(Int(fName, v.Field(i).Int()), fOpt))
			continue
		case reflect.Float32, reflect.Float64:
			dst.addProperty(fName, prepareProperty(Float(fName, v.Field(i).Float()), fOpt))
			continue
		}

		switch field.Type {
		case typeOfSecretString:
			if val, ok := v.Field(i).Interface().(SecretString); ok {
				dst.addProperty(fName, prepareProperty(Secret(fName, val.Original, val.Masked), fOpt))
			}
			continue
		case typeOfAmount:
			if val, ok := v.Field(i).Interface().(Amount); ok {
				dst.addProperty(fName, prepareProperty(Money(fName, val.Currency, val.Units), fOpt))
			}
			continue
		case typeOfTime:
			if val, ok := v.Field(i).Interface().(time.Time); ok {
				dst.addProperty(fName, prepareProperty(Time(fName, val), fOpt))
			}
			continue
		}

	}

	return nil
}

func (e *Entity) addProperty(name string, prop Property) {
	if prop.omitempty && prop.IsEmpty() {
		return
	}
	e.Properties[name] = prop
}

func prepareProperty(property Property, options fieldOptions) Property {
	if options.indexed {
		property.AsIndexed()
	}
	if options.lookup {
		property.AsLookup()
	}
	if options.omitempty {
		property.omitempty = true
	}
	return property
}

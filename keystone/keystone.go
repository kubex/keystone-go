// Package keystone is a database abstraction layer
package keystone

import (
	"reflect"
	"time"
)

// GenericResult is a map that can be used to retrieve a generic result
type GenericResult map[string]interface{}

// EntityChild is an interface that defines a child entity
type EntityChild interface {
	ChildID() string
}

// SecretString is a string that represents sensitive data
type SecretString struct {
	Masked   string `json:"masked,omitempty"`
	Original string `json:"original,omitempty"`
}

// Amount represents money
type Amount struct {
	Currency string `json:"currency"`
	Units    int64  `json:"units"`
}

var (
	typeOfTime         = reflect.TypeOf(time.Time{})
	typeOfSecretString = reflect.TypeOf(SecretString{})
	typeOfAmount       = reflect.TypeOf(Amount{})
	typeOfStringSlice  = reflect.TypeOf([]string{})
)

// String returns the original string if it exists, otherwise the masked string
func (e SecretString) String() string {
	if e.Original != "" {
		return e.Original
	}
	return e.Masked
}

// NewSecretString creates a new SecretString
func NewSecretString(original, masked string) SecretString {
	return SecretString{
		Masked:   masked,
		Original: original,
	}
}

// NewAmount creates a new Amount
func NewAmount(currency string, units int64) Amount {
	return Amount{
		Currency: currency,
		Units:    units,
	}
}

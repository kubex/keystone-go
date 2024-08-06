package keystone

import (
	"github.com/kubex/keystone-go/proto"
)

// Amount represents money
type Amount struct {
	Currency string `json:"currency"`
	Units    int64  `json:"units"`
}

// NewAmount creates a new Amount
func NewAmount(currency string, units int64) Amount {
	return Amount{
		Currency: currency,
		Units:    units,
	}
}

func (a *Amount) FromProtoValue(value *proto.Value) error {
	if value != nil {
		a.Units = value.GetInt()
		a.Currency = value.GetText()
	}
	return nil
}

func (a *Amount) ToProtoValue() *proto.Value {
	return &proto.Value{
		Text: a.Currency,
		Int:  a.Units,
	}
}

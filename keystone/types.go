package keystone

import "github.com/kubex/keystone-go/proto"

type CustomType interface {
	FromProtoValue(value *proto.Value) error
	ToProtoValue() *proto.Value
}

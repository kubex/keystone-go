package keystone

import (
	"github.com/kubex/keystone-go/proto"
)

type FindOption interface {
	Apply(config *proto.FindRequest)
}

func WhereKeyEquals(key, value string) FindOption {
	return propertyFilter{key: key, values: []string{value}, operator: proto.Operator_Equal}
}

func WhereKeyNotEquals(key, value string) FindOption {
	return propertyFilter{key: key, values: []string{value}, operator: proto.Operator_NotEqual}
}

func WhereKeyGreaterThan(key, value string) FindOption {
	return propertyFilter{key: key, values: []string{value}, operator: proto.Operator_GreaterThan}
}

func WhereKeyGreaterThanOrEquals(key, value string) FindOption {
	return propertyFilter{key: key, values: []string{value}, operator: proto.Operator_GreaterThanOrEqual}
}

func WhereKeyLessThan(key, value string) FindOption {
	return propertyFilter{key: key, values: []string{value}, operator: proto.Operator_LessThan}
}

func WhereKeyLessThanOrEquals(key, value string) FindOption {
	return propertyFilter{key: key, values: []string{value}, operator: proto.Operator_LessThanOrEqual}
}

func WhereKeyContains(key, value string) FindOption {
	return propertyFilter{key: key, values: []string{value}, operator: proto.Operator_Contains}
}

func WhereKeyNotContains(key, value string) FindOption {
	return propertyFilter{key: key, values: []string{value}, operator: proto.Operator_NotContains}
}

func WhereKeyStartsWith(key, value string) FindOption {
	return propertyFilter{key: key, values: []string{value}, operator: proto.Operator_StartsWith}
}

func WhereKeyEndsWith(key, value string) FindOption {
	return propertyFilter{key: key, values: []string{value}, operator: proto.Operator_EndsWith}
}

func WhereKeyIn(key string, value ...string) FindOption {
	return propertyFilter{key: key, values: value, operator: proto.Operator_In}
}

func WhereKeyNotIn(key, value string) FindOption {
	return propertyFilter{key: key, values: []string{value}, operator: proto.Operator_Equal}
}

func WhereKeyBetween(key, value1, value2 string) FindOption {
	return propertyFilter{key: key, values: []string{value1, value2}, operator: proto.Operator_Equal}
}

type propertyFilter struct {
	key      string
	values   []string
	operator proto.Operator
}

func (f propertyFilter) Apply(config *proto.FindRequest) {
	if config.Filters == nil {
		config.Filters = make([]*proto.PropertyFilter, 0)
	}

	var vals []*proto.Value
	for _, v := range f.values {
		vals = append(vals, &proto.Value{Text: v})
	}

	config.Filters = append(config.Filters, &proto.PropertyFilter{
		Property: &proto.Key{Key: f.key},
		Operator: f.operator,
		Values:   vals,
	})
}

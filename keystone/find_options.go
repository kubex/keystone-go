package keystone

import (
	"github.com/kubex/keystone-go/proto"
)

type FindOption interface {
	Apply(config *proto.FindRequest)
}

func WhereEquals(key string, value any) FindOption {
	return propertyFilter{key: key, values: valuesFromAny(value), operator: proto.Operator_Equal}
}

func WhereNotEquals(key string, value any) FindOption {
	return propertyFilter{key: key, values: valuesFromAny(value), operator: proto.Operator_NotEqual}
}

func WhereGreaterThan(key string, value any) FindOption {
	return propertyFilter{key: key, values: valuesFromAny(value), operator: proto.Operator_GreaterThan}
}

func WhereGreaterThanOrEquals(key string, value any) FindOption {
	return propertyFilter{key: key, values: valuesFromAny(value), operator: proto.Operator_GreaterThanOrEqual}
}

func WhereLessThan(key string, value any) FindOption {
	return propertyFilter{key: key, values: valuesFromAny(value), operator: proto.Operator_LessThan}
}

func WhereLessThanOrEquals(key string, value any) FindOption {
	return propertyFilter{key: key, values: valuesFromAny(value), operator: proto.Operator_LessThanOrEqual}
}

func WhereContains(key string, value any) FindOption {
	return propertyFilter{key: key, values: valuesFromAny(value), operator: proto.Operator_Contains}
}

func WhereNotContains(key string, value any) FindOption {
	return propertyFilter{key: key, values: valuesFromAny(value), operator: proto.Operator_NotContains}
}

func WhereStartsWith(key string, value any) FindOption {
	return propertyFilter{key: key, values: valuesFromAny(value), operator: proto.Operator_StartsWith}
}

func WhereEndsWith(key string, value string) FindOption {
	return propertyFilter{key: key, values: valuesFromAny(value), operator: proto.Operator_EndsWith}
}

func WhereIn(key string, value ...string) FindOption {
	return propertyFilter{key: key, values: valuesFromAny(value), operator: proto.Operator_In}
}

func WhereNotIn(key string, value string) FindOption {
	return propertyFilter{key: key, values: valuesFromAny(value), operator: proto.Operator_Equal}
}

func WhereBetween(key string, value1, value2 any) FindOption {
	return propertyFilter{key: key, values: valuesFromAny(value1, value2), operator: proto.Operator_Equal}
}

type propertyFilter struct {
	key      string
	values   []*proto.Value
	operator proto.Operator
}

func (f propertyFilter) Apply(config *proto.FindRequest) {
	if config.Filters == nil {
		config.Filters = make([]*proto.PropertyFilter, 0)
	}

	config.Filters = append(config.Filters, &proto.PropertyFilter{
		Property: &proto.Key{Key: f.key},
		Operator: f.operator,
		Values:   f.values,
	})
}

func valueFromAny(value any) *proto.Value {
	switch v := value.(type) {
	case string:
		return &proto.Value{Text: v}
	case int, int32, int64:
		return &proto.Value{Int: int64(v.(int))}
	case bool:
		return &proto.Value{Bool: v}
	case float64:
		return &proto.Value{Float: v}
	}
	return &proto.Value{}
}

func valuesFromAny(values ...any) []*proto.Value {
	var result []*proto.Value
	for _, v := range values {
		result = append(result, valueFromAny(v))
	}
	return result
}

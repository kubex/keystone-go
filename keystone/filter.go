package keystone

import "github.com/kubex/keystone-go/proto"

type FilterOperator string

const FilterEqual FilterOperator = "eq"
const FilterNotEqual FilterOperator = "neq"
const FilterGreaterThan FilterOperator = "gt"
const FilterGreaterThanEqual FilterOperator = "gte"
const FilterLessThan FilterOperator = "lt"
const FilterLessThanEqual FilterOperator = "lte"
const FilterLike FilterOperator = "like"
const FilterIn FilterOperator = "in"

type Option struct {
	Filters []*proto.FilterProperty
}

func NewFilterOption(filters ...*proto.FilterProperty) Option {
	return Option{
		Filters: filters,
	}
}

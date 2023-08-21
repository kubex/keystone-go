package keystone

type sortBy struct {
	property  string
	direction string
}

func (f sortBy) Apply(config *filterRequest) {
	config.SortProperty = f.property
	config.SortDirection = f.direction
}

func SortBy(property string, direction string) FindOption {
	return sortBy{property: property, direction: direction}
}

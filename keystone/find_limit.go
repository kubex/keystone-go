package keystone

type withLimit struct {
	limit  int32
	offset int32
}

func (f withLimit) Apply(config *filterRequest) {
	config.Limit = f.limit
	config.Offset = f.offset
}

func Limit(limit, offset int32) FindOption {
	return withLimit{limit: limit, offset: offset}
}

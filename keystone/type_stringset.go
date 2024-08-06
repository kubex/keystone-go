package keystone

type StringSet struct {
	values map[string]bool

	toAdd           map[string]bool
	toRemove        map[string]bool
	replaceExisting bool
}

func (s *StringSet) Clear() {
	s.values = make(map[string]bool)
	s.toAdd = make(map[string]bool)
	s.toRemove = make(map[string]bool)
}

func (s *StringSet) Add(value string) {
	s.toAdd[value] = true
	delete(s.toRemove, value)
}

func (s *StringSet) Remove(value string) {
	s.toRemove[value] = true
	delete(s.toAdd, value)
}

func (s *StringSet) Values() []string {
	var values []string
	for value := range s.values {
		values = append(values, value)
	}
	return values
}

func (s *StringSet) Has(value string) bool {
	_, ok := s.values[value]
	return ok
}

func (s *StringSet) ReplaceWith(values ...string) {
	s.Clear()
	s.replaceExisting = true
	s.applyValues(values...)
}

func (s *StringSet) applyValues(values ...string) {
	for _, value := range values {
		s.values[value] = true
	}
}

func (s *StringSet) IsEmpty() bool {
	return len(s.values) == 0
}

func (s *StringSet) ToAdd() []string {
	var values []string
	for value := range s.toAdd {
		values = append(values, value)
	}
	return values
}

func (s *StringSet) ToRemove() []string {
	var values []string
	for value := range s.toRemove {
		values = append(values, value)
	}
	return values
}

func (s *StringSet) ReplaceExisting() bool {
	return s.replaceExisting
}

func (s *StringSet) Diff(values ...string) []string {
	check := make(map[string]bool, len(values))
	for _, x := range values {
		check[x] = s.Has(x)
	}
	var diff []string
	for x := range s.values {
		if _, ok := check[x]; !ok {
			diff = append(diff, x)
		}
	}
	for x, matched := range check {
		if !matched {
			diff = append(diff, x)
		}
	}
	return diff
}

func NewStringSet(values ...string) StringSet {
	v := StringSet{}
	v.Clear()
	v.applyValues(values...)

	return v
}

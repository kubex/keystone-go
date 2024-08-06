package keystone

type IntSet struct {
	values map[int64]bool

	toAdd           map[int64]bool
	toRemove        map[int64]bool
	replaceExisting bool
}

func (s *IntSet) Clear() {
	s.values = make(map[int64]bool)
	s.toAdd = make(map[int64]bool)
	s.toRemove = make(map[int64]bool)
}

func (s *IntSet) Add(value int64) {
	s.toAdd[value] = true
	delete(s.toRemove, value)
}

func (s *IntSet) Remove(value int64) {
	s.toRemove[value] = true
	delete(s.toAdd, value)
}

func (s *IntSet) Values() []int64 {
	var values []int64
	for value := range s.values {
		values = append(values, value)
	}
	return values
}

func (s *IntSet) Has(value int64) bool {
	_, ok := s.values[value]
	return ok
}

func (s *IntSet) ReplaceWith(values ...int64) {
	s.Clear()
	s.replaceExisting = true
	s.applyValues(values...)
}

func (s *IntSet) applyValues(values ...int64) {
	for _, value := range values {
		s.values[value] = true
	}
}

func (s *IntSet) IsEmpty() bool {
	return len(s.values) == 0
}

func (s *IntSet) ToAdd() []int64 {
	var values []int64
	for value := range s.toAdd {
		values = append(values, value)
	}
	return values
}

func (s *IntSet) ToRemove() []int64 {
	var values []int64
	for value := range s.toRemove {
		values = append(values, value)
	}
	return values
}

func (s *IntSet) ReplaceExisting() bool {
	return s.replaceExisting
}

func (s *IntSet) Diff(values ...int64) []int64 {
	check := make(map[int64]bool, len(values))
	for _, x := range values {
		check[x] = s.Has(x)
	}
	var diff []int64
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

func NewIntSet(values ...int64) IntSet {
	v := IntSet{}
	v.Clear()
	v.applyValues(values...)
	return v
}

package keystone

type IntSet struct {
	values map[int]bool

	toAdd           map[int]bool
	toRemove        map[int]bool
	replaceExisting bool
}

func (s *IntSet) Clear() {
	s.values = make(map[int]bool)
	s.toAdd = make(map[int]bool)
	s.toRemove = make(map[int]bool)
}

func (s *IntSet) Add(value int) {
	s.toAdd[value] = true
	delete(s.toRemove, value)
}

func (s *IntSet) Remove(value int) {
	s.toRemove[value] = true
	delete(s.toAdd, value)
}

func (s *IntSet) Values() []int {
	var values []int
	for value := range s.values {
		values = append(values, value)
	}
	return values
}

func (s *IntSet) Has(value int) bool {
	_, ok := s.values[value]
	return ok
}

func (s *IntSet) ReplaceWith(values ...int) {
	s.Clear()
	s.replaceExisting = true
	s.applyValues(values...)
}

func (s *IntSet) applyValues(values ...int) {
	for _, value := range values {
		s.values[value] = true
	}
}

func (s *IntSet) IsEmpty() bool {
	return len(s.values) == 0
}

func (s *IntSet) ToAdd() []int {
	var values []int
	for value := range s.toAdd {
		values = append(values, value)
	}
	return values
}

func (s *IntSet) ToRemove() []int {
	var values []int
	for value := range s.toRemove {
		values = append(values, value)
	}
	return values
}

func (s *IntSet) ReplaceExisting() bool {
	return s.replaceExisting
}

func NewIntSet(values ...int) IntSet {
	v := IntSet{}
	v.Clear()
	v.applyValues(values...)
	return v
}

package proto

type EntityResponseIDSort []*EntityResponse

func (a EntityResponseIDSort) Len() int      { return len(a) }
func (a EntityResponseIDSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a EntityResponseIDSort) Less(i, j int) bool {
	return a[i].GetEntity().GetEntityId() < a[j].GetEntity().GetEntityId()
}

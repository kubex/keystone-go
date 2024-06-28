package keystone

import (
	"github.com/kubex/keystone-go/proto"
	"testing"
)

func TestEntityResponseToDst(t *testing.T) {

	dst := &testBaseEntity{}

	children := []*proto.EntityChild{
		{
			Type:  &proto.Key{Key: "children"},
			Cid:   "cid1",
			Value: 1,
			Data: map[string][]byte{
				"ChildName": []byte(`"child1"`),
				"Data":      []byte(`"data1"`),
			},
		},
		{
			Type:  &proto.Key{Key: "children"},
			Cid:   "cid2",
			Value: 2,
			Data: map[string][]byte{
				"ChildName": []byte(`"child2"`),
				"Data":      []byte(`"data2"`),
			},
		},
	}

	err := entityResponseToDst(nil, children, dst, "")
	if err != nil {
		t.Error(err)
	}

	if len(dst.Children) != 2 {
		t.Error("Expected 2 Children, got", len(dst.Children))
	} else if dst.Children[0].ChildName != "child1" {
		t.Error("Expected child1, got", dst.Children[0].ChildName)
	} else if dst.Children[1].ChildName != "child2" {
		t.Error("Expected child2, got", dst.Children[1].ChildName)
	}

}

type testBaseEntity struct {
	BaseEntity
	BasicProp string
	Children  []*testBaseChildEntity
}

type testBaseChildEntity struct {
	BaseNestedChild
	Data      []byte
	ChildName string
}

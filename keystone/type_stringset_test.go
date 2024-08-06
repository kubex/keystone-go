package keystone

import (
	"log"
	"testing"
)

func TestStringSet_Diff(t *testing.T) {
	orig := NewStringSet("a", "b")
	log.Println(orig.Diff("b", "c"))
}

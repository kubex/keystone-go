package keystone

import (
	"github.com/kubex/keystone-go/proto"
	"log"
	"testing"
)

// Test the function Unmarshal
func TestUnmarshal(t *testing.T) {
	s := unmarshalTest{}
	log.Println(Unmarshal(&proto.EntityResponse{
		EntityId: "random-uuid-1234",
		Properties: []*proto.Property{
			{Name: "id", Text: "abslfdwuflwkdh"},
			{Name: "address_country_code", Text: "US"},
		},
	}, &s))
	log.Println(s)
}

type unmarshalTest struct {
	ID                 string `keystone:"_entity_id"`
	AddressCountryCode string `keystone:""`
}

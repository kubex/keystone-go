package keystone

import (
	"github.com/kubex/keystone-go/proto"
	"log"
	"testing"
	"time"
)

// Test the function Unmarshal
func TestUnmarshal(t *testing.T) {
	s := unmarshalTest{}
	log.Println(Unmarshal(&proto.EntityResponse{
		EntityId: "random-uuid-1234",
		Properties: []*proto.Property{
			{Name: "id", Value: &proto.Value{Text: "abslfdwuflwkdh"}},
			{Name: "address_country_code", Value: &proto.Value{Text: "US"}},
		},
	}, &s))
	log.Println(s)
}

type unmarshalTest struct {
	baseStruct
	AddressCountryCode string       `keystone:""`
	EmailAddress       SecretString `keystone:",indexed"`
	second             secondStruct `keystone:",indexed,lookup,omitempty"`
}

type baseStruct struct {
	ID string `keystone:"_entity_id"`
}

type secondStruct struct {
	DateCreated time.Time `keystone:"_created"`
}

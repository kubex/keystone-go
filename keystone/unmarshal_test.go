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
			{Name: "id", Text: "abslfdwuflwkdh"},
			{Name: "address_country_code", Text: "US"},
		},
	}, &s))
	log.Println(s)
}

type unmarshalTest struct {
	baseStruct
	AddressCountryCode string `keystone:""`
	EmailAddress       SecretString
	second             secondStruct
}

type baseStruct struct {
	ID string `keystone:"_entity_id"`
}

type secondStruct struct {
	DateCreated time.Time `keystone:"_created"`
}

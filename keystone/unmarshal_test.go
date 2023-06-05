package keystone

import (
	"github.com/kubex/keystone-go/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
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
			{Name: "name", Value: &proto.Value{Text: "Putin"}},
			{Name: "created", Value: &proto.Value{Time: timestamppb.New(time.Now())}},
		},
	}, &s))
	log.Println(s.ID)
	log.Println(s.AddressCountryCode)
	log.Println(s.Second)
}

type unmarshalTest struct {
	BaseStruct
	AddressCountryCode string       `keystone:""`
	EmailAddress       SecretString `keystone:",indexed"`
	Second             SecondStruct `keystone:",indexed,lookup,omitempty"`
}

type BaseStruct struct {
	ID string `keystone:"id"`
}

type SecondStruct struct {
	ID          string `keystone:"id"`
	Name        string
	DateCreated time.Time `keystone:"created"`
}

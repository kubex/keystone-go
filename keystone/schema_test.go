package keystone

import (
	"encoding/json"
	"github.com/kubex/keystone-go/proto"
	"log"
	"testing"
	"time"
)

func TestTypeToSchema(t *testing.T) {
	by, e := json.Marshal(testSchemaType{})
	log.Println(string(by), e)
	result := typeToSchema(testSchemaType{})
	log.Println(result)
}

// testSchemaType is a test type for testing the schema generation
type testSchemaType struct {
	testSchemaTypeBase
	AddressCountryCode string               `keystone:""`
	EmailAddress       SecretString         `keystone:",indexed,pii,lookup"`
	Second             testSchemaTypeNested `keystone:",indexed,lookup,omitempty"`
}

func (t testSchemaType) GetKeystoneDefinition() TypeDefinition {
	return TypeDefinition{
		Name:        "Test Schema Type",
		Description: "A test schema type for testing the schema generation",
		Options:     []proto.Schema_Option{proto.Schema_StoreMutations},
		Singular:    "Customer",
		Plural:      "Customers",
	}
}

type testSchemaTypeBase struct {
	ID string `keystone:"id,unique"`
}

type testSchemaTypeNested struct {
	ID          string `keystone:"id"`
	Name        string
	Person      testSchemaTypeDoubleNested `keystone:"person"`
	DateCreated time.Time                  `keystone:"created"`
}

type testSchemaTypeDoubleNested struct {
	FirstName string `keystone:",pii"`
	LastName  string `keystone:",pii"`
	TotalPaid Amount
}

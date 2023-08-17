package keystone

import (
	"reflect"
	"testing"

	"github.com/kubex/keystone-go/proto"
)

func TestUnmarshal(t *testing.T) {
	addr := &Address{}
	if err := unmarshal(testAddressResponse, addr); err != nil {
		t.Error(err)
	}
	checkAddress(t, addr)
}

func TestUnmarshalGeneric(t *testing.T) {
	gr := GenericResult{}
	if err := unmarshalGeneric(testAddressResponse, gr); err != nil {
		t.Error(err)
	}
	if gr["_entity_id"] != "xxx-xxxx" {
		t.Error("_entity_id not set")
	}
	if gr["line1"] != "123 Fake St." {
		t.Error("line1 not set")
	}
	if gr["line2"] != "Apt. 2" {
		t.Error("line2 not set")
	}
	if gr["city"] != "Springfield" {
		t.Error("city not set")
	}
}

func TestMakeEntityPropertyMap(t *testing.T) {
	result := makeEntityPropertyMap(testAddressResponse)
	if len(result) != 4 {
		t.Error("wrong number of properties")
	}
	if result[EntityIDKey].Value.Text != "xxx-xxxx" {
		t.Error("_entity_id not set")
	}
	if result["line1"].Value.Text != "123 Fake St." {
		t.Error("line1 not set")
	}
	if result["line2"].Value.Text != "Apt. 2" {
		t.Error("line2 not set")
	}
	if result["city"].Value.Text != "Springfield" {
		t.Error("city not set")
	}
}

func TestEntityResponseToDst(t *testing.T) {
	addr := &Address{}
	if err := entityResponseToDst(makeEntityPropertyMap(testAddressResponse), nil, addr, ""); err != nil {
		t.Error(err)
	}
	checkAddress(t, addr)
}

func TestSetFieldValue(t *testing.T) {
	addr := &Address{}
	addressValue := reflect.ValueOf(addr)
	for addressValue.Kind() == reflect.Pointer {
		addressValue = addressValue.Elem()
	}
	propMap := map[string]*proto.EntityProperty{
		"line1": {Value: &proto.Value{Text: "123 Fake St."}},
		"line2": {Value: &proto.Value{Text: "Apt. 2"}},
		"city":  {Value: &proto.Value{Text: "Springfield"}},
	}
	for i := 0; i < addressValue.NumField(); i++ {
		field := addressValue.Type().Field(i)
		fieldValue := addressValue.Field(i)
		fieldOpt := getFieldOptions(field, "")
		setFieldValue(field, fieldValue, fieldOpt, propMap)
	}
	checkAddress(t, addr)
}

func checkAddress(t *testing.T, addr *Address) {
	if addr.Line1 != "123 Fake St." {
		t.Errorf("Expected addr.Line1 to be '123 Fake St.', got '%s'", addr.Line1)
	}
	if addr.Line2 != "Apt. 2" {
		t.Errorf("Expected addr.Line2 to be 'Apt. 2', got '%s'", addr.Line2)
	}
	if addr.City != "Springfield" {
		t.Errorf("Expected addr.City to be 'Springfield', got '%s'", addr.City)
	}
}

package keystone

import (
	"reflect"
	"testing"

	"github.com/kubex/keystone-go/proto"
)

func TestUnmarshal(t *testing.T) {
	addr := &Address{}
	if err := Unmarshal(testAddressResponse, addr); err != nil {
		t.Error(err)
	}
	checkAddress(t, addr)
}

func TestUnmarshalAppendPointer(t *testing.T) {
	var addresses []*Address
	err := UnmarshalAppend(&addresses, testAddressResponse, testAddressResponse)
	if err != nil {
		t.Error(err)
	} else if addresses == nil {
		t.Error("addresses is nil")
	} else if len(addresses) != 2 {
		t.Error("addresses is not 2")
	} else {
		for _, a := range addresses {
			checkAddress(t, a)
		}
	}
}

func TestUnmarshalAppend(t *testing.T) {
	var addresses []Address
	err := UnmarshalAppend(&addresses, testAddressResponse, testAddressResponse)
	if err != nil {
		t.Error(err)
	} else if addresses == nil {
		t.Error("addresses is nil")
	} else if len(addresses) != 2 {
		t.Error("addresses is not 2")
	} else {
		for _, a := range addresses {
			checkAddress(t, &a)
		}
	}
}

func TestUnmarshalAppendE(t *testing.T) {
	var addresses []Address
	matchA := Address{Line1: "abc"}
	addresses = append(addresses, matchA)
	err := UnmarshalAppend(&addresses, testAddressResponse, testAddressResponse)
	if err != nil {
		t.Error(err)
	} else if addresses == nil {
		t.Error("addresses is nil")
	} else if len(addresses) != 3 {
		t.Error("addresses is not 2, its ", len(addresses))
	} else {
		for x, a := range addresses {
			if x == 0 {
				if reflect.DeepEqual(matchA, a) {
					t.Error("matchA not set")
				}
				continue
			}
			checkAddress(t, &a)
		}
	}
}

func TestUnmarshalGeneric(t *testing.T) {
	gr := GenericResult{}
	if err := UnmarshalGeneric(testAddressResponse, gr); err != nil {
		t.Error(err)
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

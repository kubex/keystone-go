package keystone

import (
	"context"
	"log"
	"testing"
)

func TestActorRetrieveByUnique(t *testing.T) {
	c := NewConnection(ksClient, "vendor", "appid", "accessToken")
	actor := c.Actor("test-workspace", "123.45.67.89", "user-1234", "User Agent Unknown")

	cst := &Customer{}

	if err := actor.Get(
		context.Background(),
		RetrieveByUnique{UniqueID: "user-233", Property: "user_id"},
		cst,
		WithProperties("address~"),
		WithDecryptedProperties("name", "email", "city", "state", "country", "postcode", "amount_paid", "lead_date", "user_id"),
		WithLabels(),
		WithSummary(),
		WithDatum(),
		WithChildren("line_items"),
	); err != nil {
		t.Error(err)
	}

	//actor.Marshal(cst, "testing actor")
	//
	log.Println(cst)
}

func TestActorRetrieveByEntityID(t *testing.T) {
	c := NewConnection(ksClient, "vendor", "appid", "accessToken")
	actor := c.Actor("test-workspace", "123.45.67.89", "user-1234", "User Agent Unknown")

	gr := GenericResult{}
	if err := actor.Get(
		context.Background(),
		RetrieveByEntityID{EntityID: "14nA6UwmK7zAYsxm"},
		gr,
		WithProperties("address~"),
		WithDecryptedProperties("name", "email", "city", "state", "country", "postcode", "amount_paid", "lead_date"),
		WithLabels(),
		WithSummary(),
		WithDatum(),
	); err != nil {
		t.Error(err)
	}
	log.Println(gr)
}

func TestActorFind(t *testing.T) {
	c := NewConnection(ksClient, "vendor", "appid", "accessToken")
	actor := c.Actor("test-workspace", "123.45.67.89", "user-1234", "User Agent Unknown")

	resp, err := actor.Find(
		context.Background(),
		"Customer",
		[]string{},
		WhereKeyEquals("name", "John Doe"),
	)

	log.Println(resp, err)
}

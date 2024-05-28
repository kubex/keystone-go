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
		ByUniqueProperty(cst, "user-233", "user_id"),
		cst,
		RetrieveOptions(
			WithProperties("address~"),
			WithDecryptedProperties("name", "email", "city", "state", "country", "postcode", "amount_paid", "lead_date", "user_id"),
			WithLabels(),
			WithSummary(),
			WithDatum(),
			WithChildren("line_items"),
			WithRelationships("user"),
		),
	); err != nil {
		t.Error(err)
	}

	//actor.Mutate(cst, "testing actor")
	//
	log.Println(cst)
	log.Println(cst.GetKeystoneRelationships())
}

func TestActorRetrieveByEntityID(t *testing.T) {
	c := NewConnection(ksClient, "vendor", "appid", "accessToken")
	actor := c.Actor("test-workspace", "123.45.67.89", "user-1234", "User Agent Unknown")

	gr := &Customer{}
	if err := actor.Get(
		context.Background(),
		ByEntityID(gr, "14nA6UwmK7zAYsxm"),
		gr,
		RetrieveOptions(
			WithProperties("address~"),
			WithDecryptedProperties("name", "email", "city", "state", "country", "postcode", "amount_paid", "lead_date"),
			WithLabels(),
			WithSummary(),
			WithDatum(),
		),
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
		nil,
		WhereEquals("name", "John Doe"),
	)

	t.Log(resp, err)
}

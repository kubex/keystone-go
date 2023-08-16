package keystone

import (
	"encoding/json"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/kubex/keystone-go/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestWrite(t *testing.T) {
	cust := &Customer{
		//ID:            "enzfUSpdK7z5JMpq",
		Name:          NewSecretString("John Doe", "J**n D*e"),
		Email:         NewSecretString("john.doe@gmail.com", "j*******@gma**.com"),
		Company:       "Chargehive Ltd",
		Phone:         "0791736u63434",
		City:          "Portsmouth",
		StreetName:    "New Street",
		StreetAddress: "41",
		Postcode:      "PO1 3AG",
		Timezone:      "BST",
		State:         "Hampshire",
		HasPaid:       true,
		Country:       "UK",
		CountryCode:   "GB",
		James:         "Eagle",
		AmountPaid:    NewAmount("USD", 123),
		LeadDate:      time.Now(),
		UserID:        "user-237",
		Address: Address{
			Line1: "123 Old Street",
			Line2: "Line 2 is optional",
			City:  "Southampton",
		},
		References: []string{"ref-1", "ref-2"},
		LineItems: []*LineItem{
			{Name: "foo"},
			{Name: "bar"},
		},
	}
	cust.AddKeystoneLabel("foo", "bar")
	cust.AddKeystoneLink("gcs", "logo", "https://storage.googleapis.com/keystone-assets/keystone-logo.png")
	cust.AddKeystoneRelationship("user", "customer", map[string]string{"foo": "bar"}, time.Now())
	if err := getTestActor(nil).Mutate(cust, "Creating Customer via Mutate"); err != nil {
		t.Error(err)
	}
}

func TestConnection(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			writeCustomers()
			wg.Done()
		}()
		time.Sleep(time.Millisecond * 100)
	}
	wg.Wait()
}

func writeCustomers() {
	log.Println("Marshalling")
	for i := 0; i < 10; i++ {
		getTestActor(nil).Mutate(FakeCustomer(), "Faker Customer x")
	}
	return
}

func xx(t *testing.T) {
	type address struct {
		Street string `json:"street"`
		City   string `json:"city"`
		State  string `json:"state"`
		Zip    string `json:"zip"`
	}

	addr1Bytes, _ := json.Marshal(address{
		Street: "E Rio Salado Pkwy",
		City:   "Tempe",
		State:  "Arizona",
		Zip:    "85281",
	})

	uid := "user-id-1"

	m := &proto.MutateRequest{
		Authorization: &proto.Authorization{WorkspaceId: "workspace", Source: &proto.VendorApp{
			VendorId: "chive",
			AppId:    "keystone",
		}},
		EntityId: "",
		Schema: &proto.Key{
			Key: "customer",
		},
		Mutation: &proto.Mutation{
			Mutator: &proto.Mutator{
				UserAgent: "Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_8_5) AppleWebKit/602.21 (KHTML, like Gecko) Chrome/49.0.2155.373 Safari/534",
				RemoteIp:  "123.45.67.89",
				UserId:    "1cfe0b4f-36b9-4576-b6e9-473f7e358e24",
				Client:    "Golang/SDK",
			},
			Comment: "Customer Creation",
			Properties: []*proto.EntityProperty{
				{
					Property: &proto.Key{
						Key: "first_name",
					},
					Value: &proto.Value{
						Text:       "J**n",
						SecureText: "John",
					},
				},
				{
					Property: &proto.Key{
						Key: "last_name",
					},
					Value: &proto.Value{
						Text: "Smith",
					},
				},
				{
					Property: &proto.Key{
						Key: "external_id",
					},
					Value: &proto.Value{
						Text: uid,
					},
				},
				{
					Property: &proto.Key{
						Key: "transaction_id",
					},
					Value: &proto.Value{
						Text: "last-transaction-id-1234",
					},
				},
				{
					Property: &proto.Key{
						Key: "total_paid",
					},
					Value: &proto.Value{
						Text: "GBP",
						Int:  1298,
					},
				},
				{
					Property: &proto.Key{
						Key: "subscriptions",
					},
					Value: &proto.Value{
						Int: 2,
					},
				},
				{
					Property: &proto.Key{
						Key: "first_paid",
					},
					Value: &proto.Value{
						Time: timestamppb.New(time.Now()),
					},
				},
				{
					Property: &proto.Key{
						Key: "fraud_rating",
					},
					Value: &proto.Value{
						Float: 0.12,
					},
				},
				{
					Property: &proto.Key{
						Key: "renewing",
					},
					Value: &proto.Value{
						Bool: true,
					},
				},
			},
			RemoveProperties: nil,
			Logs: []*proto.EntityLog{
				{
					Actor:     "Random Thing",
					Level:     proto.LogLevel_Info,
					Message:   "This is a log message",
					Reference: "ref1234",
					TraceId:   "trace-1234-53234-32427",
				},
			},
			Events: []*proto.EntityEvent{
				{
					Type: &proto.Key{
						Key: "creation",
					},
				},
			},
			Children: []*proto.EntityChild{
				{
					Type: &proto.Key{Key: "address"},
					Cid:  "address1",
					Data: addr1Bytes,
				},
			},
			RemoveChildren: []*proto.EntityChild{
				{
					Type: &proto.Key{Key: "address"},
					Cid:  "cuhf4cezor00Zm",
				},
				{
					Type: &proto.Key{Key: "address"},
					Cid:  "cuhf4gu2a1t4ya",
				},
			},
			Relationships: []*proto.EntityRelationship{
				{
					Relationship: &proto.Key{
						Key: "parent",
					},
					TargetId: "missing-id",
				},
			},
			//RemoveRelationships: nil,
			Links: []*proto.EntityLink{
				{
					Type: &proto.Key{
						Key: "website",
					},
					Location: "https://microsoft.com",
					Name:     "Microsoft Site",
				},
			},
			//RemoveLinks:         nil,
			Labels: []*proto.EntityLabel{
				{Name: "account-manager", Value: "frank_turner"},
			},
			RemoveLabels: []*proto.EntityLabel{
				{Name: "paid"},
			},
			Datum: []byte("This is a test, could be a webhook or something"),
		},
	}

	getTestActor(nil).Mutate(m, "Faker Customer x")
}

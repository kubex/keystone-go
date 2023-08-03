package keystone

import (
	"context"
	"encoding/json"
	"github.com/kubex/keystone-go/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"os"
	"testing"
	"time"
)

func TestConnection(t *testing.T) {

	kHost := os.Getenv("KEYSTONE_SERVICE_HOST")
	kPort := os.Getenv("KEYSTONE_SERVICE_PORT")
	if kHost == "" {
		kHost = "127.0.0.1"
	}
	if kPort == "" {
		kPort = "50031"
	}

	ksGrpcConn, err := grpc.Dial(kHost+":"+kPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	log.Println(kHost + ":" + kPort)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	ksClient := proto.NewKeystoneClient(ksGrpcConn)
	c := NewConnection(ksClient, "vendor", "appid", "accessToken")
	actor := c.Actor("test-workspace", "123.45.67.89", "user-1234", "User Agent Unknown")

	c.RegisterTypes(testSchemaType{}, Customer{})
	c.SyncSchema().Wait()

	log.Println("Marshalling")
	actor.Marshal(Customer{
		//ID:            "23of2DIcK7WUli7A",
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
		AmountPaid:    NewAmount("USD", 12353),
		LeadDate:      time.Now(),
		UserID:        "user-237",
	}, "Creating Customer via Marshal")
}

func xx(t *testing.T) {

	kHost := os.Getenv("KEYSTONE_SERVICE_HOST")
	kPort := os.Getenv("KEYSTONE_SERVICE_PORT")
	if kHost == "" {
		kHost = "127.0.0.1"
	}
	if kPort == "" {
		kPort = "50031"
	}

	ksGrpcConn, err := grpc.Dial(kHost+":"+kPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	log.Println(kHost + ":" + kPort)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	ksClient := proto.NewKeystoneClient(ksGrpcConn)

	c := NewConnection(ksClient, "vendor", "appid", "accessToken")

	c.RegisterTypes(testSchemaType{}, Customer{})
	c.SyncSchema().Wait()

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
					Key: &proto.Key{
						Key: "first_name",
					},
					Value: &proto.Value{
						Text:       "J**n",
						SecureText: "John",
					},
				},
				{
					Key: &proto.Key{
						Key: "last_name",
					},
					Value: &proto.Value{
						Text: "Smith",
					},
				},
				{
					Key: &proto.Key{
						Key: "external_id",
					},
					Value: &proto.Value{
						Text: uid,
					},
				},
				{
					Key: &proto.Key{
						Key: "transaction_id",
					},
					Value: &proto.Value{
						Text: "last-transaction-id-1234",
					},
				},
				{
					Key: &proto.Key{
						Key: "total_paid",
					},
					Value: &proto.Value{
						Text: "GBP",
						Int:  1298,
					},
				},
				{
					Key: &proto.Key{
						Key: "subscriptions",
					},
					Value: &proto.Value{
						Int: 2,
					},
				},
				{
					Key: &proto.Key{
						Key: "first_paid",
					},
					Value: &proto.Value{
						Time: timestamppb.New(time.Now()),
					},
				},
				{
					Key: &proto.Key{
						Key: "fraud_rating",
					},
					Value: &proto.Value{
						Float: 0.12,
					},
				},
				{
					Key: &proto.Key{
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
	c.ProtoClient().Mutate(context.Background(), m)
}

type Customer struct {
	ID               string       `keystone:"_entity_id" json:",omitempty"`
	Name             SecretString `keystone:",indexed,personal,omitempty" json:",omitempty"`
	Email            SecretString `keystone:",indexed,omitempty" json:",omitempty"`
	Company          string       `keystone:",omitempty" json:",omitempty"`
	Phone            string       `keystone:",indexed,omitempty" json:",omitempty"`
	HasPaid          bool         `keystone:",omitempty" json:",omitempty"`
	AvatarUrl        string       `keystone:",omitempty" json:",omitempty"`
	City             string       `keystone:",indexed,omitempty" json:",omitempty"`
	StreetName       string       `keystone:",omitempty" json:",omitempty"`
	StreetAddress    string       `keystone:",omitempty" json:",omitempty"`
	SecondaryAddress string       `keystone:",omitempty" json:",omitempty"`
	BuildingNumber   string       `keystone:",omitempty" json:",omitempty"`
	Postcode         string       `keystone:",lookup,omitempty" json:",omitempty"`
	Zipcode          string       `keystone:",omitempty" json:",omitempty"`
	Timezone         string       `keystone:",omitempty" json:",omitempty"`
	State            string       `keystone:",omitempty" json:",omitempty"`
	StateAbbr        string       `keystone:",omitempty" json:",omitempty"`
	Country          string       `keystone:",omitempty" json:",omitempty"`
	CountryCode      string       `keystone:",omitempty" json:",omitempty"`
	Latitude         float64      `keystone:",omitempty" json:",omitempty"`
	Longitude        float64      `keystone:",omitempty" json:",omitempty"`
	AmountPaid       Amount       `keystone:",omitempty" json:",omitempty"`
	LeadDate         time.Time    `keystone:",omitempty" json:",omitempty"`
	UserID           string       `keystone:",unique,omitempty" json:",omitempty"`
}

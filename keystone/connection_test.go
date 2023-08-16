package keystone

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ggwhite/go-masker"
	"github.com/kubex/keystone-go/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"syreclabs.com/go/faker"
)

var ksClient proto.KeystoneClient

func init() {
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
	ksClient = proto.NewKeystoneClient(ksGrpcConn)
}

func TestWrite(t *testing.T) {
	c := NewConnection(ksClient, "vendor", "appid", "accessToken")
	actor := c.Actor("test-workspace", "123.45.67.89", "user-1234", "User Agent Unknown")

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
	if err := actor.Marshal(cust, "Creating Customer via Marshal"); err != nil {
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
	c := NewConnection(ksClient, "vendor", "appid", "accessToken")
	actor := c.Actor("test-workspace", "123.45.67.89", "user-1234", "User Agent Unknown")

	c.RegisterTypes( /*testSchemaType{},*/ Customer{})
	c.SyncSchema().Wait()

	log.Println("Marshalling")
	for i := 0; i < 10; i++ {
		actor.Marshal(FakeCustomer(), "Faker Customer x")
	}
	return
}

func xx(t *testing.T) {
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
	c.ProtoClient().Mutate(context.Background(), m)
}

type Customer struct {
	EntityLogger
	EntityEvents
	EntityLabels
	EntityLinks
	EntityRelationships
	ID                  string       `keystone:"_entity_id" json:",omitempty"`
	Name                SecretString `keystone:",indexed,personal,omitempty" json:",omitempty"`
	Email               SecretString `keystone:",indexed,omitempty" json:",omitempty"`
	Company             string       `keystone:",omitempty" json:",omitempty"`
	Phone               string       `keystone:",omitempty" json:",omitempty"`
	HasPaid             bool         `keystone:",omitempty" json:",omitempty"`
	AvatarUrl           string       `keystone:",omitempty" json:",omitempty"`
	City                string       `keystone:",omitempty" json:",omitempty"`
	StreetName          string       `keystone:",omitempty" json:",omitempty"`
	StreetAddress       string       `keystone:",omitempty" json:",omitempty"`
	SecondaryAddress    string       `keystone:",omitempty" json:",omitempty"`
	BuildingNumber      string       `keystone:",omitempty" json:",omitempty"`
	Postcode            string       `keystone:",lookup,omitempty" json:",omitempty"`
	Zipcode             string       `keystone:",omitempty" json:",omitempty"`
	Timezone            string       `keystone:",omitempty" json:",omitempty"`
	State               string       `keystone:",omitempty" json:",omitempty"`
	StateAbbr           string       `keystone:",omitempty" json:",omitempty"`
	James               string       `keystone:",omitempty,indexed" json:",omitempty"`
	Country             string       `keystone:",omitempty" json:",omitempty"`
	CountryCode         string       `keystone:",omitempty" json:",omitempty"`
	Latitude            float64      `keystone:",omitempty" json:",omitempty"`
	Longitude           float64      `keystone:",omitempty" json:",omitempty"`
	AmountPaid          Amount       `keystone:",indexed,omitempty" json:",omitempty"`
	LeadDate            time.Time    `keystone:",omitempty" json:",omitempty"`
	UserID              string       `keystone:",unique,omitempty" json:",omitempty"`
	Address             Address
	References          []string
	LineItems           []*LineItem // TODO: Store as children?
	DiscountedLineItems []LineItem  // TODO: Store as children?
}

type LineItem struct {
	ID   string
	Name string
}

func (l LineItem) ChildID() string {
	return l.ID
}

type Address struct {
	Line1 string
	Line2 string
	City  string
}

func FakeCustomer() Customer {

	name := faker.Name().Name()
	email := faker.Internet().Email()

	c := Customer{}
	c.Name = NewSecretString(name, masker.Name(name))
	c.Email = NewSecretString(email, masker.Email(email))
	c.Company = faker.Company().Name()
	c.Phone = faker.PhoneNumber().String()
	c.AvatarUrl = faker.Avatar().Url("png", 100, 100)
	c.City = faker.Address().City()
	c.StreetName = faker.Address().StreetName()
	c.StreetAddress = faker.Address().StreetAddress()
	c.SecondaryAddress = faker.Address().SecondaryAddress()
	c.BuildingNumber = faker.Address().BuildingNumber()
	c.Postcode = faker.Address().Postcode()
	c.Zipcode = faker.Address().ZipCode()
	c.Timezone = faker.Address().TimeZone()
	c.State = faker.Address().State()
	c.Country = faker.Address().Country()
	c.CountryCode = faker.Address().CountryCode()
	c.Latitude = float64(faker.Address().Latitude())
	c.Longitude = float64(faker.Address().Longitude())
	c.AmountPaid = NewAmount("GBP", int64(faker.Commerce().Price()*100))
	c.LeadDate = faker.Time().Birthday(0, 5)
	c.UserID = faker.RandomString(10)

	c.LogDebug("Created customer", "REF123", "trace-id", "mr-man", nil)

	return c
}

func TestFind(t *testing.T) {
	c := NewConnection(ksClient, "vendor", "appid", "accessToken")
	actor := c.Actor("test-workspace", "123.45.67.89", "user-1234", "User Agent Unknown")

	cst := Customer{}
	c.RegisterTypes(cst)
	c.SyncSchema().Wait()
	schema, _ := actor.connection.registerType(cst)

	res, err := c.ProtoClient().Find(context.Background(), &proto.FindRequest{
		Authorization: actor.authorization(),
		Schema:        &proto.Key{Key: schema.Type, Source: schema.Source},
		Properties: []*proto.PropertyRequest{
			{
				Properties: []string{"name"},
			},
		},
		Filters: []*proto.PropertyFilter{
			{
				Property: &proto.Key{Key: "name"},
				Operator: proto.Operator_Contains,
				Values: []*proto.Value{
					{Text: "j"},
				},
			},
			//{
			//	Property: &proto.Key{Key: "james"},
			//	Operator: proto.Operator_Between,
			//	Values: []*proto.Value{
			//		{Text: "A"},
			//		{Text: "F"},
			//	},
			//},
			//{
			//	Property: &proto.Key{Key: "user_id"},
			//	Operator: proto.Operator_In,
			//	Values: []*proto.Value{
			//		{Text: "user-2331"},
			//		{Text: "user-233"},
			//		{Text: "user-212"},
			//		{Text: "user-233120230810175925"},
			//	},
			//},
			//{
			//	Property: &proto.Key{Key: "postcode"},
			//	Operator: proto.Operator_In,
			//	Values: []*proto.Value{
			//		//{Text: "PO1 3AG"},
			//		{Text: "PO2 1AG"},
			//	},
			//},
		},
		//Labels: []*proto.EntityLabel{
		//	{
		//		Name: "active",
		//	},
		//},

	})
	log.Println(res, err)

}

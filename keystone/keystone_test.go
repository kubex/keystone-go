package keystone

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ggwhite/go-masker"
	"github.com/kubex/keystone-go/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"syreclabs.com/go/faker"
)

var testAddressResponse = &proto.EntityResponse{
	Entity: &proto.Entity{EntityId: "xxx-xxxx"},
	Properties: []*proto.EntityProperty{
		{Property: &proto.Key{Key: "line1"}, Value: &proto.Value{Text: "123 Fake St."}},
		{Property: &proto.Key{Key: "line2"}, Value: &proto.Value{Text: "Apt. 2"}},
		{Property: &proto.Key{Key: "city"}, Value: &proto.Value{Text: "Springfield"}},
	},
}

var testAddressResponseWithPostCode = &proto.EntityResponse{
	Entity: &proto.Entity{EntityId: "xxx-xxxx"},
	Properties: []*proto.EntityProperty{
		{Property: &proto.Key{Key: "line1"}, Value: &proto.Value{Text: "123 Fake St."}},
		{Property: &proto.Key{Key: "line2"}, Value: &proto.Value{Text: "Apt. 2"}},
		{Property: &proto.Key{Key: "city"}, Value: &proto.Value{Text: "Springfield"}},
		{Property: &proto.Key{Key: "post_code"}, Value: &proto.Value{Text: "Springfield"}},
	},
}

var ksClient proto.KeystoneClient

func init() {
	host := os.Getenv("KEYSTONE_SERVICE_HOST")
	port := os.Getenv("KEYSTONE_SERVICE_PORT")
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "50051"
	}

	ksGrpcConn, err := grpc.Dial(host+":"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	log.Println(host + ":" + port)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	ksClient = proto.NewKeystoneClient(ksGrpcConn)
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
	AvatarURL           string       `keystone:",omitempty" json:",omitempty"`
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
	LineItems           []*LineItem
	DiscountedLineItems []LineItem
}

type LineItem struct {
	ID   string
	Name string
}

func (l LineItem) ChildID() string {
	return l.ID
}

type Address struct {
	BaseEntity
	Line1 string
	Line2 string
	City  string
}

type MockConnector struct {
	mutateError error
}

func (m *MockConnector) Define(ctx context.Context, in *proto.SchemaRequest, opts ...grpc.CallOption) (*proto.Schema, error) {
	log.Println(ctx, in, opts)
	return &proto.Schema{}, nil
}

func (m *MockConnector) Mutate(ctx context.Context, in *proto.MutateRequest, opts ...grpc.CallOption) (*proto.MutateResponse, error) {
	log.Println(ctx, in, opts)
	r := &proto.MutateResponse{
		Success: true,
	}
	if m.mutateError != nil {
		r.Success = false
		r.ErrorMessage = m.mutateError.Error()
		r.ErrorCode = 500
	}

	return r, nil
}
func (m *MockConnector) Retrieve(ctx context.Context, in *proto.EntityRequest, opts ...grpc.CallOption) (*proto.EntityResponse, error) {
	log.Println(ctx, in, opts)
	return &proto.EntityResponse{}, nil
}
func (m *MockConnector) Logs(ctx context.Context, in *proto.LogRequest, opts ...grpc.CallOption) (*proto.LogsResponse, error) {
	log.Println(ctx, in, opts)
	return &proto.LogsResponse{}, nil
}
func (m *MockConnector) Events(ctx context.Context, in *proto.EventRequest, opts ...grpc.CallOption) (*proto.EventsResponse, error) {
	log.Println(ctx, in, opts)
	return &proto.EventsResponse{}, nil
}
func (m *MockConnector) Find(ctx context.Context, in *proto.FindRequest, opts ...grpc.CallOption) (*proto.FindResponse, error) {
	log.Println(ctx, in, opts)
	return &proto.FindResponse{}, nil
}

/*func (m *MockConnector) ADSList(ctx context.Context, in *proto.ADSListRequest, opts ...grpc.CallOption) (*proto.ADSListResponse, error) {
	log.Println(ctx, in, opts)
	return &proto.ADSListResponse{}, nil
}
func (m *MockConnector) ApplyADS(ctx context.Context, in *proto.ADS, opts ...grpc.CallOption) (*proto.GenericResponse, error) {
	log.Println(ctx, in, opts)
	return &proto.GenericResponse{}, nil
}*/

func FakeCustomer() *Customer {

	name := faker.Name().Name()
	email := faker.Internet().Email()

	c := &Customer{}
	c.Name = NewSecretString(name, masker.Name(name))
	c.Email = NewSecretString(email, masker.Email(email))
	c.Company = faker.Company().Name()
	c.Phone = faker.PhoneNumber().String()
	c.AvatarURL = faker.Avatar().Url("png", 100, 100)
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

	c.AddKeystoneEvent("custom-event", map[string]string{"foo": "bar"})
	c.AddKeystoneLabel("active", "true")
	c.AddKeystoneLink("gcs", "Raw Data", "https://google.com")
	c.AddKeystoneRelationship("src", "targ", map[string]string{"a": "b"}, time.Now())

	c.LogDebug("Created customer", "REF123", "trace-id", "mr-man", nil)

	return c
}

func getTestActor(client proto.KeystoneClient) *Actor {
	if client == nil {
		client = ksClient
	}
	c := NewConnection(client, "vendor", "appid", "accessToken")
	actor := c.Actor("test-workspace", "123.45.67.89", "user-1234", "User Agent Unknown")
	return &actor
}

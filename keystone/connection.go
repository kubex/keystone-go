package keystone

import (
	"context"
	"github.com/kubex/keystone-go/proto"
	"log"
	"reflect"
	"sync"
)

type Connection struct {
	client        proto.KeystoneClient
	appID         proto.VendorApp
	token         string
	typeRegister  map[reflect.Type]*proto.Schema
	registerQueue map[reflect.Type]bool // true if the type is processing registration
}

func NewConnection(client proto.KeystoneClient, vendorID, appID, accessToken string) *Connection {
	return &Connection{
		client:        client,
		appID:         proto.VendorApp{VendorId: vendorID, AppId: appID},
		token:         accessToken,
		typeRegister:  make(map[reflect.Type]*proto.Schema),
		registerQueue: make(map[reflect.Type]bool),
	}
}

func (c *Connection) ProtoClient() proto.KeystoneClient {
	return c.client
}

func (c *Connection) authorization() *proto.Authorization {
	return &proto.Authorization{
		Source: &c.appID,
		Token:  c.token,
	}
}

func (c *Connection) Actor(workspaceID, remoteIP, userID, userAgent string) Actor {
	return Actor{
		connection:  c,
		workspaceID: workspaceID,
		mutator: &proto.Mutator{
			UserAgent: userAgent,
			RemoteIp:  remoteIP,
			UserId:    userID,
			Client:    "Keystone Go-Client",
		},
	}
}

// RegisterTypes registers the given types with the connection, returning the number of new types registered
func (c *Connection) RegisterTypes(types ...interface{}) int {
	registered := 0
	for _, t := range types {
		if c.registerType(t) {
			registered++
		}
	}
	return registered
}

// registerType returns true if the type is already registered
func (c *Connection) registerType(t interface{}) bool {
	typ := reflect.TypeOf(t)
	if _, ok := c.typeRegister[typ]; !ok {
		c.typeRegister[typ] = typeToSchema(t)
		c.registerQueue[typ] = false
		return false
	}
	return true
}

func (c *Connection) SyncSchema() *sync.WaitGroup {
	wg := &sync.WaitGroup{}
	wg.Add(len(c.registerQueue))
	go func() {
		for typ, processing := range c.registerQueue {
			if !processing {
				if toRegister, ok := c.typeRegister[typ]; ok {
					log.Println("Registering type", typ)
					resp, err := c.ProtoClient().Define(context.Background(), toRegister)
					log.Println(resp, err)
				}
				wg.Done()
			}
		}
	}()
	return wg
}

type Actor struct {
	connection  *Connection
	workspaceID string
	mutator     *proto.Mutator
}

func (a *Actor) SetClient(client string) {
	a.mutator.Client = client
}

func (a *Actor) Marshal(entity interface{}) {
	log.Println("Processing Marshal request")
	if !a.connection.registerType(entity) {
		// wait for the type to be registered with the keystone server
		a.connection.SyncSchema().Wait()
	}
	log.Println("Marshalling entity", entity)
}

func (c *Connection) Retrieve(ctx context.Context, workspaceID, entityId string, retrieveProperties []string) (*proto.EntityResponse, error) {
	return nil, nil
}

func (c *Connection) Lookup(ctx context.Context, workspaceID, idLookup string) ([]*proto.EntityResponse, error) {
	return nil, nil
}

func (c *Connection) Find(ctx context.Context, workspaceID, entityType string, retrieveProperties []string, options ...Option) ([]*proto.EntityResponse, error) {
	return nil, nil
}

package keystone

import (
	"context"
	"log"
	"reflect"
	"sync"

	"github.com/kubex/keystone-go/proto"
)

// Connection is a connection to a keystone server
type Connection struct {
	client        proto.KeystoneClient
	appID         proto.VendorApp
	token         string
	typeRegister  map[reflect.Type]*proto.Schema
	registerQueue map[reflect.Type]bool // true if the type is processing registration
}

// NewConnection creates a new connection to a keystone server
func NewConnection(client proto.KeystoneClient, vendorID, appID, accessToken string) *Connection {
	return &Connection{
		client:        client,
		appID:         proto.VendorApp{VendorId: vendorID, AppId: appID},
		token:         accessToken,
		typeRegister:  make(map[reflect.Type]*proto.Schema),
		registerQueue: make(map[reflect.Type]bool),
	}
}

// ProtoClient returns the underlying proto client
func (c *Connection) ProtoClient() proto.KeystoneClient {
	return c.client
}

func (c *Connection) authorization() *proto.Authorization {
	return &proto.Authorization{
		Source: &c.appID,
		Token:  c.token,
	}
}

// Actor returns an actor for the given workspace, remote IP, user ID, and user agent
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
		if _, reg := c.registerType(t); !reg {
			registered++
		}
	}
	return registered
}

// registerType returns true if the type is already registered
func (c *Connection) registerType(t interface{}) (*proto.Schema, bool) {
	typ := reflect.TypeOf(t)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	schema, ok := c.typeRegister[typ]
	if !ok {
		newSchema := typeToSchema(t)
		c.typeRegister[typ] = newSchema
		c.registerQueue[typ] = false
		return newSchema, false
	}
	return schema, true
}

// SyncSchema syncs the schema with the server
func (c *Connection) SyncSchema() *sync.WaitGroup {
	wg := &sync.WaitGroup{}
	wg.Add(len(c.registerQueue))
	go func() {
		for typ, processing := range c.registerQueue {
			if !processing {
				if toRegister, ok := c.typeRegister[typ]; ok {
					log.Println("Registering type", typ)
					resp, err := c.ProtoClient().Define(context.Background(), &proto.SchemaRequest{
						Authorization: c.authorization(),
						Schema:        toRegister,
					})
					if err == nil {
						c.typeRegister[typ].Id = resp.GetId()
						c.typeRegister[typ].Name = resp.GetName()
						c.typeRegister[typ].Source = resp.GetSource()
						c.typeRegister[typ].Type = resp.GetType()
						c.typeRegister[typ].Properties = resp.GetProperties()
						c.typeRegister[typ].Options = resp.GetOptions()
						c.typeRegister[typ].Singular = resp.GetSingular()
						c.typeRegister[typ].Plural = resp.GetPlural()
					}
				}
				wg.Done()
			}
		}
	}()
	return wg
}

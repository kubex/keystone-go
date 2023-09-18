package keystone

import (
	"context"
	"github.com/packaged/logger/v3/logger"
	"go.uber.org/zap"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/kubex/keystone-go/proto"
	"google.golang.org/grpc"
)

// Connection is a connection to a keystone server
type Connection struct {
	client        proto.KeystoneClient
	logger        *logger.Logger
	timeLogConfig *logger.TimedLogConfig
	appID         proto.VendorApp
	token         string
	typeRegister  map[reflect.Type]schemaDef
	registerQueue map[reflect.Type]bool // true if the type is processing registration
}

// NewConnection creates a new connection to a keystone server
func NewConnection(client proto.KeystoneClient, vendorID, appID, accessToken string) *Connection {
	return &Connection{
		timeLogConfig: &logger.TimedLogConfig{
			ErrorDuration: time.Minute,
			WarnDuration:  30 * time.Second,
			InfoDuration:  2 * time.Second,
			DebugDuration: 100 * time.Millisecond,
		},
		logger:        logger.I(),
		client:        client,
		appID:         proto.VendorApp{VendorId: vendorID, AppId: appID},
		token:         accessToken,
		typeRegister:  make(map[reflect.Type]schemaDef),
		registerQueue: make(map[reflect.Type]bool),
	}
}

func (c *Connection) Define(ctx context.Context, in *proto.SchemaRequest, opts ...grpc.CallOption) (*proto.Schema, error) {
	tl := c.timeLogConfig.NewLog("Define", zap.String("schema", in.GetSchema().GetType()))
	resp, err := c.client.Define(ctx, in, opts...)
	c.logger.TimedLog(tl)
	return resp, err
}

func (c *Connection) Mutate(ctx context.Context, in *proto.MutateRequest, opts ...grpc.CallOption) (*proto.MutateResponse, error) {
	tl := c.timeLogConfig.NewLog("Mutate", zap.String("EntityId", in.GetEntityId()))
	resp, err := c.client.Mutate(ctx, in, opts...)
	c.logger.TimedLog(tl)
	return resp, err
}

func (c *Connection) Retrieve(ctx context.Context, in *proto.EntityRequest, opts ...grpc.CallOption) (*proto.EntityResponse, error) {
	tl := c.timeLogConfig.NewLog("Retrieve", zap.String("EntityId", in.GetEntityId()))
	resp, err := c.client.Retrieve(ctx, in, opts...)
	c.logger.TimedLog(tl)
	return resp, err
}

func (c *Connection) Logs(ctx context.Context, in *proto.LogRequest, opts ...grpc.CallOption) (*proto.LogsResponse, error) {
	tl := c.timeLogConfig.NewLog("Logs", zap.String("EntityId", in.GetEntityId()))
	resp, err := c.client.Logs(ctx, in, opts...)
	c.logger.TimedLog(tl)
	return resp, err
}

func (c *Connection) Events(ctx context.Context, in *proto.EventRequest, opts ...grpc.CallOption) (*proto.EventsResponse, error) {
	tl := c.timeLogConfig.NewLog("Events", zap.String("EntityId", in.GetEntityId()))
	resp, err := c.client.Events(ctx, in, opts...)
	c.logger.TimedLog(tl)
	return resp, err
}

func (c *Connection) Find(ctx context.Context, in *proto.FindRequest, opts ...grpc.CallOption) (*proto.FindResponse, error) {
	tl := c.timeLogConfig.NewLog("Find", zap.String("schema", in.GetSchema().GetKey()))
	resp, err := c.client.Find(ctx, in, opts...)
	c.logger.TimedLog(tl)
	return resp, err
}

func (c *Connection) ADSList(ctx context.Context, in *proto.ADSListRequest, opts ...grpc.CallOption) (*proto.ADSListResponse, error) {
	tl := c.timeLogConfig.NewLog("ADSList", zap.String("schema", in.GetSchema().GetKey()), zap.String("ADS", in.GetAdsName()))
	resp, err := c.client.ADSList(ctx, in, opts...)
	c.logger.TimedLog(tl)
	return resp, err
}

func (c *Connection) ApplyADS(ctx context.Context, in *proto.ADS, opts ...grpc.CallOption) (*proto.GenericResponse, error) {
	tl := c.timeLogConfig.NewLog("ADSList", zap.String("schema", in.GetSchema().GetKey()), zap.String("ADS", in.GetName()))
	resp, err := c.client.ApplyADS(ctx, in, opts...)
	c.logger.TimedLog(tl)
	return resp, err
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

	sDef, ok := c.typeRegister[typ]
	if !ok {
		newDef := typeToSchema(t)
		c.typeRegister[typ] = newDef
		c.registerQueue[typ] = false
		return newDef.schema, false
	}
	return sDef.schema, true
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
					resp, err := c.Define(context.Background(), &proto.SchemaRequest{
						Authorization: c.authorization(),
						Schema:        toRegister.schema,
					})
					if err == nil {
						c.typeRegister[typ].schema.Id = resp.GetId()
						c.typeRegister[typ].schema.Name = resp.GetName()
						c.typeRegister[typ].schema.Source = resp.GetSource()
						c.typeRegister[typ].schema.Type = resp.GetType()
						c.typeRegister[typ].schema.Properties = resp.GetProperties()
						c.typeRegister[typ].schema.Options = resp.GetOptions()
						c.typeRegister[typ].schema.Singular = resp.GetSingular()
						c.typeRegister[typ].schema.Plural = resp.GetPlural()
					}

					if toRegister.definition.ActiveDataSets != nil {
						for _, ads := range toRegister.definition.ActiveDataSets {
							adsResp, adsErr := c.ApplyADS(context.Background(), ads)
							if adsErr != nil || !adsResp.GetSuccess() {
								log.Println("Error applying ADS", adsErr, adsResp)
							}
						}
					}

				}
				wg.Done()
			}
		}
	}()
	return wg
}

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: keystone.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Keystone_Define_FullMethodName       = "/kubex.keystone.Keystone/Define"
	Keystone_Mutate_FullMethodName       = "/kubex.keystone.Keystone/Mutate"
	Keystone_Retrieve_FullMethodName     = "/kubex.keystone.Keystone/Retrieve"
	Keystone_RetrieveView_FullMethodName = "/kubex.keystone.Keystone/RetrieveView"
	Keystone_Logs_FullMethodName         = "/kubex.keystone.Keystone/Logs"
	Keystone_Events_FullMethodName       = "/kubex.keystone.Keystone/Events"
	Keystone_Find_FullMethodName         = "/kubex.keystone.Keystone/Find"
)

// KeystoneClient is the client API for Keystone service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KeystoneClient interface {
	Define(ctx context.Context, in *SchemaRequest, opts ...grpc.CallOption) (*Schema, error)
	Mutate(ctx context.Context, in *MutateRequest, opts ...grpc.CallOption) (*MutateResponse, error)
	Retrieve(ctx context.Context, in *EntityRequest, opts ...grpc.CallOption) (*EntityResponse, error)
	RetrieveView(ctx context.Context, in *EntityViewRequest, opts ...grpc.CallOption) (*EntitiesResponse, error)
	Logs(ctx context.Context, in *LogRequest, opts ...grpc.CallOption) (*LogsResponse, error)
	Events(ctx context.Context, in *EventRequest, opts ...grpc.CallOption) (*EventsResponse, error)
	Find(ctx context.Context, in *FindRequest, opts ...grpc.CallOption) (*FindResponse, error)
}

type keystoneClient struct {
	cc grpc.ClientConnInterface
}

func NewKeystoneClient(cc grpc.ClientConnInterface) KeystoneClient {
	return &keystoneClient{cc}
}

func (c *keystoneClient) Define(ctx context.Context, in *SchemaRequest, opts ...grpc.CallOption) (*Schema, error) {
	out := new(Schema)
	err := c.cc.Invoke(ctx, Keystone_Define_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keystoneClient) Mutate(ctx context.Context, in *MutateRequest, opts ...grpc.CallOption) (*MutateResponse, error) {
	out := new(MutateResponse)
	err := c.cc.Invoke(ctx, Keystone_Mutate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keystoneClient) Retrieve(ctx context.Context, in *EntityRequest, opts ...grpc.CallOption) (*EntityResponse, error) {
	out := new(EntityResponse)
	err := c.cc.Invoke(ctx, Keystone_Retrieve_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keystoneClient) RetrieveView(ctx context.Context, in *EntityViewRequest, opts ...grpc.CallOption) (*EntitiesResponse, error) {
	out := new(EntitiesResponse)
	err := c.cc.Invoke(ctx, Keystone_RetrieveView_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keystoneClient) Logs(ctx context.Context, in *LogRequest, opts ...grpc.CallOption) (*LogsResponse, error) {
	out := new(LogsResponse)
	err := c.cc.Invoke(ctx, Keystone_Logs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keystoneClient) Events(ctx context.Context, in *EventRequest, opts ...grpc.CallOption) (*EventsResponse, error) {
	out := new(EventsResponse)
	err := c.cc.Invoke(ctx, Keystone_Events_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keystoneClient) Find(ctx context.Context, in *FindRequest, opts ...grpc.CallOption) (*FindResponse, error) {
	out := new(FindResponse)
	err := c.cc.Invoke(ctx, Keystone_Find_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KeystoneServer is the server API for Keystone service.
// All implementations must embed UnimplementedKeystoneServer
// for forward compatibility
type KeystoneServer interface {
	Define(context.Context, *SchemaRequest) (*Schema, error)
	Mutate(context.Context, *MutateRequest) (*MutateResponse, error)
	Retrieve(context.Context, *EntityRequest) (*EntityResponse, error)
	RetrieveView(context.Context, *EntityViewRequest) (*EntitiesResponse, error)
	Logs(context.Context, *LogRequest) (*LogsResponse, error)
	Events(context.Context, *EventRequest) (*EventsResponse, error)
	Find(context.Context, *FindRequest) (*FindResponse, error)
	mustEmbedUnimplementedKeystoneServer()
}

// UnimplementedKeystoneServer must be embedded to have forward compatible implementations.
type UnimplementedKeystoneServer struct {
}

func (UnimplementedKeystoneServer) Define(context.Context, *SchemaRequest) (*Schema, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Define not implemented")
}
func (UnimplementedKeystoneServer) Mutate(context.Context, *MutateRequest) (*MutateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Mutate not implemented")
}
func (UnimplementedKeystoneServer) Retrieve(context.Context, *EntityRequest) (*EntityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Retrieve not implemented")
}
func (UnimplementedKeystoneServer) RetrieveView(context.Context, *EntityViewRequest) (*EntitiesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RetrieveView not implemented")
}
func (UnimplementedKeystoneServer) Logs(context.Context, *LogRequest) (*LogsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logs not implemented")
}
func (UnimplementedKeystoneServer) Events(context.Context, *EventRequest) (*EventsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Events not implemented")
}
func (UnimplementedKeystoneServer) Find(context.Context, *FindRequest) (*FindResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Find not implemented")
}
func (UnimplementedKeystoneServer) mustEmbedUnimplementedKeystoneServer() {}

// UnsafeKeystoneServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KeystoneServer will
// result in compilation errors.
type UnsafeKeystoneServer interface {
	mustEmbedUnimplementedKeystoneServer()
}

func RegisterKeystoneServer(s grpc.ServiceRegistrar, srv KeystoneServer) {
	s.RegisterService(&Keystone_ServiceDesc, srv)
}

func _Keystone_Define_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SchemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeystoneServer).Define(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Keystone_Define_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeystoneServer).Define(ctx, req.(*SchemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Keystone_Mutate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MutateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeystoneServer).Mutate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Keystone_Mutate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeystoneServer).Mutate(ctx, req.(*MutateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Keystone_Retrieve_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EntityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeystoneServer).Retrieve(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Keystone_Retrieve_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeystoneServer).Retrieve(ctx, req.(*EntityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Keystone_RetrieveView_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EntityViewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeystoneServer).RetrieveView(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Keystone_RetrieveView_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeystoneServer).RetrieveView(ctx, req.(*EntityViewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Keystone_Logs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeystoneServer).Logs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Keystone_Logs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeystoneServer).Logs(ctx, req.(*LogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Keystone_Events_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeystoneServer).Events(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Keystone_Events_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeystoneServer).Events(ctx, req.(*EventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Keystone_Find_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeystoneServer).Find(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Keystone_Find_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeystoneServer).Find(ctx, req.(*FindRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Keystone_ServiceDesc is the grpc.ServiceDesc for Keystone service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Keystone_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "kubex.keystone.Keystone",
	HandlerType: (*KeystoneServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Define",
			Handler:    _Keystone_Define_Handler,
		},
		{
			MethodName: "Mutate",
			Handler:    _Keystone_Mutate_Handler,
		},
		{
			MethodName: "Retrieve",
			Handler:    _Keystone_Retrieve_Handler,
		},
		{
			MethodName: "RetrieveView",
			Handler:    _Keystone_RetrieveView_Handler,
		},
		{
			MethodName: "Logs",
			Handler:    _Keystone_Logs_Handler,
		},
		{
			MethodName: "Events",
			Handler:    _Keystone_Events_Handler,
		},
		{
			MethodName: "Find",
			Handler:    _Keystone_Find_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "keystone.proto",
}

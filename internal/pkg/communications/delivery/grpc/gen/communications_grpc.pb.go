// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.3
// source: communications.proto

package gen

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Communications_AddReaction_FullMethodName     = "/communications.Communications/AddReaction"
	Communications_GetMatchList_FullMethodName    = "/communications.Communications/GetMatchList"
	Communications_GetReactionList_FullMethodName = "/communications.Communications/GetReactionList"
)

// CommunicationsClient is the client API for Communications service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommunicationsClient interface {
	AddReaction(ctx context.Context, in *AddReactionRequest, opts ...grpc.CallOption) (*AddReactionResponse, error)
	GetMatchList(ctx context.Context, in *GetMatchListRequest, opts ...grpc.CallOption) (*GetMatchListResponse, error)
	GetReactionList(ctx context.Context, in *GetReactionListRequest, opts ...grpc.CallOption) (*GetReactionListResponse, error)
}

type communicationsClient struct {
	cc grpc.ClientConnInterface
}

func NewCommunicationsClient(cc grpc.ClientConnInterface) CommunicationsClient {
	return &communicationsClient{cc}
}

func (c *communicationsClient) AddReaction(ctx context.Context, in *AddReactionRequest, opts ...grpc.CallOption) (*AddReactionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddReactionResponse)
	err := c.cc.Invoke(ctx, Communications_AddReaction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *communicationsClient) GetMatchList(ctx context.Context, in *GetMatchListRequest, opts ...grpc.CallOption) (*GetMatchListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetMatchListResponse)
	err := c.cc.Invoke(ctx, Communications_GetMatchList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *communicationsClient) GetReactionList(ctx context.Context, in *GetReactionListRequest, opts ...grpc.CallOption) (*GetReactionListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetReactionListResponse)
	err := c.cc.Invoke(ctx, Communications_GetReactionList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommunicationsServer is the server API for Communications service.
// All implementations must embed UnimplementedCommunicationsServer
// for forward compatibility.
type CommunicationsServer interface {
	AddReaction(context.Context, *AddReactionRequest) (*AddReactionResponse, error)
	GetMatchList(context.Context, *GetMatchListRequest) (*GetMatchListResponse, error)
	GetReactionList(context.Context, *GetReactionListRequest) (*GetReactionListResponse, error)
	mustEmbedUnimplementedCommunicationsServer()
}

// UnimplementedCommunicationsServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCommunicationsServer struct{}

func (UnimplementedCommunicationsServer) AddReaction(context.Context, *AddReactionRequest) (*AddReactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddReaction not implemented")
}
func (UnimplementedCommunicationsServer) GetMatchList(context.Context, *GetMatchListRequest) (*GetMatchListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMatchList not implemented")
}
func (UnimplementedCommunicationsServer) GetReactionList(context.Context, *GetReactionListRequest) (*GetReactionListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReactionList not implemented")
}
func (UnimplementedCommunicationsServer) mustEmbedUnimplementedCommunicationsServer() {}
func (UnimplementedCommunicationsServer) testEmbeddedByValue()                        {}

// UnsafeCommunicationsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommunicationsServer will
// result in compilation errors.
type UnsafeCommunicationsServer interface {
	mustEmbedUnimplementedCommunicationsServer()
}

func RegisterCommunicationsServer(s grpc.ServiceRegistrar, srv CommunicationsServer) {
	// If the following call pancis, it indicates UnimplementedCommunicationsServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Communications_ServiceDesc, srv)
}

func _Communications_AddReaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddReactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunicationsServer).AddReaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Communications_AddReaction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunicationsServer).AddReaction(ctx, req.(*AddReactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Communications_GetMatchList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMatchListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunicationsServer).GetMatchList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Communications_GetMatchList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunicationsServer).GetMatchList(ctx, req.(*GetMatchListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Communications_GetReactionList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReactionListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunicationsServer).GetReactionList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Communications_GetReactionList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunicationsServer).GetReactionList(ctx, req.(*GetReactionListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Communications_ServiceDesc is the grpc.ServiceDesc for Communications service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Communications_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "communications.Communications",
	HandlerType: (*CommunicationsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddReaction",
			Handler:    _Communications_AddReaction_Handler,
		},
		{
			MethodName: "GetMatchList",
			Handler:    _Communications_GetMatchList_Handler,
		},
		{
			MethodName: "GetReactionList",
			Handler:    _Communications_GetReactionList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "communications.proto",
}
// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.29.1
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
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Communications_AddReaction_FullMethodName            = "/communications.Communications/AddReaction"
	Communications_GetMatchList_FullMethodName           = "/communications.Communications/GetMatchList"
	Communications_GetReactionList_FullMethodName        = "/communications.Communications/GetReactionList"
	Communications_GetMatchTime_FullMethodName           = "/communications.Communications/GetMatchTime"
	Communications_GetMatchesBySearch_FullMethodName     = "/communications.Communications/GetMatchesBySearch"
	Communications_UpdateOrCreateReaction_FullMethodName = "/communications.Communications/UpdateOrCreateReaction"
	Communications_CheckMatchExists_FullMethodName       = "/communications.Communications/CheckMatchExists"
)

// CommunicationsClient is the client API for Communications service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommunicationsClient interface {
	AddReaction(ctx context.Context, in *AddReactionRequest, opts ...grpc.CallOption) (*AddReactionResponse, error)
	GetMatchList(ctx context.Context, in *GetMatchListRequest, opts ...grpc.CallOption) (*GetMatchListResponse, error)
	GetReactionList(ctx context.Context, in *GetReactionListRequest, opts ...grpc.CallOption) (*GetReactionListResponse, error)
	GetMatchTime(ctx context.Context, in *GetMatchTimeRequest, opts ...grpc.CallOption) (*GetMatchTimeResponse, error)
	GetMatchesBySearch(ctx context.Context, in *GetMatchesBySearchRequest, opts ...grpc.CallOption) (*GetMatchesBySearchResponse, error)
	UpdateOrCreateReaction(ctx context.Context, in *UpdateOrCreateReactionRequest, opts ...grpc.CallOption) (*UpdateOrCreateReactionResponse, error)
	CheckMatchExists(ctx context.Context, in *CheckMatchExistsRequest, opts ...grpc.CallOption) (*CheckMatchExistsResponse, error)
}

type communicationsClient struct {
	cc grpc.ClientConnInterface
}

func NewCommunicationsClient(cc grpc.ClientConnInterface) CommunicationsClient {
	return &communicationsClient{cc}
}

func (c *communicationsClient) AddReaction(ctx context.Context, in *AddReactionRequest, opts ...grpc.CallOption) (*AddReactionResponse, error) {
	out := new(AddReactionResponse)
	err := c.cc.Invoke(ctx, Communications_AddReaction_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *communicationsClient) GetMatchList(ctx context.Context, in *GetMatchListRequest, opts ...grpc.CallOption) (*GetMatchListResponse, error) {
	out := new(GetMatchListResponse)
	err := c.cc.Invoke(ctx, Communications_GetMatchList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *communicationsClient) GetReactionList(ctx context.Context, in *GetReactionListRequest, opts ...grpc.CallOption) (*GetReactionListResponse, error) {
	out := new(GetReactionListResponse)
	err := c.cc.Invoke(ctx, Communications_GetReactionList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *communicationsClient) GetMatchTime(ctx context.Context, in *GetMatchTimeRequest, opts ...grpc.CallOption) (*GetMatchTimeResponse, error) {
	out := new(GetMatchTimeResponse)
	err := c.cc.Invoke(ctx, Communications_GetMatchTime_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *communicationsClient) GetMatchesBySearch(ctx context.Context, in *GetMatchesBySearchRequest, opts ...grpc.CallOption) (*GetMatchesBySearchResponse, error) {
	out := new(GetMatchesBySearchResponse)
	err := c.cc.Invoke(ctx, Communications_GetMatchesBySearch_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *communicationsClient) UpdateOrCreateReaction(ctx context.Context, in *UpdateOrCreateReactionRequest, opts ...grpc.CallOption) (*UpdateOrCreateReactionResponse, error) {
	out := new(UpdateOrCreateReactionResponse)
	err := c.cc.Invoke(ctx, Communications_UpdateOrCreateReaction_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *communicationsClient) CheckMatchExists(ctx context.Context, in *CheckMatchExistsRequest, opts ...grpc.CallOption) (*CheckMatchExistsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CheckMatchExistsResponse)
	err := c.cc.Invoke(ctx, Communications_CheckMatchExists_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommunicationsServer is the server API for Communications service.
// All implementations must embed UnimplementedCommunicationsServer
// for forward compatibility
type CommunicationsServer interface {
	AddReaction(context.Context, *AddReactionRequest) (*AddReactionResponse, error)
	GetMatchList(context.Context, *GetMatchListRequest) (*GetMatchListResponse, error)
	GetReactionList(context.Context, *GetReactionListRequest) (*GetReactionListResponse, error)
	GetMatchTime(context.Context, *GetMatchTimeRequest) (*GetMatchTimeResponse, error)
	GetMatchesBySearch(context.Context, *GetMatchesBySearchRequest) (*GetMatchesBySearchResponse, error)
	UpdateOrCreateReaction(context.Context, *UpdateOrCreateReactionRequest) (*UpdateOrCreateReactionResponse, error)
	CheckMatchExists(context.Context, *CheckMatchExistsRequest) (*CheckMatchExistsResponse, error)
	mustEmbedUnimplementedCommunicationsServer()
}

// UnimplementedCommunicationsServer must be embedded to have forward compatible implementations.
type UnimplementedCommunicationsServer struct {
}

func (UnimplementedCommunicationsServer) AddReaction(context.Context, *AddReactionRequest) (*AddReactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddReaction not implemented")
}
func (UnimplementedCommunicationsServer) GetMatchList(context.Context, *GetMatchListRequest) (*GetMatchListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMatchList not implemented")
}
func (UnimplementedCommunicationsServer) GetReactionList(context.Context, *GetReactionListRequest) (*GetReactionListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReactionList not implemented")
}
func (UnimplementedCommunicationsServer) GetMatchTime(context.Context, *GetMatchTimeRequest) (*GetMatchTimeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMatchTime not implemented")
}
func (UnimplementedCommunicationsServer) GetMatchesBySearch(context.Context, *GetMatchesBySearchRequest) (*GetMatchesBySearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMatchesBySearch not implemented")
}
func (UnimplementedCommunicationsServer) UpdateOrCreateReaction(context.Context, *UpdateOrCreateReactionRequest) (*UpdateOrCreateReactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateOrCreateReaction not implemented")
}
func (UnimplementedCommunicationsServer) CheckMatchExists(context.Context, *CheckMatchExistsRequest) (*CheckMatchExistsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckMatchExists not implemented")
}
func (UnimplementedCommunicationsServer) mustEmbedUnimplementedCommunicationsServer() {}

// UnsafeCommunicationsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommunicationsServer will
// result in compilation errors.
type UnsafeCommunicationsServer interface {
	mustEmbedUnimplementedCommunicationsServer()
}

func RegisterCommunicationsServer(s grpc.ServiceRegistrar, srv CommunicationsServer) {
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

func _Communications_GetMatchTime_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMatchTimeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunicationsServer).GetMatchTime(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Communications_GetMatchTime_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunicationsServer).GetMatchTime(ctx, req.(*GetMatchTimeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Communications_GetMatchesBySearch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMatchesBySearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunicationsServer).GetMatchesBySearch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Communications_GetMatchesBySearch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunicationsServer).GetMatchesBySearch(ctx, req.(*GetMatchesBySearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Communications_UpdateOrCreateReaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateOrCreateReactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunicationsServer).UpdateOrCreateReaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Communications_UpdateOrCreateReaction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunicationsServer).UpdateOrCreateReaction(ctx, req.(*UpdateOrCreateReactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Communications_CheckMatchExists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckMatchExistsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunicationsServer).CheckMatchExists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Communications_CheckMatchExists_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunicationsServer).CheckMatchExists(ctx, req.(*CheckMatchExistsRequest))
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
		{
			MethodName: "GetMatchTime",
			Handler:    _Communications_GetMatchTime_Handler,
		},
		{
			MethodName: "GetMatchesBySearch",
			Handler:    _Communications_GetMatchesBySearch_Handler,
		},
		{
			MethodName: "UpdateOrCreateReaction",
			Handler:    _Communications_UpdateOrCreateReaction_Handler,
		},
		{
			MethodName: "CheckMatchExists",
			Handler:    _Communications_CheckMatchExists_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "communications.proto",
}

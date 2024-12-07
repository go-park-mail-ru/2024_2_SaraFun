// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.3
// source: survey.proto

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
	Survey_AddSurvey_FullMethodName      = "/survey.Survey/AddSurvey"
	Survey_GetSurveyInfo_FullMethodName  = "/survey.Survey/GetSurveyInfo"
	Survey_AddQuestion_FullMethodName    = "/survey.Survey/AddQuestion"
	Survey_UpdateQuestion_FullMethodName = "/survey.Survey/UpdateQuestion"
	Survey_DeleteQuestion_FullMethodName = "/survey.Survey/DeleteQuestion"
	Survey_GetQuestions_FullMethodName   = "/survey.Survey/GetQuestions"
)

// SurveyClient is the client API for Survey service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SurveyClient interface {
	AddSurvey(ctx context.Context, in *AddSurveyRequest, opts ...grpc.CallOption) (*AddSurveyResponse, error)
	GetSurveyInfo(ctx context.Context, in *GetSurveyInfoRequest, opts ...grpc.CallOption) (*GetSurveyInfoResponse, error)
	AddQuestion(ctx context.Context, in *AddQuestionRequest, opts ...grpc.CallOption) (*AddQuestionResponse, error)
	UpdateQuestion(ctx context.Context, in *UpdateQuestionRequest, opts ...grpc.CallOption) (*UpdateQuestionResponse, error)
	DeleteQuestion(ctx context.Context, in *DeleteQuestionRequest, opts ...grpc.CallOption) (*DeleteQuestionResponse, error)
	GetQuestions(ctx context.Context, in *GetQuestionsRequest, opts ...grpc.CallOption) (*GetQuestionResponse, error)
}

type surveyClient struct {
	cc grpc.ClientConnInterface
}

func NewSurveyClient(cc grpc.ClientConnInterface) SurveyClient {
	return &surveyClient{cc}
}

func (c *surveyClient) AddSurvey(ctx context.Context, in *AddSurveyRequest, opts ...grpc.CallOption) (*AddSurveyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddSurveyResponse)
	err := c.cc.Invoke(ctx, Survey_AddSurvey_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *surveyClient) GetSurveyInfo(ctx context.Context, in *GetSurveyInfoRequest, opts ...grpc.CallOption) (*GetSurveyInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetSurveyInfoResponse)
	err := c.cc.Invoke(ctx, Survey_GetSurveyInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *surveyClient) AddQuestion(ctx context.Context, in *AddQuestionRequest, opts ...grpc.CallOption) (*AddQuestionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddQuestionResponse)
	err := c.cc.Invoke(ctx, Survey_AddQuestion_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *surveyClient) UpdateQuestion(ctx context.Context, in *UpdateQuestionRequest, opts ...grpc.CallOption) (*UpdateQuestionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateQuestionResponse)
	err := c.cc.Invoke(ctx, Survey_UpdateQuestion_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *surveyClient) DeleteQuestion(ctx context.Context, in *DeleteQuestionRequest, opts ...grpc.CallOption) (*DeleteQuestionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteQuestionResponse)
	err := c.cc.Invoke(ctx, Survey_DeleteQuestion_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *surveyClient) GetQuestions(ctx context.Context, in *GetQuestionsRequest, opts ...grpc.CallOption) (*GetQuestionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetQuestionResponse)
	err := c.cc.Invoke(ctx, Survey_GetQuestions_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SurveyServer is the server API for Survey service.
// All implementations must embed UnimplementedSurveyServer
// for forward compatibility.
type SurveyServer interface {
	AddSurvey(context.Context, *AddSurveyRequest) (*AddSurveyResponse, error)
	GetSurveyInfo(context.Context, *GetSurveyInfoRequest) (*GetSurveyInfoResponse, error)
	AddQuestion(context.Context, *AddQuestionRequest) (*AddQuestionResponse, error)
	UpdateQuestion(context.Context, *UpdateQuestionRequest) (*UpdateQuestionResponse, error)
	DeleteQuestion(context.Context, *DeleteQuestionRequest) (*DeleteQuestionResponse, error)
	GetQuestions(context.Context, *GetQuestionsRequest) (*GetQuestionResponse, error)
	mustEmbedUnimplementedSurveyServer()
}

// UnimplementedSurveyServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSurveyServer struct{}

func (UnimplementedSurveyServer) AddSurvey(context.Context, *AddSurveyRequest) (*AddSurveyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddSurvey not implemented")
}
func (UnimplementedSurveyServer) GetSurveyInfo(context.Context, *GetSurveyInfoRequest) (*GetSurveyInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSurveyInfo not implemented")
}
func (UnimplementedSurveyServer) AddQuestion(context.Context, *AddQuestionRequest) (*AddQuestionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddQuestion not implemented")
}
func (UnimplementedSurveyServer) UpdateQuestion(context.Context, *UpdateQuestionRequest) (*UpdateQuestionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateQuestion not implemented")
}
func (UnimplementedSurveyServer) DeleteQuestion(context.Context, *DeleteQuestionRequest) (*DeleteQuestionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteQuestion not implemented")
}
func (UnimplementedSurveyServer) GetQuestions(context.Context, *GetQuestionsRequest) (*GetQuestionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetQuestions not implemented")
}
func (UnimplementedSurveyServer) mustEmbedUnimplementedSurveyServer() {}
func (UnimplementedSurveyServer) testEmbeddedByValue()                {}

// UnsafeSurveyServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SurveyServer will
// result in compilation errors.
type UnsafeSurveyServer interface {
	mustEmbedUnimplementedSurveyServer()
}

func RegisterSurveyServer(s grpc.ServiceRegistrar, srv SurveyServer) {
	// If the following call pancis, it indicates UnimplementedSurveyServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Survey_ServiceDesc, srv)
}

func _Survey_AddSurvey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddSurveyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SurveyServer).AddSurvey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Survey_AddSurvey_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SurveyServer).AddSurvey(ctx, req.(*AddSurveyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Survey_GetSurveyInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSurveyInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SurveyServer).GetSurveyInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Survey_GetSurveyInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SurveyServer).GetSurveyInfo(ctx, req.(*GetSurveyInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Survey_AddQuestion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddQuestionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SurveyServer).AddQuestion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Survey_AddQuestion_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SurveyServer).AddQuestion(ctx, req.(*AddQuestionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Survey_UpdateQuestion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateQuestionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SurveyServer).UpdateQuestion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Survey_UpdateQuestion_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SurveyServer).UpdateQuestion(ctx, req.(*UpdateQuestionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Survey_DeleteQuestion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteQuestionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SurveyServer).DeleteQuestion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Survey_DeleteQuestion_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SurveyServer).DeleteQuestion(ctx, req.(*DeleteQuestionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Survey_GetQuestions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetQuestionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SurveyServer).GetQuestions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Survey_GetQuestions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SurveyServer).GetQuestions(ctx, req.(*GetQuestionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Survey_ServiceDesc is the grpc.ServiceDesc for Survey service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Survey_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "survey.Survey",
	HandlerType: (*SurveyServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddSurvey",
			Handler:    _Survey_AddSurvey_Handler,
		},
		{
			MethodName: "GetSurveyInfo",
			Handler:    _Survey_GetSurveyInfo_Handler,
		},
		{
			MethodName: "AddQuestion",
			Handler:    _Survey_AddQuestion_Handler,
		},
		{
			MethodName: "UpdateQuestion",
			Handler:    _Survey_UpdateQuestion_Handler,
		},
		{
			MethodName: "DeleteQuestion",
			Handler:    _Survey_DeleteQuestion_Handler,
		},
		{
			MethodName: "GetQuestions",
			Handler:    _Survey_GetQuestions_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "survey.proto",
}

// Code generated by MockGen. DO NOT EDIT.
// Source: communications_grpc.pb.go

// Package mock_gen is a generated GoMock package.
package mock_gen

import (
	context "context"
	reflect "reflect"

	gen "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockCommunicationsClient is a mock of CommunicationsClient interface.
type MockCommunicationsClient struct {
	ctrl     *gomock.Controller
	recorder *MockCommunicationsClientMockRecorder
}

// MockCommunicationsClientMockRecorder is the mock recorder for MockCommunicationsClient.
type MockCommunicationsClientMockRecorder struct {
	mock *MockCommunicationsClient
}

// NewMockCommunicationsClient creates a new mock instance.
func NewMockCommunicationsClient(ctrl *gomock.Controller) *MockCommunicationsClient {
	mock := &MockCommunicationsClient{ctrl: ctrl}
	mock.recorder = &MockCommunicationsClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommunicationsClient) EXPECT() *MockCommunicationsClientMockRecorder {
	return m.recorder
}

// AddReaction mocks base method.
func (m *MockCommunicationsClient) AddReaction(ctx context.Context, in *gen.AddReactionRequest, opts ...grpc.CallOption) (*gen.AddReactionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddReaction", varargs...)
	ret0, _ := ret[0].(*gen.AddReactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddReaction indicates an expected call of AddReaction.
func (mr *MockCommunicationsClientMockRecorder) AddReaction(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddReaction", reflect.TypeOf((*MockCommunicationsClient)(nil).AddReaction), varargs...)
}

// CheckMatchExists mocks base method.
func (m *MockCommunicationsClient) CheckMatchExists(ctx context.Context, in *gen.CheckMatchExistsRequest, opts ...grpc.CallOption) (*gen.CheckMatchExistsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CheckMatchExists", varargs...)
	ret0, _ := ret[0].(*gen.CheckMatchExistsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckMatchExists indicates an expected call of CheckMatchExists.
func (mr *MockCommunicationsClientMockRecorder) CheckMatchExists(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckMatchExists", reflect.TypeOf((*MockCommunicationsClient)(nil).CheckMatchExists), varargs...)
}

// GetMatchList mocks base method.
func (m *MockCommunicationsClient) GetMatchList(ctx context.Context, in *gen.GetMatchListRequest, opts ...grpc.CallOption) (*gen.GetMatchListResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetMatchList", varargs...)
	ret0, _ := ret[0].(*gen.GetMatchListResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMatchList indicates an expected call of GetMatchList.
func (mr *MockCommunicationsClientMockRecorder) GetMatchList(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMatchList", reflect.TypeOf((*MockCommunicationsClient)(nil).GetMatchList), varargs...)
}

// GetMatchTime mocks base method.
func (m *MockCommunicationsClient) GetMatchTime(ctx context.Context, in *gen.GetMatchTimeRequest, opts ...grpc.CallOption) (*gen.GetMatchTimeResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetMatchTime", varargs...)
	ret0, _ := ret[0].(*gen.GetMatchTimeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMatchTime indicates an expected call of GetMatchTime.
func (mr *MockCommunicationsClientMockRecorder) GetMatchTime(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMatchTime", reflect.TypeOf((*MockCommunicationsClient)(nil).GetMatchTime), varargs...)
}

// GetMatchesBySearch mocks base method.
func (m *MockCommunicationsClient) GetMatchesBySearch(ctx context.Context, in *gen.GetMatchesBySearchRequest, opts ...grpc.CallOption) (*gen.GetMatchesBySearchResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetMatchesBySearch", varargs...)
	ret0, _ := ret[0].(*gen.GetMatchesBySearchResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMatchesBySearch indicates an expected call of GetMatchesBySearch.
func (mr *MockCommunicationsClientMockRecorder) GetMatchesBySearch(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMatchesBySearch", reflect.TypeOf((*MockCommunicationsClient)(nil).GetMatchesBySearch), varargs...)
}

// GetReactionList mocks base method.
func (m *MockCommunicationsClient) GetReactionList(ctx context.Context, in *gen.GetReactionListRequest, opts ...grpc.CallOption) (*gen.GetReactionListResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetReactionList", varargs...)
	ret0, _ := ret[0].(*gen.GetReactionListResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReactionList indicates an expected call of GetReactionList.
func (mr *MockCommunicationsClientMockRecorder) GetReactionList(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReactionList", reflect.TypeOf((*MockCommunicationsClient)(nil).GetReactionList), varargs...)
}

// UpdateOrCreateReaction mocks base method.
func (m *MockCommunicationsClient) UpdateOrCreateReaction(ctx context.Context, in *gen.UpdateOrCreateReactionRequest, opts ...grpc.CallOption) (*gen.UpdateOrCreateReactionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateOrCreateReaction", varargs...)
	ret0, _ := ret[0].(*gen.UpdateOrCreateReactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateOrCreateReaction indicates an expected call of UpdateOrCreateReaction.
func (mr *MockCommunicationsClientMockRecorder) UpdateOrCreateReaction(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrCreateReaction", reflect.TypeOf((*MockCommunicationsClient)(nil).UpdateOrCreateReaction), varargs...)
}

// MockCommunicationsServer is a mock of CommunicationsServer interface.
type MockCommunicationsServer struct {
	ctrl     *gomock.Controller
	recorder *MockCommunicationsServerMockRecorder
}

// MockCommunicationsServerMockRecorder is the mock recorder for MockCommunicationsServer.
type MockCommunicationsServerMockRecorder struct {
	mock *MockCommunicationsServer
}

// NewMockCommunicationsServer creates a new mock instance.
func NewMockCommunicationsServer(ctrl *gomock.Controller) *MockCommunicationsServer {
	mock := &MockCommunicationsServer{ctrl: ctrl}
	mock.recorder = &MockCommunicationsServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommunicationsServer) EXPECT() *MockCommunicationsServerMockRecorder {
	return m.recorder
}

// AddReaction mocks base method.
func (m *MockCommunicationsServer) AddReaction(arg0 context.Context, arg1 *gen.AddReactionRequest) (*gen.AddReactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddReaction", arg0, arg1)
	ret0, _ := ret[0].(*gen.AddReactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddReaction indicates an expected call of AddReaction.
func (mr *MockCommunicationsServerMockRecorder) AddReaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddReaction", reflect.TypeOf((*MockCommunicationsServer)(nil).AddReaction), arg0, arg1)
}

// CheckMatchExists mocks base method.
func (m *MockCommunicationsServer) CheckMatchExists(arg0 context.Context, arg1 *gen.CheckMatchExistsRequest) (*gen.CheckMatchExistsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckMatchExists", arg0, arg1)
	ret0, _ := ret[0].(*gen.CheckMatchExistsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckMatchExists indicates an expected call of CheckMatchExists.
func (mr *MockCommunicationsServerMockRecorder) CheckMatchExists(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckMatchExists", reflect.TypeOf((*MockCommunicationsServer)(nil).CheckMatchExists), arg0, arg1)
}

// GetMatchList mocks base method.
func (m *MockCommunicationsServer) GetMatchList(arg0 context.Context, arg1 *gen.GetMatchListRequest) (*gen.GetMatchListResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMatchList", arg0, arg1)
	ret0, _ := ret[0].(*gen.GetMatchListResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMatchList indicates an expected call of GetMatchList.
func (mr *MockCommunicationsServerMockRecorder) GetMatchList(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMatchList", reflect.TypeOf((*MockCommunicationsServer)(nil).GetMatchList), arg0, arg1)
}

// GetMatchTime mocks base method.
func (m *MockCommunicationsServer) GetMatchTime(arg0 context.Context, arg1 *gen.GetMatchTimeRequest) (*gen.GetMatchTimeResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMatchTime", arg0, arg1)
	ret0, _ := ret[0].(*gen.GetMatchTimeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMatchTime indicates an expected call of GetMatchTime.
func (mr *MockCommunicationsServerMockRecorder) GetMatchTime(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMatchTime", reflect.TypeOf((*MockCommunicationsServer)(nil).GetMatchTime), arg0, arg1)
}

// GetMatchesBySearch mocks base method.
func (m *MockCommunicationsServer) GetMatchesBySearch(arg0 context.Context, arg1 *gen.GetMatchesBySearchRequest) (*gen.GetMatchesBySearchResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMatchesBySearch", arg0, arg1)
	ret0, _ := ret[0].(*gen.GetMatchesBySearchResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMatchesBySearch indicates an expected call of GetMatchesBySearch.
func (mr *MockCommunicationsServerMockRecorder) GetMatchesBySearch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMatchesBySearch", reflect.TypeOf((*MockCommunicationsServer)(nil).GetMatchesBySearch), arg0, arg1)
}

// GetReactionList mocks base method.
func (m *MockCommunicationsServer) GetReactionList(arg0 context.Context, arg1 *gen.GetReactionListRequest) (*gen.GetReactionListResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReactionList", arg0, arg1)
	ret0, _ := ret[0].(*gen.GetReactionListResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReactionList indicates an expected call of GetReactionList.
func (mr *MockCommunicationsServerMockRecorder) GetReactionList(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReactionList", reflect.TypeOf((*MockCommunicationsServer)(nil).GetReactionList), arg0, arg1)
}

// UpdateOrCreateReaction mocks base method.
func (m *MockCommunicationsServer) UpdateOrCreateReaction(arg0 context.Context, arg1 *gen.UpdateOrCreateReactionRequest) (*gen.UpdateOrCreateReactionResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrCreateReaction", arg0, arg1)
	ret0, _ := ret[0].(*gen.UpdateOrCreateReactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateOrCreateReaction indicates an expected call of UpdateOrCreateReaction.
func (mr *MockCommunicationsServerMockRecorder) UpdateOrCreateReaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrCreateReaction", reflect.TypeOf((*MockCommunicationsServer)(nil).UpdateOrCreateReaction), arg0, arg1)
}

// mustEmbedUnimplementedCommunicationsServer mocks base method.
func (m *MockCommunicationsServer) mustEmbedUnimplementedCommunicationsServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedCommunicationsServer")
}

// mustEmbedUnimplementedCommunicationsServer indicates an expected call of mustEmbedUnimplementedCommunicationsServer.
func (mr *MockCommunicationsServerMockRecorder) mustEmbedUnimplementedCommunicationsServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedCommunicationsServer", reflect.TypeOf((*MockCommunicationsServer)(nil).mustEmbedUnimplementedCommunicationsServer))
}

// MockUnsafeCommunicationsServer is a mock of UnsafeCommunicationsServer interface.
type MockUnsafeCommunicationsServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeCommunicationsServerMockRecorder
}

// MockUnsafeCommunicationsServerMockRecorder is the mock recorder for MockUnsafeCommunicationsServer.
type MockUnsafeCommunicationsServerMockRecorder struct {
	mock *MockUnsafeCommunicationsServer
}

// NewMockUnsafeCommunicationsServer creates a new mock instance.
func NewMockUnsafeCommunicationsServer(ctrl *gomock.Controller) *MockUnsafeCommunicationsServer {
	mock := &MockUnsafeCommunicationsServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeCommunicationsServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeCommunicationsServer) EXPECT() *MockUnsafeCommunicationsServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedCommunicationsServer mocks base method.
func (m *MockUnsafeCommunicationsServer) mustEmbedUnimplementedCommunicationsServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedCommunicationsServer")
}

// mustEmbedUnimplementedCommunicationsServer indicates an expected call of mustEmbedUnimplementedCommunicationsServer.
func (mr *MockUnsafeCommunicationsServerMockRecorder) mustEmbedUnimplementedCommunicationsServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedCommunicationsServer", reflect.TypeOf((*MockUnsafeCommunicationsServer)(nil).mustEmbedUnimplementedCommunicationsServer))
}

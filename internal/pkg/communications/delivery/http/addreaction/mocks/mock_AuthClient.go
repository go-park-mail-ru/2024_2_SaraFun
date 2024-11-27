// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen (interfaces: AuthClient)

// Package addreaction_mocks is a generated GoMock package.
package addreaction_mocks

import (
	context "context"
	reflect "reflect"

	gen "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockAuthClient is a mock of AuthClient interface.
type MockAuthClient struct {
	ctrl     *gomock.Controller
	recorder *MockAuthClientMockRecorder
}

// MockAuthClientMockRecorder is the mock recorder for MockAuthClient.
type MockAuthClientMockRecorder struct {
	mock *MockAuthClient
}

// NewMockAuthClient creates a new mock instance.
func NewMockAuthClient(ctrl *gomock.Controller) *MockAuthClient {
	mock := &MockAuthClient{ctrl: ctrl}
	mock.recorder = &MockAuthClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthClient) EXPECT() *MockAuthClientMockRecorder {
	return m.recorder
}

// CheckSession mocks base method.
func (m *MockAuthClient) CheckSession(arg0 context.Context, arg1 *gen.CheckSessionRequest, arg2 ...grpc.CallOption) (*gen.CheckSessionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CheckSession", varargs...)
	ret0, _ := ret[0].(*gen.CheckSessionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckSession indicates an expected call of CheckSession.
func (mr *MockAuthClientMockRecorder) CheckSession(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckSession", reflect.TypeOf((*MockAuthClient)(nil).CheckSession), varargs...)
}

// CreateSession mocks base method.
func (m *MockAuthClient) CreateSession(arg0 context.Context, arg1 *gen.CreateSessionRequest, arg2 ...grpc.CallOption) (*gen.CreateSessionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateSession", varargs...)
	ret0, _ := ret[0].(*gen.CreateSessionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockAuthClientMockRecorder) CreateSession(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockAuthClient)(nil).CreateSession), varargs...)
}

// DeleteSession mocks base method.
func (m *MockAuthClient) DeleteSession(arg0 context.Context, arg1 *gen.DeleteSessionRequest, arg2 ...grpc.CallOption) (*gen.DeleteSessionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteSession", varargs...)
	ret0, _ := ret[0].(*gen.DeleteSessionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSession indicates an expected call of DeleteSession.
func (mr *MockAuthClientMockRecorder) DeleteSession(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSession", reflect.TypeOf((*MockAuthClient)(nil).DeleteSession), varargs...)
}

// GetUserIDBySessionID mocks base method.
func (m *MockAuthClient) GetUserIDBySessionID(arg0 context.Context, arg1 *gen.GetUserIDBySessionIDRequest, arg2 ...grpc.CallOption) (*gen.GetUserIDBYSessionIDResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetUserIDBySessionID", varargs...)
	ret0, _ := ret[0].(*gen.GetUserIDBYSessionIDResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserIDBySessionID indicates an expected call of GetUserIDBySessionID.
func (mr *MockAuthClientMockRecorder) GetUserIDBySessionID(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserIDBySessionID", reflect.TypeOf((*MockAuthClient)(nil).GetUserIDBySessionID), varargs...)
}

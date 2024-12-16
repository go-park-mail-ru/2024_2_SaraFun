// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/http/logout (interfaces: SessionService)

// Package sign_up_mocks is a generated GoMock package.
package sign_up_mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSessionService is a mock of SessionService interface.
type MockSessionService struct {
	ctrl     *gomock.Controller
	recorder *MockSessionServiceMockRecorder
}

// MockSessionServiceMockRecorder is the mock recorder for MockSessionService.
type MockSessionServiceMockRecorder struct {
	mock *MockSessionService
}

// NewMockSessionService creates a new mock instance.
func NewMockSessionService(ctrl *gomock.Controller) *MockSessionService {
	mock := &MockSessionService{ctrl: ctrl}
	mock.recorder = &MockSessionServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionService) EXPECT() *MockSessionServiceMockRecorder {
	return m.recorder
}

// DeleteSession mocks base method.
func (m *MockSessionService) DeleteSession(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSession", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSession indicates an expected call of DeleteSession.
func (mr *MockSessionServiceMockRecorder) DeleteSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSession", reflect.TypeOf((*MockSessionService)(nil).DeleteSession), arg0, arg1)
}

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/websockets/delivery/setConnection (interfaces: UseCase)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	websocket "github.com/gorilla/websocket"
)

// MockUseCase is a mock of UseCase interface.
type MockUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUseCaseMockRecorder
}

// MockUseCaseMockRecorder is the mock recorder for MockUseCase.
type MockUseCaseMockRecorder struct {
	mock *MockUseCase
}

// NewMockUseCase creates a new mock instance.
func NewMockUseCase(ctrl *gomock.Controller) *MockUseCase {
	mock := &MockUseCase{ctrl: ctrl}
	mock.recorder = &MockUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUseCase) EXPECT() *MockUseCaseMockRecorder {
	return m.recorder
}

// AddConnection mocks base method.
func (m *MockUseCase) AddConnection(arg0 context.Context, arg1 *websocket.Conn, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddConnection", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddConnection indicates an expected call of AddConnection.
func (mr *MockUseCaseMockRecorder) AddConnection(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddConnection", reflect.TypeOf((*MockUseCase)(nil).AddConnection), arg0, arg1, arg2)
}

// DeleteConnection mocks base method.
func (m *MockUseCase) DeleteConnection(arg0 context.Context, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteConnection", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteConnection indicates an expected call of DeleteConnection.
func (mr *MockUseCaseMockRecorder) DeleteConnection(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteConnection", reflect.TypeOf((*MockUseCase)(nil).DeleteConnection), arg0, arg1)
}

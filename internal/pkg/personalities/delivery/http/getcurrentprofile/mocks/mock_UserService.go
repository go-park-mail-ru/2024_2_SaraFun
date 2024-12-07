// Code generated by MockGen. DO NOT EDIT.
// Source: sparkit/internal/handlers/getcurrentprofile (interfaces: UserService)

// Package sign_up_mocks is a generated GoMock package.
package sign_up_mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserService is a mock of UserService interface.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// GetProfileIdByUserId mocks base method.
func (m *MockUserService) GetProfileIdByUserId(arg0 context.Context, arg1 int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfileIdByUserId", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfileIdByUserId indicates an expected call of GetProfileIdByUserId.
func (mr *MockUserServiceMockRecorder) GetProfileIdByUserId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfileIdByUserId", reflect.TypeOf((*MockUserService)(nil).GetProfileIdByUserId), arg0, arg1)
}

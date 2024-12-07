// Code generated by MockGen. DO NOT EDIT.
// Source: sparkit/internal/handlers/getuserlist (interfaces: ProfileService)

// Package getuserlist_mocks is a generated GoMock package.
package getuserlist_mocks

import (
	context "context"
	models "github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockProfileService is a mock of ProfileService interface.
type MockProfileService struct {
	ctrl     *gomock.Controller
	recorder *MockProfileServiceMockRecorder
}

// MockProfileServiceMockRecorder is the mock recorder for MockProfileService.
type MockProfileServiceMockRecorder struct {
	mock *MockProfileService
}

// NewMockProfileService creates a new mock instance.
func NewMockProfileService(ctrl *gomock.Controller) *MockProfileService {
	mock := &MockProfileService{ctrl: ctrl}
	mock.recorder = &MockProfileServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProfileService) EXPECT() *MockProfileServiceMockRecorder {
	return m.recorder
}

// GetProfile mocks base method.
func (m *MockProfileService) GetProfile(arg0 context.Context, arg1 int) (models.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfile", arg0, arg1)
	ret0, _ := ret[0].(models.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfile indicates an expected call of GetProfile.
func (mr *MockProfileServiceMockRecorder) GetProfile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfile", reflect.TypeOf((*MockProfileService)(nil).GetProfile), arg0, arg1)
}

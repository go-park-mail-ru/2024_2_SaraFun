// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc (interfaces: UserUsecase)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	models "github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockUserUsecase is a mock of UserUsecase interface.
type MockUserUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockUserUsecaseMockRecorder
}

// MockUserUsecaseMockRecorder is the mock recorder for MockUserUsecase.
type MockUserUsecaseMockRecorder struct {
	mock *MockUserUsecase
}

// NewMockUserUsecase creates a new mock instance.
func NewMockUserUsecase(ctrl *gomock.Controller) *MockUserUsecase {
	mock := &MockUserUsecase{ctrl: ctrl}
	mock.recorder = &MockUserUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserUsecase) EXPECT() *MockUserUsecaseMockRecorder {
	return m.recorder
}

// ChangePassword mocks base method.
func (m *MockUserUsecase) ChangePassword(arg0 context.Context, arg1 int, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangePassword", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangePassword indicates an expected call of ChangePassword.
func (mr *MockUserUsecaseMockRecorder) ChangePassword(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangePassword", reflect.TypeOf((*MockUserUsecase)(nil).ChangePassword), arg0, arg1, arg2)
}

// CheckPassword mocks base method.
func (m *MockUserUsecase) CheckPassword(arg0 context.Context, arg1, arg2 string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckPassword", arg0, arg1, arg2)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckPassword indicates an expected call of CheckPassword.
func (mr *MockUserUsecaseMockRecorder) CheckPassword(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckPassword", reflect.TypeOf((*MockUserUsecase)(nil).CheckPassword), arg0, arg1, arg2)
}

// CheckUsernameExists mocks base method.
func (m *MockUserUsecase) CheckUsernameExists(arg0 context.Context, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUsernameExists", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUsernameExists indicates an expected call of CheckUsernameExists.
func (mr *MockUserUsecaseMockRecorder) CheckUsernameExists(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUsernameExists", reflect.TypeOf((*MockUserUsecase)(nil).CheckUsernameExists), arg0, arg1)
}

// GetFeedList mocks base method.
func (m *MockUserUsecase) GetFeedList(arg0 context.Context, arg1 int, arg2 []int) ([]models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFeedList", arg0, arg1, arg2)
	ret0, _ := ret[0].([]models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFeedList indicates an expected call of GetFeedList.
func (mr *MockUserUsecaseMockRecorder) GetFeedList(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeedList", reflect.TypeOf((*MockUserUsecase)(nil).GetFeedList), arg0, arg1, arg2)
}

// GetProfileIdByUserId mocks base method.
func (m *MockUserUsecase) GetProfileIdByUserId(arg0 context.Context, arg1 int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfileIdByUserId", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfileIdByUserId indicates an expected call of GetProfileIdByUserId.
func (mr *MockUserUsecaseMockRecorder) GetProfileIdByUserId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfileIdByUserId", reflect.TypeOf((*MockUserUsecase)(nil).GetProfileIdByUserId), arg0, arg1)
}

// GetUserIdByUsername mocks base method.
func (m *MockUserUsecase) GetUserIdByUsername(arg0 context.Context, arg1 string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserIdByUsername", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserIdByUsername indicates an expected call of GetUserIdByUsername.
func (mr *MockUserUsecaseMockRecorder) GetUserIdByUsername(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserIdByUsername", reflect.TypeOf((*MockUserUsecase)(nil).GetUserIdByUsername), arg0, arg1)
}

// GetUsernameByUserId mocks base method.
func (m *MockUserUsecase) GetUsernameByUserId(arg0 context.Context, arg1 int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsernameByUserId", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsernameByUserId indicates an expected call of GetUsernameByUserId.
func (mr *MockUserUsecaseMockRecorder) GetUsernameByUserId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsernameByUserId", reflect.TypeOf((*MockUserUsecase)(nil).GetUsernameByUserId), arg0, arg1)
}

// RegisterUser mocks base method.
func (m *MockUserUsecase) RegisterUser(arg0 context.Context, arg1 models.User) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUser", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterUser indicates an expected call of RegisterUser.
func (mr *MockUserUsecaseMockRecorder) RegisterUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUser", reflect.TypeOf((*MockUserUsecase)(nil).RegisterUser), arg0, arg1)
}

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc (interfaces: UseCase)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	models "github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	gomock "github.com/golang/mock/gomock"
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

// AddBalance mocks base method.
func (m *MockUseCase) AddBalance(arg0 context.Context, arg1, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddBalance", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddBalance indicates an expected call of AddBalance.
func (mr *MockUseCaseMockRecorder) AddBalance(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBalance", reflect.TypeOf((*MockUseCase)(nil).AddBalance), arg0, arg1, arg2)
}

// AddDailyLikesCount mocks base method.
func (m *MockUseCase) AddDailyLikesCount(arg0 context.Context, arg1, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDailyLikesCount", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddDailyLikesCount indicates an expected call of AddDailyLikesCount.
func (mr *MockUseCaseMockRecorder) AddDailyLikesCount(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDailyLikesCount", reflect.TypeOf((*MockUseCase)(nil).AddDailyLikesCount), arg0, arg1, arg2)
}

// AddPurchasedLikesCount mocks base method.
func (m *MockUseCase) AddPurchasedLikesCount(arg0 context.Context, arg1, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPurchasedLikesCount", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddPurchasedLikesCount indicates an expected call of AddPurchasedLikesCount.
func (mr *MockUseCaseMockRecorder) AddPurchasedLikesCount(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPurchasedLikesCount", reflect.TypeOf((*MockUseCase)(nil).AddPurchasedLikesCount), arg0, arg1, arg2)
}

// ChangeBalance mocks base method.
func (m *MockUseCase) ChangeBalance(arg0 context.Context, arg1, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeBalance", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeBalance indicates an expected call of ChangeBalance.
func (mr *MockUseCaseMockRecorder) ChangeBalance(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeBalance", reflect.TypeOf((*MockUseCase)(nil).ChangeBalance), arg0, arg1, arg2)
}

// ChangeDailyLikeCount mocks base method.
func (m *MockUseCase) ChangeDailyLikeCount(arg0 context.Context, arg1, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeDailyLikeCount", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeDailyLikeCount indicates an expected call of ChangeDailyLikeCount.
func (mr *MockUseCaseMockRecorder) ChangeDailyLikeCount(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeDailyLikeCount", reflect.TypeOf((*MockUseCase)(nil).ChangeDailyLikeCount), arg0, arg1, arg2)
}

// ChangePurchasedLikeCount mocks base method.
func (m *MockUseCase) ChangePurchasedLikeCount(arg0 context.Context, arg1, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangePurchasedLikeCount", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangePurchasedLikeCount indicates an expected call of ChangePurchasedLikeCount.
func (mr *MockUseCaseMockRecorder) ChangePurchasedLikeCount(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangePurchasedLikeCount", reflect.TypeOf((*MockUseCase)(nil).ChangePurchasedLikeCount), arg0, arg1, arg2)
}

// CheckBalance mocks base method.
func (m *MockUseCase) CheckBalance(arg0 context.Context, arg1, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckBalance", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckBalance indicates an expected call of CheckBalance.
func (mr *MockUseCaseMockRecorder) CheckBalance(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckBalance", reflect.TypeOf((*MockUseCase)(nil).CheckBalance), arg0, arg1, arg2)
}

// CreateProduct mocks base method.
func (m *MockUseCase) CreateProduct(arg0 context.Context, arg1 models.Product) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProduct", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProduct indicates an expected call of CreateProduct.
func (mr *MockUseCaseMockRecorder) CreateProduct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProduct", reflect.TypeOf((*MockUseCase)(nil).CreateProduct), arg0, arg1)
}

// GetBalance mocks base method.
func (m *MockUseCase) GetBalance(arg0 context.Context, arg1 int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalance", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBalance indicates an expected call of GetBalance.
func (mr *MockUseCaseMockRecorder) GetBalance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalance", reflect.TypeOf((*MockUseCase)(nil).GetBalance), arg0, arg1)
}

// GetDailyLikesCount mocks base method.
func (m *MockUseCase) GetDailyLikesCount(arg0 context.Context, arg1 int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDailyLikesCount", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDailyLikesCount indicates an expected call of GetDailyLikesCount.
func (mr *MockUseCaseMockRecorder) GetDailyLikesCount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDailyLikesCount", reflect.TypeOf((*MockUseCase)(nil).GetDailyLikesCount), arg0, arg1)
}

// GetProduct mocks base method.
func (m *MockUseCase) GetProduct(arg0 context.Context, arg1 string) (models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProduct", arg0, arg1)
	ret0, _ := ret[0].(models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProduct indicates an expected call of GetProduct.
func (mr *MockUseCaseMockRecorder) GetProduct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProduct", reflect.TypeOf((*MockUseCase)(nil).GetProduct), arg0, arg1)
}

// GetProducts mocks base method.
func (m *MockUseCase) GetProducts(arg0 context.Context) ([]models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProducts", arg0)
	ret0, _ := ret[0].([]models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProducts indicates an expected call of GetProducts.
func (mr *MockUseCaseMockRecorder) GetProducts(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProducts", reflect.TypeOf((*MockUseCase)(nil).GetProducts), arg0)
}

// GetPurchasedLikesCount mocks base method.
func (m *MockUseCase) GetPurchasedLikesCount(arg0 context.Context, arg1 int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPurchasedLikesCount", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPurchasedLikesCount indicates an expected call of GetPurchasedLikesCount.
func (mr *MockUseCaseMockRecorder) GetPurchasedLikesCount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPurchasedLikesCount", reflect.TypeOf((*MockUseCase)(nil).GetPurchasedLikesCount), arg0, arg1)
}

// SetDailyLikeCountToAll mocks base method.
func (m *MockUseCase) SetDailyLikeCountToAll(arg0 context.Context, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetDailyLikeCountToAll", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetDailyLikeCountToAll indicates an expected call of SetDailyLikeCountToAll.
func (mr *MockUseCaseMockRecorder) SetDailyLikeCountToAll(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDailyLikeCountToAll", reflect.TypeOf((*MockUseCase)(nil).SetDailyLikeCountToAll), arg0, arg1)
}
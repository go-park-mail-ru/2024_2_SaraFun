// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/usecase (interfaces: Repository)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	models "github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// AddActivity mocks base method.
func (m *MockRepository) AddActivity(arg0 context.Context, arg1 models.Activity) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddActivity", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddActivity indicates an expected call of AddActivity.
func (mr *MockRepositoryMockRecorder) AddActivity(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddActivity", reflect.TypeOf((*MockRepository)(nil).AddActivity), arg0, arg1)
}

// AddAward mocks base method.
func (m *MockRepository) AddAward(arg0 context.Context, arg1 models.Award) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddAward", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddAward indicates an expected call of AddAward.
func (mr *MockRepositoryMockRecorder) AddAward(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAward", reflect.TypeOf((*MockRepository)(nil).AddAward), arg0, arg1)
}

// AddBalance mocks base method.
func (m *MockRepository) AddBalance(arg0 context.Context, arg1, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddBalance", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddBalance indicates an expected call of AddBalance.
func (mr *MockRepositoryMockRecorder) AddBalance(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBalance", reflect.TypeOf((*MockRepository)(nil).AddBalance), arg0, arg1, arg2)
}

// AddDailyLikeCount mocks base method.
func (m *MockRepository) AddDailyLikeCount(arg0 context.Context, arg1, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDailyLikeCount", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddDailyLikeCount indicates an expected call of AddDailyLikeCount.
func (mr *MockRepositoryMockRecorder) AddDailyLikeCount(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDailyLikeCount", reflect.TypeOf((*MockRepository)(nil).AddDailyLikeCount), arg0, arg1, arg2)
}

// AddPurchasedLikeCount mocks base method.
func (m *MockRepository) AddPurchasedLikeCount(arg0 context.Context, arg1, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPurchasedLikeCount", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddPurchasedLikeCount indicates an expected call of AddPurchasedLikeCount.
func (mr *MockRepositoryMockRecorder) AddPurchasedLikeCount(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPurchasedLikeCount", reflect.TypeOf((*MockRepository)(nil).AddPurchasedLikeCount), arg0, arg1, arg2)
}

// ChangeBalance mocks base method.
func (m *MockRepository) ChangeBalance(arg0 context.Context, arg1, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeBalance", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeBalance indicates an expected call of ChangeBalance.
func (mr *MockRepositoryMockRecorder) ChangeBalance(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeBalance", reflect.TypeOf((*MockRepository)(nil).ChangeBalance), arg0, arg1, arg2)
}

// ChangeDailyLikeCount mocks base method.
func (m *MockRepository) ChangeDailyLikeCount(arg0 context.Context, arg1, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeDailyLikeCount", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeDailyLikeCount indicates an expected call of ChangeDailyLikeCount.
func (mr *MockRepositoryMockRecorder) ChangeDailyLikeCount(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeDailyLikeCount", reflect.TypeOf((*MockRepository)(nil).ChangeDailyLikeCount), arg0, arg1, arg2)
}

// ChangePurchasedLikeCount mocks base method.
func (m *MockRepository) ChangePurchasedLikeCount(arg0 context.Context, arg1, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangePurchasedLikeCount", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangePurchasedLikeCount indicates an expected call of ChangePurchasedLikeCount.
func (mr *MockRepositoryMockRecorder) ChangePurchasedLikeCount(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangePurchasedLikeCount", reflect.TypeOf((*MockRepository)(nil).ChangePurchasedLikeCount), arg0, arg1, arg2)
}

// CreateProduct mocks base method.
func (m *MockRepository) CreateProduct(arg0 context.Context, arg1 models.Product) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProduct", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProduct indicates an expected call of CreateProduct.
func (mr *MockRepositoryMockRecorder) CreateProduct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProduct", reflect.TypeOf((*MockRepository)(nil).CreateProduct), arg0, arg1)
}

// GetActivity mocks base method.
func (m *MockRepository) GetActivity(arg0 context.Context, arg1 int) (models.Activity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActivity", arg0, arg1)
	ret0, _ := ret[0].(models.Activity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActivity indicates an expected call of GetActivity.
func (mr *MockRepositoryMockRecorder) GetActivity(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActivity", reflect.TypeOf((*MockRepository)(nil).GetActivity), arg0, arg1)
}

// GetActivityDay mocks base method.
func (m *MockRepository) GetActivityDay(arg0 context.Context, arg1 int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActivityDay", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActivityDay indicates an expected call of GetActivityDay.
func (mr *MockRepositoryMockRecorder) GetActivityDay(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActivityDay", reflect.TypeOf((*MockRepository)(nil).GetActivityDay), arg0, arg1)
}

// GetAwardByDayNumber mocks base method.
func (m *MockRepository) GetAwardByDayNumber(arg0 context.Context, arg1 int) (models.Award, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAwardByDayNumber", arg0, arg1)
	ret0, _ := ret[0].(models.Award)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAwardByDayNumber indicates an expected call of GetAwardByDayNumber.
func (mr *MockRepositoryMockRecorder) GetAwardByDayNumber(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAwardByDayNumber", reflect.TypeOf((*MockRepository)(nil).GetAwardByDayNumber), arg0, arg1)
}

// GetAwards mocks base method.
func (m *MockRepository) GetAwards(arg0 context.Context) ([]models.Award, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAwards", arg0)
	ret0, _ := ret[0].([]models.Award)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAwards indicates an expected call of GetAwards.
func (mr *MockRepositoryMockRecorder) GetAwards(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAwards", reflect.TypeOf((*MockRepository)(nil).GetAwards), arg0)
}

// GetBalance mocks base method.
func (m *MockRepository) GetBalance(arg0 context.Context, arg1 int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalance", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBalance indicates an expected call of GetBalance.
func (mr *MockRepositoryMockRecorder) GetBalance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalance", reflect.TypeOf((*MockRepository)(nil).GetBalance), arg0, arg1)
}

// GetDailyLikesCount mocks base method.
func (m *MockRepository) GetDailyLikesCount(arg0 context.Context, arg1 int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDailyLikesCount", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDailyLikesCount indicates an expected call of GetDailyLikesCount.
func (mr *MockRepositoryMockRecorder) GetDailyLikesCount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDailyLikesCount", reflect.TypeOf((*MockRepository)(nil).GetDailyLikesCount), arg0, arg1)
}

// GetProduct mocks base method.
func (m *MockRepository) GetProduct(arg0 context.Context, arg1 string) (models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProduct", arg0, arg1)
	ret0, _ := ret[0].(models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProduct indicates an expected call of GetProduct.
func (mr *MockRepositoryMockRecorder) GetProduct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProduct", reflect.TypeOf((*MockRepository)(nil).GetProduct), arg0, arg1)
}

// GetProducts mocks base method.
func (m *MockRepository) GetProducts(arg0 context.Context) ([]models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProducts", arg0)
	ret0, _ := ret[0].([]models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProducts indicates an expected call of GetProducts.
func (mr *MockRepositoryMockRecorder) GetProducts(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProducts", reflect.TypeOf((*MockRepository)(nil).GetProducts), arg0)
}

// GetPurchasedLikesCount mocks base method.
func (m *MockRepository) GetPurchasedLikesCount(arg0 context.Context, arg1 int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPurchasedLikesCount", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPurchasedLikesCount indicates an expected call of GetPurchasedLikesCount.
func (mr *MockRepositoryMockRecorder) GetPurchasedLikesCount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPurchasedLikesCount", reflect.TypeOf((*MockRepository)(nil).GetPurchasedLikesCount), arg0, arg1)
}

// SetBalance mocks base method.
func (m *MockRepository) SetBalance(arg0 context.Context, arg1, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetBalance", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetBalance indicates an expected call of SetBalance.
func (mr *MockRepositoryMockRecorder) SetBalance(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBalance", reflect.TypeOf((*MockRepository)(nil).SetBalance), arg0, arg1, arg2)
}

// SetDailyLikesCount mocks base method.
func (m *MockRepository) SetDailyLikesCount(arg0 context.Context, arg1, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetDailyLikesCount", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetDailyLikesCount indicates an expected call of SetDailyLikesCount.
func (mr *MockRepositoryMockRecorder) SetDailyLikesCount(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDailyLikesCount", reflect.TypeOf((*MockRepository)(nil).SetDailyLikesCount), arg0, arg1, arg2)
}

// SetDailyLikesCountToAll mocks base method.
func (m *MockRepository) SetDailyLikesCountToAll(arg0 context.Context, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetDailyLikesCountToAll", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetDailyLikesCountToAll indicates an expected call of SetDailyLikesCountToAll.
func (mr *MockRepositoryMockRecorder) SetDailyLikesCountToAll(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDailyLikesCountToAll", reflect.TypeOf((*MockRepository)(nil).SetDailyLikesCountToAll), arg0, arg1)
}

// SetPurchasedLikesCount mocks base method.
func (m *MockRepository) SetPurchasedLikesCount(arg0 context.Context, arg1, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetPurchasedLikesCount", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetPurchasedLikesCount indicates an expected call of SetPurchasedLikesCount.
func (mr *MockRepositoryMockRecorder) SetPurchasedLikesCount(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPurchasedLikesCount", reflect.TypeOf((*MockRepository)(nil).SetPurchasedLikesCount), arg0, arg1, arg2)
}

// UpdateActivity mocks base method.
func (m *MockRepository) UpdateActivity(arg0 context.Context, arg1 int, arg2 models.Activity) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateActivity", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateActivity indicates an expected call of UpdateActivity.
func (mr *MockRepositoryMockRecorder) UpdateActivity(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateActivity", reflect.TypeOf((*MockRepository)(nil).UpdateActivity), arg0, arg1, arg2)
}

// UpdateProduct mocks base method.
func (m *MockRepository) UpdateProduct(arg0 context.Context, arg1 string, arg2 models.Product) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProduct", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProduct indicates an expected call of UpdateProduct.
func (mr *MockRepositoryMockRecorder) UpdateProduct(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProduct", reflect.TypeOf((*MockRepository)(nil).UpdateProduct), arg0, arg1, arg2)
}

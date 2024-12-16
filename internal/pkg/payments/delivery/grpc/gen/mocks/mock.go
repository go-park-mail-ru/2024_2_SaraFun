// Code generated by MockGen. DO NOT EDIT.
// Source: payments_grpc.pb.go

// Package mock_gen is a generated GoMock package.
package mock_gen

import (
	context "context"
	reflect "reflect"

	gen "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc/gen"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockPaymentClient is a mock of PaymentClient interface.
type MockPaymentClient struct {
	ctrl     *gomock.Controller
	recorder *MockPaymentClientMockRecorder
}

// MockPaymentClientMockRecorder is the mock recorder for MockPaymentClient.
type MockPaymentClientMockRecorder struct {
	mock *MockPaymentClient
}

// NewMockPaymentClient creates a new mock instance.
func NewMockPaymentClient(ctrl *gomock.Controller) *MockPaymentClient {
	mock := &MockPaymentClient{ctrl: ctrl}
	mock.recorder = &MockPaymentClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPaymentClient) EXPECT() *MockPaymentClientMockRecorder {
	return m.recorder
}

// BuyLikes mocks base method.
func (m *MockPaymentClient) BuyLikes(ctx context.Context, in *gen.BuyLikesRequest, opts ...grpc.CallOption) (*gen.BuyLikesResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "BuyLikes", varargs...)
	ret0, _ := ret[0].(*gen.BuyLikesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BuyLikes indicates an expected call of BuyLikes.
func (mr *MockPaymentClientMockRecorder) BuyLikes(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuyLikes", reflect.TypeOf((*MockPaymentClient)(nil).BuyLikes), varargs...)
}

// ChangeBalance mocks base method.
func (m *MockPaymentClient) ChangeBalance(ctx context.Context, in *gen.ChangeBalanceRequest, opts ...grpc.CallOption) (*gen.ChangeBalanceResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ChangeBalance", varargs...)
	ret0, _ := ret[0].(*gen.ChangeBalanceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChangeBalance indicates an expected call of ChangeBalance.
func (mr *MockPaymentClientMockRecorder) ChangeBalance(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeBalance", reflect.TypeOf((*MockPaymentClient)(nil).ChangeBalance), varargs...)
}

// ChangePurchasedLikesBalance mocks base method.
func (m *MockPaymentClient) ChangePurchasedLikesBalance(ctx context.Context, in *gen.ChangePurchasedLikesBalanceRequest, opts ...grpc.CallOption) (*gen.ChangePurchasedLikesBalanceResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ChangePurchasedLikesBalance", varargs...)
	ret0, _ := ret[0].(*gen.ChangePurchasedLikesBalanceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChangePurchasedLikesBalance indicates an expected call of ChangePurchasedLikesBalance.
func (mr *MockPaymentClientMockRecorder) ChangePurchasedLikesBalance(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangePurchasedLikesBalance", reflect.TypeOf((*MockPaymentClient)(nil).ChangePurchasedLikesBalance), varargs...)
}

// CheckAndSpendLike mocks base method.
func (m *MockPaymentClient) CheckAndSpendLike(ctx context.Context, in *gen.CheckAndSpendLikeRequest, opts ...grpc.CallOption) (*gen.CheckAndSpendLikeResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CheckAndSpendLike", varargs...)
	ret0, _ := ret[0].(*gen.CheckAndSpendLikeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckAndSpendLike indicates an expected call of CheckAndSpendLike.
func (mr *MockPaymentClientMockRecorder) CheckAndSpendLike(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAndSpendLike", reflect.TypeOf((*MockPaymentClient)(nil).CheckAndSpendLike), varargs...)
}

// CreateBalances mocks base method.
func (m *MockPaymentClient) CreateBalances(ctx context.Context, in *gen.CreateBalancesRequest, opts ...grpc.CallOption) (*gen.CreateBalancesResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateBalances", varargs...)
	ret0, _ := ret[0].(*gen.CreateBalancesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBalances indicates an expected call of CreateBalances.
func (mr *MockPaymentClientMockRecorder) CreateBalances(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBalances", reflect.TypeOf((*MockPaymentClient)(nil).CreateBalances), varargs...)
}

// CreateProduct mocks base method.
func (m *MockPaymentClient) CreateProduct(ctx context.Context, in *gen.CreateProductRequest, opts ...grpc.CallOption) (*gen.CreateProductResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateProduct", varargs...)
	ret0, _ := ret[0].(*gen.CreateProductResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProduct indicates an expected call of CreateProduct.
func (mr *MockPaymentClientMockRecorder) CreateProduct(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProduct", reflect.TypeOf((*MockPaymentClient)(nil).CreateProduct), varargs...)
}

// GetAllBalance mocks base method.
func (m *MockPaymentClient) GetAllBalance(ctx context.Context, in *gen.GetAllBalanceRequest, opts ...grpc.CallOption) (*gen.GetAllBalanceResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAllBalance", varargs...)
	ret0, _ := ret[0].(*gen.GetAllBalanceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllBalance indicates an expected call of GetAllBalance.
func (mr *MockPaymentClientMockRecorder) GetAllBalance(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllBalance", reflect.TypeOf((*MockPaymentClient)(nil).GetAllBalance), varargs...)
}

// GetBalance mocks base method.
func (m *MockPaymentClient) GetBalance(ctx context.Context, in *gen.GetBalanceRequest, opts ...grpc.CallOption) (*gen.GetBalanceResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetBalance", varargs...)
	ret0, _ := ret[0].(*gen.GetBalanceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBalance indicates an expected call of GetBalance.
func (mr *MockPaymentClientMockRecorder) GetBalance(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalance", reflect.TypeOf((*MockPaymentClient)(nil).GetBalance), varargs...)
}

// GetDailyLikeBalance mocks base method.
func (m *MockPaymentClient) GetDailyLikeBalance(ctx context.Context, in *gen.GetDailyLikeBalanceRequest, opts ...grpc.CallOption) (*gen.GetDailyLikeBalanceResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetDailyLikeBalance", varargs...)
	ret0, _ := ret[0].(*gen.GetDailyLikeBalanceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDailyLikeBalance indicates an expected call of GetDailyLikeBalance.
func (mr *MockPaymentClientMockRecorder) GetDailyLikeBalance(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDailyLikeBalance", reflect.TypeOf((*MockPaymentClient)(nil).GetDailyLikeBalance), varargs...)
}

// GetProducts mocks base method.
func (m *MockPaymentClient) GetProducts(ctx context.Context, in *gen.GetProductsRequest, opts ...grpc.CallOption) (*gen.GetProductsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetProducts", varargs...)
	ret0, _ := ret[0].(*gen.GetProductsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProducts indicates an expected call of GetProducts.
func (mr *MockPaymentClientMockRecorder) GetProducts(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProducts", reflect.TypeOf((*MockPaymentClient)(nil).GetProducts), varargs...)
}

// GetPurchasedLikeBalance mocks base method.
func (m *MockPaymentClient) GetPurchasedLikeBalance(ctx context.Context, in *gen.GetPurchasedLikeBalanceRequest, opts ...grpc.CallOption) (*gen.GetPurchasedLikeBalanceResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetPurchasedLikeBalance", varargs...)
	ret0, _ := ret[0].(*gen.GetPurchasedLikeBalanceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPurchasedLikeBalance indicates an expected call of GetPurchasedLikeBalance.
func (mr *MockPaymentClientMockRecorder) GetPurchasedLikeBalance(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPurchasedLikeBalance", reflect.TypeOf((*MockPaymentClient)(nil).GetPurchasedLikeBalance), varargs...)
}

// RefreshDailyLikeBalance mocks base method.
func (m *MockPaymentClient) RefreshDailyLikeBalance(ctx context.Context, in *gen.RefreshDailyLikeBalanceRequest, opts ...grpc.CallOption) (*gen.RefreshDailyLikeBalanceResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RefreshDailyLikeBalance", varargs...)
	ret0, _ := ret[0].(*gen.RefreshDailyLikeBalanceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshDailyLikeBalance indicates an expected call of RefreshDailyLikeBalance.
func (mr *MockPaymentClientMockRecorder) RefreshDailyLikeBalance(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshDailyLikeBalance", reflect.TypeOf((*MockPaymentClient)(nil).RefreshDailyLikeBalance), varargs...)
}

// MockPaymentServer is a mock of PaymentServer interface.
type MockPaymentServer struct {
	ctrl     *gomock.Controller
	recorder *MockPaymentServerMockRecorder
}

// MockPaymentServerMockRecorder is the mock recorder for MockPaymentServer.
type MockPaymentServerMockRecorder struct {
	mock *MockPaymentServer
}

// NewMockPaymentServer creates a new mock instance.
func NewMockPaymentServer(ctrl *gomock.Controller) *MockPaymentServer {
	mock := &MockPaymentServer{ctrl: ctrl}
	mock.recorder = &MockPaymentServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPaymentServer) EXPECT() *MockPaymentServerMockRecorder {
	return m.recorder
}

// BuyLikes mocks base method.
func (m *MockPaymentServer) BuyLikes(arg0 context.Context, arg1 *gen.BuyLikesRequest) (*gen.BuyLikesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuyLikes", arg0, arg1)
	ret0, _ := ret[0].(*gen.BuyLikesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BuyLikes indicates an expected call of BuyLikes.
func (mr *MockPaymentServerMockRecorder) BuyLikes(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuyLikes", reflect.TypeOf((*MockPaymentServer)(nil).BuyLikes), arg0, arg1)
}

// ChangeBalance mocks base method.
func (m *MockPaymentServer) ChangeBalance(arg0 context.Context, arg1 *gen.ChangeBalanceRequest) (*gen.ChangeBalanceResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeBalance", arg0, arg1)
	ret0, _ := ret[0].(*gen.ChangeBalanceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChangeBalance indicates an expected call of ChangeBalance.
func (mr *MockPaymentServerMockRecorder) ChangeBalance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeBalance", reflect.TypeOf((*MockPaymentServer)(nil).ChangeBalance), arg0, arg1)
}

// ChangePurchasedLikesBalance mocks base method.
func (m *MockPaymentServer) ChangePurchasedLikesBalance(arg0 context.Context, arg1 *gen.ChangePurchasedLikesBalanceRequest) (*gen.ChangePurchasedLikesBalanceResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangePurchasedLikesBalance", arg0, arg1)
	ret0, _ := ret[0].(*gen.ChangePurchasedLikesBalanceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChangePurchasedLikesBalance indicates an expected call of ChangePurchasedLikesBalance.
func (mr *MockPaymentServerMockRecorder) ChangePurchasedLikesBalance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangePurchasedLikesBalance", reflect.TypeOf((*MockPaymentServer)(nil).ChangePurchasedLikesBalance), arg0, arg1)
}

// CheckAndSpendLike mocks base method.
func (m *MockPaymentServer) CheckAndSpendLike(arg0 context.Context, arg1 *gen.CheckAndSpendLikeRequest) (*gen.CheckAndSpendLikeResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAndSpendLike", arg0, arg1)
	ret0, _ := ret[0].(*gen.CheckAndSpendLikeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckAndSpendLike indicates an expected call of CheckAndSpendLike.
func (mr *MockPaymentServerMockRecorder) CheckAndSpendLike(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAndSpendLike", reflect.TypeOf((*MockPaymentServer)(nil).CheckAndSpendLike), arg0, arg1)
}

// CreateBalances mocks base method.
func (m *MockPaymentServer) CreateBalances(arg0 context.Context, arg1 *gen.CreateBalancesRequest) (*gen.CreateBalancesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBalances", arg0, arg1)
	ret0, _ := ret[0].(*gen.CreateBalancesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBalances indicates an expected call of CreateBalances.
func (mr *MockPaymentServerMockRecorder) CreateBalances(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBalances", reflect.TypeOf((*MockPaymentServer)(nil).CreateBalances), arg0, arg1)
}

// CreateProduct mocks base method.
func (m *MockPaymentServer) CreateProduct(arg0 context.Context, arg1 *gen.CreateProductRequest) (*gen.CreateProductResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProduct", arg0, arg1)
	ret0, _ := ret[0].(*gen.CreateProductResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProduct indicates an expected call of CreateProduct.
func (mr *MockPaymentServerMockRecorder) CreateProduct(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProduct", reflect.TypeOf((*MockPaymentServer)(nil).CreateProduct), arg0, arg1)
}

// GetAllBalance mocks base method.
func (m *MockPaymentServer) GetAllBalance(arg0 context.Context, arg1 *gen.GetAllBalanceRequest) (*gen.GetAllBalanceResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllBalance", arg0, arg1)
	ret0, _ := ret[0].(*gen.GetAllBalanceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllBalance indicates an expected call of GetAllBalance.
func (mr *MockPaymentServerMockRecorder) GetAllBalance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllBalance", reflect.TypeOf((*MockPaymentServer)(nil).GetAllBalance), arg0, arg1)
}

// GetBalance mocks base method.
func (m *MockPaymentServer) GetBalance(arg0 context.Context, arg1 *gen.GetBalanceRequest) (*gen.GetBalanceResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalance", arg0, arg1)
	ret0, _ := ret[0].(*gen.GetBalanceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBalance indicates an expected call of GetBalance.
func (mr *MockPaymentServerMockRecorder) GetBalance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalance", reflect.TypeOf((*MockPaymentServer)(nil).GetBalance), arg0, arg1)
}

// GetDailyLikeBalance mocks base method.
func (m *MockPaymentServer) GetDailyLikeBalance(arg0 context.Context, arg1 *gen.GetDailyLikeBalanceRequest) (*gen.GetDailyLikeBalanceResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDailyLikeBalance", arg0, arg1)
	ret0, _ := ret[0].(*gen.GetDailyLikeBalanceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDailyLikeBalance indicates an expected call of GetDailyLikeBalance.
func (mr *MockPaymentServerMockRecorder) GetDailyLikeBalance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDailyLikeBalance", reflect.TypeOf((*MockPaymentServer)(nil).GetDailyLikeBalance), arg0, arg1)
}

// GetProducts mocks base method.
func (m *MockPaymentServer) GetProducts(arg0 context.Context, arg1 *gen.GetProductsRequest) (*gen.GetProductsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProducts", arg0, arg1)
	ret0, _ := ret[0].(*gen.GetProductsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProducts indicates an expected call of GetProducts.
func (mr *MockPaymentServerMockRecorder) GetProducts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProducts", reflect.TypeOf((*MockPaymentServer)(nil).GetProducts), arg0, arg1)
}

// GetPurchasedLikeBalance mocks base method.
func (m *MockPaymentServer) GetPurchasedLikeBalance(arg0 context.Context, arg1 *gen.GetPurchasedLikeBalanceRequest) (*gen.GetPurchasedLikeBalanceResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPurchasedLikeBalance", arg0, arg1)
	ret0, _ := ret[0].(*gen.GetPurchasedLikeBalanceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPurchasedLikeBalance indicates an expected call of GetPurchasedLikeBalance.
func (mr *MockPaymentServerMockRecorder) GetPurchasedLikeBalance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPurchasedLikeBalance", reflect.TypeOf((*MockPaymentServer)(nil).GetPurchasedLikeBalance), arg0, arg1)
}

// RefreshDailyLikeBalance mocks base method.
func (m *MockPaymentServer) RefreshDailyLikeBalance(arg0 context.Context, arg1 *gen.RefreshDailyLikeBalanceRequest) (*gen.RefreshDailyLikeBalanceResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshDailyLikeBalance", arg0, arg1)
	ret0, _ := ret[0].(*gen.RefreshDailyLikeBalanceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshDailyLikeBalance indicates an expected call of RefreshDailyLikeBalance.
func (mr *MockPaymentServerMockRecorder) RefreshDailyLikeBalance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshDailyLikeBalance", reflect.TypeOf((*MockPaymentServer)(nil).RefreshDailyLikeBalance), arg0, arg1)
}

// mustEmbedUnimplementedPaymentServer mocks base method.
func (m *MockPaymentServer) mustEmbedUnimplementedPaymentServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedPaymentServer")
}

// mustEmbedUnimplementedPaymentServer indicates an expected call of mustEmbedUnimplementedPaymentServer.
func (mr *MockPaymentServerMockRecorder) mustEmbedUnimplementedPaymentServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedPaymentServer", reflect.TypeOf((*MockPaymentServer)(nil).mustEmbedUnimplementedPaymentServer))
}

// MockUnsafePaymentServer is a mock of UnsafePaymentServer interface.
type MockUnsafePaymentServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafePaymentServerMockRecorder
}

// MockUnsafePaymentServerMockRecorder is the mock recorder for MockUnsafePaymentServer.
type MockUnsafePaymentServerMockRecorder struct {
	mock *MockUnsafePaymentServer
}

// NewMockUnsafePaymentServer creates a new mock instance.
func NewMockUnsafePaymentServer(ctrl *gomock.Controller) *MockUnsafePaymentServer {
	mock := &MockUnsafePaymentServer{ctrl: ctrl}
	mock.recorder = &MockUnsafePaymentServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafePaymentServer) EXPECT() *MockUnsafePaymentServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedPaymentServer mocks base method.
func (m *MockUnsafePaymentServer) mustEmbedUnimplementedPaymentServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedPaymentServer")
}

// mustEmbedUnimplementedPaymentServer indicates an expected call of mustEmbedUnimplementedPaymentServer.
func (mr *MockUnsafePaymentServerMockRecorder) mustEmbedUnimplementedPaymentServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedPaymentServer", reflect.TypeOf((*MockUnsafePaymentServer)(nil).mustEmbedUnimplementedPaymentServer))
}

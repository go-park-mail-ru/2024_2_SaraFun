// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.3
// source: payments.proto

package gen

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Payment_GetDailyLikeBalance_FullMethodName         = "/payments.Payment/GetDailyLikeBalance"
	Payment_GetPurchasedLikeBalance_FullMethodName     = "/payments.Payment/GetPurchasedLikeBalance"
	Payment_GetBalance_FullMethodName                  = "/payments.Payment/GetBalance"
	Payment_RefreshDailyLikeBalance_FullMethodName     = "/payments.Payment/RefreshDailyLikeBalance"
	Payment_ChangeBalance_FullMethodName               = "/payments.Payment/ChangeBalance"
	Payment_CheckAndSpendLike_FullMethodName           = "/payments.Payment/CheckAndSpendLike"
	Payment_ChangePurchasedLikesBalance_FullMethodName = "/payments.Payment/ChangePurchasedLikesBalance"
	Payment_GetAllBalance_FullMethodName               = "/payments.Payment/GetAllBalance"
	Payment_CreateBalances_FullMethodName              = "/payments.Payment/CreateBalances"
	Payment_BuyLikes_FullMethodName                    = "/payments.Payment/BuyLikes"
	Payment_CreateProduct_FullMethodName               = "/payments.Payment/CreateProduct"
	Payment_GetProducts_FullMethodName                 = "/payments.Payment/GetProducts"
	Payment_AddAward_FullMethodName                    = "/payments.Payment/AddAward"
	Payment_GetAwards_FullMethodName                   = "/payments.Payment/GetAwards"
	Payment_UpdateActivity_FullMethodName              = "/payments.Payment/UpdateActivity"
	Payment_CreateActivity_FullMethodName              = "/payments.Payment/CreateActivity"
)

// PaymentClient is the client API for Payment service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PaymentClient interface {
	GetDailyLikeBalance(ctx context.Context, in *GetDailyLikeBalanceRequest, opts ...grpc.CallOption) (*GetDailyLikeBalanceResponse, error)
	GetPurchasedLikeBalance(ctx context.Context, in *GetPurchasedLikeBalanceRequest, opts ...grpc.CallOption) (*GetPurchasedLikeBalanceResponse, error)
	GetBalance(ctx context.Context, in *GetBalanceRequest, opts ...grpc.CallOption) (*GetBalanceResponse, error)
	RefreshDailyLikeBalance(ctx context.Context, in *RefreshDailyLikeBalanceRequest, opts ...grpc.CallOption) (*RefreshDailyLikeBalanceResponse, error)
	ChangeBalance(ctx context.Context, in *ChangeBalanceRequest, opts ...grpc.CallOption) (*ChangeBalanceResponse, error)
	CheckAndSpendLike(ctx context.Context, in *CheckAndSpendLikeRequest, opts ...grpc.CallOption) (*CheckAndSpendLikeResponse, error)
	ChangePurchasedLikesBalance(ctx context.Context, in *ChangePurchasedLikesBalanceRequest, opts ...grpc.CallOption) (*ChangePurchasedLikesBalanceResponse, error)
	GetAllBalance(ctx context.Context, in *GetAllBalanceRequest, opts ...grpc.CallOption) (*GetAllBalanceResponse, error)
	CreateBalances(ctx context.Context, in *CreateBalancesRequest, opts ...grpc.CallOption) (*CreateBalancesResponse, error)
	BuyLikes(ctx context.Context, in *BuyLikesRequest, opts ...grpc.CallOption) (*BuyLikesResponse, error)
	CreateProduct(ctx context.Context, in *CreateProductRequest, opts ...grpc.CallOption) (*CreateProductResponse, error)
	GetProducts(ctx context.Context, in *GetProductsRequest, opts ...grpc.CallOption) (*GetProductsResponse, error)
	AddAward(ctx context.Context, in *AddAwardRequest, opts ...grpc.CallOption) (*AddAwardResponse, error)
	GetAwards(ctx context.Context, in *GetAwardsRequest, opts ...grpc.CallOption) (*GetAwardsResponse, error)
	UpdateActivity(ctx context.Context, in *UpdateActivityRequest, opts ...grpc.CallOption) (*UpdateActivityResponse, error)
	CreateActivity(ctx context.Context, in *CreateActivityRequest, opts ...grpc.CallOption) (*CreateActivityResponse, error)
}

type paymentClient struct {
	cc grpc.ClientConnInterface
}

func NewPaymentClient(cc grpc.ClientConnInterface) PaymentClient {
	return &paymentClient{cc}
}

func (c *paymentClient) GetDailyLikeBalance(ctx context.Context, in *GetDailyLikeBalanceRequest, opts ...grpc.CallOption) (*GetDailyLikeBalanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetDailyLikeBalanceResponse)
	err := c.cc.Invoke(ctx, Payment_GetDailyLikeBalance_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentClient) GetPurchasedLikeBalance(ctx context.Context, in *GetPurchasedLikeBalanceRequest, opts ...grpc.CallOption) (*GetPurchasedLikeBalanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPurchasedLikeBalanceResponse)
	err := c.cc.Invoke(ctx, Payment_GetPurchasedLikeBalance_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentClient) GetBalance(ctx context.Context, in *GetBalanceRequest, opts ...grpc.CallOption) (*GetBalanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetBalanceResponse)
	err := c.cc.Invoke(ctx, Payment_GetBalance_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentClient) RefreshDailyLikeBalance(ctx context.Context, in *RefreshDailyLikeBalanceRequest, opts ...grpc.CallOption) (*RefreshDailyLikeBalanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RefreshDailyLikeBalanceResponse)
	err := c.cc.Invoke(ctx, Payment_RefreshDailyLikeBalance_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentClient) ChangeBalance(ctx context.Context, in *ChangeBalanceRequest, opts ...grpc.CallOption) (*ChangeBalanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ChangeBalanceResponse)
	err := c.cc.Invoke(ctx, Payment_ChangeBalance_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentClient) CheckAndSpendLike(ctx context.Context, in *CheckAndSpendLikeRequest, opts ...grpc.CallOption) (*CheckAndSpendLikeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CheckAndSpendLikeResponse)
	err := c.cc.Invoke(ctx, Payment_CheckAndSpendLike_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentClient) ChangePurchasedLikesBalance(ctx context.Context, in *ChangePurchasedLikesBalanceRequest, opts ...grpc.CallOption) (*ChangePurchasedLikesBalanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ChangePurchasedLikesBalanceResponse)
	err := c.cc.Invoke(ctx, Payment_ChangePurchasedLikesBalance_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentClient) GetAllBalance(ctx context.Context, in *GetAllBalanceRequest, opts ...grpc.CallOption) (*GetAllBalanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAllBalanceResponse)
	err := c.cc.Invoke(ctx, Payment_GetAllBalance_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentClient) CreateBalances(ctx context.Context, in *CreateBalancesRequest, opts ...grpc.CallOption) (*CreateBalancesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateBalancesResponse)
	err := c.cc.Invoke(ctx, Payment_CreateBalances_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentClient) BuyLikes(ctx context.Context, in *BuyLikesRequest, opts ...grpc.CallOption) (*BuyLikesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BuyLikesResponse)
	err := c.cc.Invoke(ctx, Payment_BuyLikes_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentClient) CreateProduct(ctx context.Context, in *CreateProductRequest, opts ...grpc.CallOption) (*CreateProductResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateProductResponse)
	err := c.cc.Invoke(ctx, Payment_CreateProduct_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentClient) GetProducts(ctx context.Context, in *GetProductsRequest, opts ...grpc.CallOption) (*GetProductsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetProductsResponse)
	err := c.cc.Invoke(ctx, Payment_GetProducts_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentClient) AddAward(ctx context.Context, in *AddAwardRequest, opts ...grpc.CallOption) (*AddAwardResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddAwardResponse)
	err := c.cc.Invoke(ctx, Payment_AddAward_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentClient) GetAwards(ctx context.Context, in *GetAwardsRequest, opts ...grpc.CallOption) (*GetAwardsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAwardsResponse)
	err := c.cc.Invoke(ctx, Payment_GetAwards_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentClient) UpdateActivity(ctx context.Context, in *UpdateActivityRequest, opts ...grpc.CallOption) (*UpdateActivityResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateActivityResponse)
	err := c.cc.Invoke(ctx, Payment_UpdateActivity_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentClient) CreateActivity(ctx context.Context, in *CreateActivityRequest, opts ...grpc.CallOption) (*CreateActivityResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateActivityResponse)
	err := c.cc.Invoke(ctx, Payment_CreateActivity_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PaymentServer is the server API for Payment service.
// All implementations must embed UnimplementedPaymentServer
// for forward compatibility.
type PaymentServer interface {
	GetDailyLikeBalance(context.Context, *GetDailyLikeBalanceRequest) (*GetDailyLikeBalanceResponse, error)
	GetPurchasedLikeBalance(context.Context, *GetPurchasedLikeBalanceRequest) (*GetPurchasedLikeBalanceResponse, error)
	GetBalance(context.Context, *GetBalanceRequest) (*GetBalanceResponse, error)
	RefreshDailyLikeBalance(context.Context, *RefreshDailyLikeBalanceRequest) (*RefreshDailyLikeBalanceResponse, error)
	ChangeBalance(context.Context, *ChangeBalanceRequest) (*ChangeBalanceResponse, error)
	CheckAndSpendLike(context.Context, *CheckAndSpendLikeRequest) (*CheckAndSpendLikeResponse, error)
	ChangePurchasedLikesBalance(context.Context, *ChangePurchasedLikesBalanceRequest) (*ChangePurchasedLikesBalanceResponse, error)
	GetAllBalance(context.Context, *GetAllBalanceRequest) (*GetAllBalanceResponse, error)
	CreateBalances(context.Context, *CreateBalancesRequest) (*CreateBalancesResponse, error)
	BuyLikes(context.Context, *BuyLikesRequest) (*BuyLikesResponse, error)
	CreateProduct(context.Context, *CreateProductRequest) (*CreateProductResponse, error)
	GetProducts(context.Context, *GetProductsRequest) (*GetProductsResponse, error)
	AddAward(context.Context, *AddAwardRequest) (*AddAwardResponse, error)
	GetAwards(context.Context, *GetAwardsRequest) (*GetAwardsResponse, error)
	UpdateActivity(context.Context, *UpdateActivityRequest) (*UpdateActivityResponse, error)
	CreateActivity(context.Context, *CreateActivityRequest) (*CreateActivityResponse, error)
	mustEmbedUnimplementedPaymentServer()
}

// UnimplementedPaymentServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPaymentServer struct{}

func (UnimplementedPaymentServer) GetDailyLikeBalance(context.Context, *GetDailyLikeBalanceRequest) (*GetDailyLikeBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDailyLikeBalance not implemented")
}
func (UnimplementedPaymentServer) GetPurchasedLikeBalance(context.Context, *GetPurchasedLikeBalanceRequest) (*GetPurchasedLikeBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPurchasedLikeBalance not implemented")
}
func (UnimplementedPaymentServer) GetBalance(context.Context, *GetBalanceRequest) (*GetBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBalance not implemented")
}
func (UnimplementedPaymentServer) RefreshDailyLikeBalance(context.Context, *RefreshDailyLikeBalanceRequest) (*RefreshDailyLikeBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefreshDailyLikeBalance not implemented")
}
func (UnimplementedPaymentServer) ChangeBalance(context.Context, *ChangeBalanceRequest) (*ChangeBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeBalance not implemented")
}
func (UnimplementedPaymentServer) CheckAndSpendLike(context.Context, *CheckAndSpendLikeRequest) (*CheckAndSpendLikeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckAndSpendLike not implemented")
}
func (UnimplementedPaymentServer) ChangePurchasedLikesBalance(context.Context, *ChangePurchasedLikesBalanceRequest) (*ChangePurchasedLikesBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangePurchasedLikesBalance not implemented")
}
func (UnimplementedPaymentServer) GetAllBalance(context.Context, *GetAllBalanceRequest) (*GetAllBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllBalance not implemented")
}
func (UnimplementedPaymentServer) CreateBalances(context.Context, *CreateBalancesRequest) (*CreateBalancesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBalances not implemented")
}
func (UnimplementedPaymentServer) BuyLikes(context.Context, *BuyLikesRequest) (*BuyLikesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BuyLikes not implemented")
}
func (UnimplementedPaymentServer) CreateProduct(context.Context, *CreateProductRequest) (*CreateProductResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateProduct not implemented")
}
func (UnimplementedPaymentServer) GetProducts(context.Context, *GetProductsRequest) (*GetProductsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProducts not implemented")
}
func (UnimplementedPaymentServer) AddAward(context.Context, *AddAwardRequest) (*AddAwardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddAward not implemented")
}
func (UnimplementedPaymentServer) GetAwards(context.Context, *GetAwardsRequest) (*GetAwardsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAwards not implemented")
}
func (UnimplementedPaymentServer) UpdateActivity(context.Context, *UpdateActivityRequest) (*UpdateActivityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateActivity not implemented")
}
func (UnimplementedPaymentServer) CreateActivity(context.Context, *CreateActivityRequest) (*CreateActivityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateActivity not implemented")
}
func (UnimplementedPaymentServer) mustEmbedUnimplementedPaymentServer() {}
func (UnimplementedPaymentServer) testEmbeddedByValue()                 {}

// UnsafePaymentServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PaymentServer will
// result in compilation errors.
type UnsafePaymentServer interface {
	mustEmbedUnimplementedPaymentServer()
}

func RegisterPaymentServer(s grpc.ServiceRegistrar, srv PaymentServer) {
	// If the following call pancis, it indicates UnimplementedPaymentServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Payment_ServiceDesc, srv)
}

func _Payment_GetDailyLikeBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDailyLikeBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServer).GetDailyLikeBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Payment_GetDailyLikeBalance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServer).GetDailyLikeBalance(ctx, req.(*GetDailyLikeBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payment_GetPurchasedLikeBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPurchasedLikeBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServer).GetPurchasedLikeBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Payment_GetPurchasedLikeBalance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServer).GetPurchasedLikeBalance(ctx, req.(*GetPurchasedLikeBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payment_GetBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServer).GetBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Payment_GetBalance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServer).GetBalance(ctx, req.(*GetBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payment_RefreshDailyLikeBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshDailyLikeBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServer).RefreshDailyLikeBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Payment_RefreshDailyLikeBalance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServer).RefreshDailyLikeBalance(ctx, req.(*RefreshDailyLikeBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payment_ChangeBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangeBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServer).ChangeBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Payment_ChangeBalance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServer).ChangeBalance(ctx, req.(*ChangeBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payment_CheckAndSpendLike_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckAndSpendLikeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServer).CheckAndSpendLike(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Payment_CheckAndSpendLike_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServer).CheckAndSpendLike(ctx, req.(*CheckAndSpendLikeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payment_ChangePurchasedLikesBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangePurchasedLikesBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServer).ChangePurchasedLikesBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Payment_ChangePurchasedLikesBalance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServer).ChangePurchasedLikesBalance(ctx, req.(*ChangePurchasedLikesBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payment_GetAllBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServer).GetAllBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Payment_GetAllBalance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServer).GetAllBalance(ctx, req.(*GetAllBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payment_CreateBalances_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateBalancesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServer).CreateBalances(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Payment_CreateBalances_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServer).CreateBalances(ctx, req.(*CreateBalancesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payment_BuyLikes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BuyLikesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServer).BuyLikes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Payment_BuyLikes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServer).BuyLikes(ctx, req.(*BuyLikesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payment_CreateProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServer).CreateProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Payment_CreateProduct_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServer).CreateProduct(ctx, req.(*CreateProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payment_GetProducts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProductsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServer).GetProducts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Payment_GetProducts_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServer).GetProducts(ctx, req.(*GetProductsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payment_AddAward_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddAwardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServer).AddAward(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Payment_AddAward_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServer).AddAward(ctx, req.(*AddAwardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payment_GetAwards_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAwardsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServer).GetAwards(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Payment_GetAwards_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServer).GetAwards(ctx, req.(*GetAwardsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payment_UpdateActivity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateActivityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServer).UpdateActivity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Payment_UpdateActivity_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServer).UpdateActivity(ctx, req.(*UpdateActivityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payment_CreateActivity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateActivityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServer).CreateActivity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Payment_CreateActivity_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServer).CreateActivity(ctx, req.(*CreateActivityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Payment_ServiceDesc is the grpc.ServiceDesc for Payment service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Payment_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "payments.Payment",
	HandlerType: (*PaymentServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDailyLikeBalance",
			Handler:    _Payment_GetDailyLikeBalance_Handler,
		},
		{
			MethodName: "GetPurchasedLikeBalance",
			Handler:    _Payment_GetPurchasedLikeBalance_Handler,
		},
		{
			MethodName: "GetBalance",
			Handler:    _Payment_GetBalance_Handler,
		},
		{
			MethodName: "RefreshDailyLikeBalance",
			Handler:    _Payment_RefreshDailyLikeBalance_Handler,
		},
		{
			MethodName: "ChangeBalance",
			Handler:    _Payment_ChangeBalance_Handler,
		},
		{
			MethodName: "CheckAndSpendLike",
			Handler:    _Payment_CheckAndSpendLike_Handler,
		},
		{
			MethodName: "ChangePurchasedLikesBalance",
			Handler:    _Payment_ChangePurchasedLikesBalance_Handler,
		},
		{
			MethodName: "GetAllBalance",
			Handler:    _Payment_GetAllBalance_Handler,
		},
		{
			MethodName: "CreateBalances",
			Handler:    _Payment_CreateBalances_Handler,
		},
		{
			MethodName: "BuyLikes",
			Handler:    _Payment_BuyLikes_Handler,
		},
		{
			MethodName: "CreateProduct",
			Handler:    _Payment_CreateProduct_Handler,
		},
		{
			MethodName: "GetProducts",
			Handler:    _Payment_GetProducts_Handler,
		},
		{
			MethodName: "AddAward",
			Handler:    _Payment_AddAward_Handler,
		},
		{
			MethodName: "GetAwards",
			Handler:    _Payment_GetAwards_Handler,
		},
		{
			MethodName: "UpdateActivity",
			Handler:    _Payment_UpdateActivity_Handler,
		},
		{
			MethodName: "CreateActivity",
			Handler:    _Payment_CreateActivity_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "payments.proto",
}

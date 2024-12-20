package grpc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc"
	generatedPayments "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc/mocks"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestGRPCHandler_GetBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uc := mocks.NewMockUseCase(ctrl)
	logger := zap.NewNop()
	h := grpc.NewGrpcPaymentsHandler(uc, logger)

	ctx := context.WithValue(context.Background(), consts.RequestIDKey, "test_req_id")
	req := &generatedPayments.GetBalanceRequest{UserID: 10}

	t.Run("success", func(t *testing.T) {
		uc.EXPECT().GetBalance(ctx, 10).Return(100, nil)
		resp, err := h.GetBalance(ctx, req)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.Balance != 100 {
			t.Errorf("got %v, want 100", resp.Balance)
		}
	})

	t.Run("error", func(t *testing.T) {
		uc.EXPECT().GetBalance(ctx, 10).Return(0, errors.New("balance error"))
		_, err := h.GetBalance(ctx, req)
		if err == nil || !contains(err.Error(), "grpc get balance error") {
			t.Errorf("expected error got %v", err)
		}
	})
}

func TestGRPCHandler_GetDailyLikeBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uc := mocks.NewMockUseCase(ctrl)
	logger := zap.NewNop()
	h := grpc.NewGrpcPaymentsHandler(uc, logger)

	ctx := context.WithValue(context.Background(), consts.RequestIDKey, "test_req_id")
	req := &generatedPayments.GetDailyLikeBalanceRequest{UserID: 20}

	t.Run("success", func(t *testing.T) {
		uc.EXPECT().GetDailyLikesCount(ctx, 20).Return(5, nil)
		resp, err := h.GetDailyLikeBalance(ctx, req)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.Balance != 5 {
			t.Errorf("got %v, want 5", resp.Balance)
		}
	})

	t.Run("error", func(t *testing.T) {
		uc.EXPECT().GetDailyLikesCount(ctx, 20).Return(0, errors.New("likes error"))
		_, err := h.GetDailyLikeBalance(ctx, req)
		if err == nil || !contains(err.Error(), "grpc get balance error") {
			t.Errorf("expected error got %v", err)
		}
	})
}

func TestGRPCHandler_GetPurchasedLikeBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uc := mocks.NewMockUseCase(ctrl)
	logger := zap.NewNop()
	h := grpc.NewGrpcPaymentsHandler(uc, logger)

	ctx := context.WithValue(context.Background(), consts.RequestIDKey, "test_req_id")
	req := &generatedPayments.GetPurchasedLikeBalanceRequest{UserID: 30}

	t.Run("success", func(t *testing.T) {
		uc.EXPECT().GetPurchasedLikesCount(ctx, 30).Return(10, nil)
		resp, err := h.GetPurchasedLikeBalance(ctx, req)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.Balance != 10 {
			t.Errorf("got %v, want 10", resp.Balance)
		}
	})

	t.Run("error", func(t *testing.T) {
		uc.EXPECT().GetPurchasedLikesCount(ctx, 30).Return(0, errors.New("purchased error"))
		_, err := h.GetPurchasedLikeBalance(ctx, req)
		if err == nil || !contains(err.Error(), "grpc get balance error") {
			t.Errorf("expected error got %v", err)
		}
	})
}

func TestGRPCHandler_RefreshDailyLikeBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uc := mocks.NewMockUseCase(ctrl)
	logger := zap.NewNop()
	h := grpc.NewGrpcPaymentsHandler(uc, logger)

	ctx := context.Background()
	req := &generatedPayments.RefreshDailyLikeBalanceRequest{}

	t.Run("success", func(t *testing.T) {
		uc.EXPECT().SetDailyLikeCountToAll(ctx, consts.DailyLikeLimit).Return(nil)
		_, err := h.RefreshDailyLikeBalance(ctx, req)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		uc.EXPECT().SetDailyLikeCountToAll(ctx, consts.DailyLikeLimit).Return(errors.New("set error"))
		_, err := h.RefreshDailyLikeBalance(ctx, req)
		if err == nil || !contains(err.Error(), "grpc set daily like like count error") {
			t.Errorf("expected error got %v", err)
		}
	})
}

func TestGRPCHandler_ChangeBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uc := mocks.NewMockUseCase(ctrl)
	logger := zap.NewNop()
	h := grpc.NewGrpcPaymentsHandler(uc, logger)

	ctx := context.Background()
	req := &generatedPayments.ChangeBalanceRequest{UserID: 40, Amount: 100}

	t.Run("success", func(t *testing.T) {
		uc.EXPECT().ChangeBalance(ctx, 40, 100).Return(nil)
		_, err := h.ChangeBalance(ctx, req)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		uc.EXPECT().ChangeBalance(ctx, 40, 100).Return(errors.New("change error"))
		_, err := h.ChangeBalance(ctx, req)
		if err == nil || !contains(err.Error(), "grpc change balance error") {
			t.Errorf("expected error got %v", err)
		}
	})
}

func TestGRPCHandler_CheckAndSpendLike(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uc := mocks.NewMockUseCase(ctrl)
	logger := zap.NewNop()
	h := grpc.NewGrpcPaymentsHandler(uc, logger)

	ctx := context.Background()
	req := &generatedPayments.CheckAndSpendLikeRequest{UserID: 50}

	t.Run("have daily like", func(t *testing.T) {
		uc.EXPECT().GetDailyLikesCount(ctx, 50).Return(1, nil)
		uc.EXPECT().ChangeDailyLikeCount(ctx, 50, -1).Return(nil)
		_, err := h.CheckAndSpendLike(ctx, req)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("no daily like, have purchased", func(t *testing.T) {
		uc.EXPECT().GetDailyLikesCount(ctx, 50).Return(0, nil)
		uc.EXPECT().GetPurchasedLikesCount(ctx, 50).Return(2, nil)
		uc.EXPECT().ChangePurchasedLikeCount(ctx, 50, -1).Return(nil)
		_, err := h.CheckAndSpendLike(ctx, req)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("no daily like, no purchased", func(t *testing.T) {
		uc.EXPECT().GetDailyLikesCount(ctx, 50).Return(0, nil)
		uc.EXPECT().GetPurchasedLikesCount(ctx, 50).Return(0, nil)
		_, err := h.CheckAndSpendLike(ctx, req)
		if err == nil || !contains(err.Error(), "dont have likes") {
			t.Errorf("expected error got %v", err)
		}
	})

	t.Run("error get daily likes", func(t *testing.T) {
		uc.EXPECT().GetDailyLikesCount(ctx, 50).Return(0, errors.New("daily error"))
		_, err := h.CheckAndSpendLike(ctx, req)
		if err == nil || !contains(err.Error(), "grpc get daily likes count error") {
			t.Errorf("expected error got %v", err)
		}
	})

	t.Run("error get purchased likes", func(t *testing.T) {
		uc.EXPECT().GetDailyLikesCount(ctx, 50).Return(0, nil)
		uc.EXPECT().GetPurchasedLikesCount(ctx, 50).Return(0, errors.New("purchased error"))
		_, err := h.CheckAndSpendLike(ctx, req)
		if err == nil || !contains(err.Error(), "grpc get balance error") {
			t.Errorf("expected error got %v", err)
		}
	})

	t.Run("error change daily", func(t *testing.T) {
		uc.EXPECT().GetDailyLikesCount(ctx, 50).Return(1, nil)
		uc.EXPECT().ChangeDailyLikeCount(ctx, 50, -1).Return(errors.New("change error"))
		_, err := h.CheckAndSpendLike(ctx, req)
		if err == nil || !contains(err.Error(), "grpc change balance error") {
			t.Errorf("expected error got %v", err)
		}
	})

}

func TestGRPCHandler_ChangePurchasedLikesBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uc := mocks.NewMockUseCase(ctrl)
	logger := zap.NewNop()
	h := grpc.NewGrpcPaymentsHandler(uc, logger)

	ctx := context.Background()
	req := &generatedPayments.ChangePurchasedLikesBalanceRequest{UserID: 60, Amount: 10}

	t.Run("success", func(t *testing.T) {
		uc.EXPECT().ChangePurchasedLikeCount(ctx, 60, 10).Return(nil)
		_, err := h.ChangePurchasedLikesBalance(ctx, req)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		uc.EXPECT().ChangePurchasedLikeCount(ctx, 60, 10).Return(errors.New("change error"))
		_, err := h.ChangePurchasedLikesBalance(ctx, req)
		if err == nil || !contains(err.Error(), "grpc change balance error") {
			t.Errorf("expected error got %v", err)
		}
	})
}

func TestGRPCHandler_GetAllBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uc := mocks.NewMockUseCase(ctrl)
	logger := zap.NewNop()
	h := grpc.NewGrpcPaymentsHandler(uc, logger)

	ctx := context.Background()
	req := &generatedPayments.GetAllBalanceRequest{UserID: 70}

	t.Run("success", func(t *testing.T) {
		uc.EXPECT().GetDailyLikesCount(ctx, 70).Return(1, nil)
		uc.EXPECT().GetPurchasedLikesCount(ctx, 70).Return(2, nil)
		uc.EXPECT().GetBalance(ctx, 70).Return(100, nil)

		resp, err := h.GetAllBalance(ctx, req)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.DailyLikeBalance != 1 || resp.PurchasedLikeBalance != 2 || resp.MoneyBalance != 100 {
			t.Errorf("balances mismatch: got %+v", resp)
		}
	})

	t.Run("error daily", func(t *testing.T) {
		uc.EXPECT().GetDailyLikesCount(ctx, 70).Return(0, errors.New("daily error"))
		_, err := h.GetAllBalance(ctx, req)
		if err == nil || !contains(err.Error(), "grpc get balance error") {
			t.Errorf("expected error got %v", err)
		}
	})

	t.Run("error purchased", func(t *testing.T) {
		uc.EXPECT().GetDailyLikesCount(ctx, 70).Return(1, nil)
		uc.EXPECT().GetPurchasedLikesCount(ctx, 70).Return(0, errors.New("purchased error"))
		_, err := h.GetAllBalance(ctx, req)
		if err == nil || !contains(err.Error(), "grpc get balance error") {
			t.Errorf("expected error got %v", err)
		}
	})

	t.Run("error money", func(t *testing.T) {
		uc.EXPECT().GetDailyLikesCount(ctx, 70).Return(1, nil)
		uc.EXPECT().GetPurchasedLikesCount(ctx, 70).Return(2, nil)
		uc.EXPECT().GetBalance(ctx, 70).Return(0, errors.New("balance error"))
		_, err := h.GetAllBalance(ctx, req)
		if err == nil || !contains(err.Error(), "grpc get balance error") {
			t.Errorf("expected error got %v", err)
		}
	})
}

func TestGRPCHandler_CreateBalances(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uc := mocks.NewMockUseCase(ctrl)
	logger := zap.NewNop()
	h := grpc.NewGrpcPaymentsHandler(uc, logger)

	ctx := context.Background()
	req := &generatedPayments.CreateBalancesRequest{
		UserID:          80,
		MoneyAmount:     100,
		DailyAmount:     10,
		PurchasedAmount: 5,
	}

	t.Run("success", func(t *testing.T) {
		uc.EXPECT().AddBalance(ctx, 80, 100).Return(nil)
		uc.EXPECT().AddDailyLikesCount(ctx, 80, 10).Return(nil)
		uc.EXPECT().AddPurchasedLikesCount(ctx, 80, 5).Return(nil)

		_, err := h.CreateBalances(ctx, req)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("money error", func(t *testing.T) {
		uc.EXPECT().AddBalance(ctx, 80, 100).Return(errors.New("money error"))
		_, err := h.CreateBalances(ctx, req)
		if err == nil || !contains(err.Error(), "bad add balance error") {
			t.Errorf("expected error got %v", err)
		}
	})

	t.Run("daily error", func(t *testing.T) {
		uc.EXPECT().AddBalance(ctx, 80, 100).Return(nil)
		uc.EXPECT().AddDailyLikesCount(ctx, 80, 10).Return(errors.New("daily error"))
		_, err := h.CreateBalances(ctx, req)
		if err == nil || !contains(err.Error(), "bad daily likes count error") {
			t.Errorf("expected error got %v", err)
		}
	})

	t.Run("purchased error", func(t *testing.T) {
		uc.EXPECT().AddBalance(ctx, 80, 100).Return(nil)
		uc.EXPECT().AddDailyLikesCount(ctx, 80, 10).Return(nil)
		uc.EXPECT().AddPurchasedLikesCount(ctx, 80, 5).Return(errors.New("purch error"))
		_, err := h.CreateBalances(ctx, req)
		if err == nil || !contains(err.Error(), "bad purchase count error") {
			t.Errorf("expected error got %v", err)
		}
	})
}

//func TestGRPCHandler_BuyLikes(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	uc := mocks.NewMockUseCase(ctrl)
//	logger := zap.NewNop()
//	h := grpc.NewGrpcPaymentsHandler(uc, logger)
//
//	ctx := context.Background()
//	req := &generatedPayments.BuyLikesRequest{
//		Title:  "likes",
//		Amount: 100,
//		UserID: 90,
//	}
//
//	product := models.Product{Price: 10}
//
//	t.Run("success", func(t *testing.T) {
//		uc.EXPECT().GetProduct(ctx, "likes").Return(product, nil)
//		uc.EXPECT().CheckBalance(ctx, 90, 100).Return(nil)
//		// amount=100, price=10 => count=100/10=10 likes
//		uc.EXPECT().ChangeBalance(ctx, 90, -100).Return(nil)
//		uc.EXPECT().ChangePurchasedLikeCount(ctx, 90, 10).Return(nil)
//		_, err := h.BuyLikes(ctx, req)
//		if err != nil {
//			t.Errorf("unexpected error: %v", err)
//		}
//	})
//
//	t.Run("get product error", func(t *testing.T) {
//		uc.EXPECT().GetProduct(ctx, "likes").Return(models.Product{}, errors.New("prod error"))
//		_, err := h.BuyLikes(ctx, req)
//		if err == nil || !contains(err.Error(), "grpc get balance error") {
//			t.Errorf("expected error got %v", err)
//		}
//	})
//
//	t.Run("check balance error", func(t *testing.T) {
//		uc.EXPECT().GetProduct(ctx, "likes").Return(product, nil)
//		uc.EXPECT().CheckBalance(ctx, 90, 100).Return(errors.New("no money"))
//		_, err := h.BuyLikes(ctx, req)
//		if st, ok := status.FromError(err); !ok || st.Code() != codes.InvalidArgument {
//			t.Errorf("expected InvalidArgument, got %v", err)
//		}
//	})
//
//	t.Run("change balance error after success check", func(t *testing.T) {
//		uc.EXPECT().GetProduct(ctx, "likes").Return(product, nil)
//		uc.EXPECT().CheckBalance(ctx, 90, 100).Return(nil)
//		uc.EXPECT().ChangeBalance(ctx, 90, -100).Return(errors.New("change err"))
//		_, err := h.BuyLikes(ctx, req)
//		if err == nil || !contains(err.Error(), "grpc change balance error") {
//			t.Errorf("expected error got %v", err)
//		}
//	})
//
//	t.Run("change purchased error", func(t *testing.T) {
//		uc.EXPECT().GetProduct(ctx, "likes").Return(product, nil)
//		uc.EXPECT().CheckBalance(ctx, 90, 100).Return(nil)
//		uc.EXPECT().ChangeBalance(ctx, 90, -100).Return(nil)
//		uc.EXPECT().ChangePurchasedLikeCount(ctx, 90, 10).Return(errors.New("purch err"))
//		_, err := h.BuyLikes(ctx, req)
//		if err == nil || !contains(err.Error(), "grpc change balance error") {
//			t.Errorf("expected error got %v", err)
//		}
//	})
//}

func TestGRPCHandler_CreateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uc := mocks.NewMockUseCase(ctrl)
	logger := zap.NewNop()
	h := grpc.NewGrpcPaymentsHandler(uc, logger)

	ctx := context.Background()
	req := &generatedPayments.CreateProductRequest{
		Product: &generatedPayments.Product{
			Title:       "prod",
			Description: "desc",
			ImageLink:   "img",
			Price:       50,
		},
	}

	t.Run("success", func(t *testing.T) {
		uc.EXPECT().CreateProduct(ctx, models.Product{
			Title:       "prod",
			Description: "desc",
			ImageLink:   "img",
			Price:       50,
		}).Return(999, nil)
		resp, err := h.CreateProduct(ctx, req)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.ID != 999 {
			t.Errorf("got %v, want 999", resp.ID)
		}
	})

	t.Run("error", func(t *testing.T) {
		uc.EXPECT().CreateProduct(ctx, gomock.Any()).Return(0, errors.New("create prod err"))
		_, err := h.CreateProduct(ctx, req)
		if err == nil || !contains(err.Error(), "grpc create product error") {
			t.Errorf("expected error got %v", err)
		}
	})
}

func TestGRPCHandler_GetProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uc := mocks.NewMockUseCase(ctrl)
	logger := zap.NewNop()
	h := grpc.NewGrpcPaymentsHandler(uc, logger)

	ctx := context.Background()
	req := &generatedPayments.GetProductsRequest{}

	t.Run("success", func(t *testing.T) {
		products := []models.Product{
			{Title: "p1", Description: "d1", ImageLink: "i1", Price: 10},
			{Title: "p2", Description: "d2", ImageLink: "i2", Price: 20},
		}
		uc.EXPECT().GetProducts(ctx).Return(products, nil)
		resp, err := h.GetProducts(ctx, req)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(resp.Products) != 2 {
			t.Errorf("got %d products, want 2", len(resp.Products))
		}
	})

	t.Run("error", func(t *testing.T) {
		uc.EXPECT().GetProducts(ctx).Return(nil, errors.New("get prod err"))
		_, err := h.GetProducts(ctx, req)
		if err == nil || !contains(err.Error(), "grpc get products error") {
			t.Errorf("expected error got %v", err)
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

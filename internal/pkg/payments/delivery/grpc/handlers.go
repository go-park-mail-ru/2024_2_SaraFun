package grpc

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedPayments "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate mockgen -destination=./mocks/mock_usecase.go -package=mocks github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc UseCase

type UseCase interface {
	GetBalance(ctx context.Context, userID int) (int, error)
	GetDailyLikesCount(ctx context.Context, userID int) (int, error)
	GetPurchasedLikesCount(ctx context.Context, userID int) (int, error)
	SetDailyLikeCountToAll(ctx context.Context, amount int) error
	ChangeBalance(ctx context.Context, userID int, amount int) error
	ChangeDailyLikeCount(ctx context.Context, userID int, amount int) error
	ChangePurchasedLikeCount(ctx context.Context, userID int, amount int) error
	AddBalance(ctx context.Context, userID int, amount int) error
	AddDailyLikesCount(ctx context.Context, userID int, amount int) error
	AddPurchasedLikesCount(ctx context.Context, userID int, amount int) error
	GetProduct(ctx context.Context, title string) (models.Product, error)
	CheckBalance(ctx context.Context, userID int, amount int) error
	CreateProduct(ctx context.Context, product models.Product) (int, error)
	GetProducts(ctx context.Context) ([]models.Product, error)
	AddAward(ctx context.Context, award models.Award) error
	GetAwards(ctx context.Context) ([]models.Award, error)
	UpdateActivity(ctx context.Context, userID int) (string, error)
	AddActivity(ctx context.Context, userID int) error
}

type GRPCHandler struct {
	generatedPayments.PaymentServer
	uc     UseCase
	logger *zap.Logger
}

func NewGrpcPaymentsHandler(uc UseCase, logger *zap.Logger) *GRPCHandler {
	return &GRPCHandler{
		uc:     uc,
		logger: logger,
	}
}

func (h *GRPCHandler) GetBalance(ctx context.Context,
	in *generatedPayments.GetBalanceRequest) (*generatedPayments.GetBalanceResponse, error) {
	userID := int(in.UserID)
	balance, err := h.uc.GetBalance(ctx, userID)
	if err != nil {
		h.logger.Error("grpc get balance error", zap.Error(err))
		return nil, fmt.Errorf("grpc get balance error: %w", err)
	}
	response := &generatedPayments.GetBalanceResponse{
		Balance: int32(balance),
	}
	return response, nil
}

func (h *GRPCHandler) GetDailyLikeBalance(ctx context.Context,
	in *generatedPayments.GetDailyLikeBalanceRequest) (*generatedPayments.GetDailyLikeBalanceResponse, error) {
	userID := int(in.UserID)
	balance, err := h.uc.GetDailyLikesCount(ctx, userID)
	if err != nil {
		h.logger.Error("grpc get balance error", zap.Error(err))
		return nil, fmt.Errorf("grpc get balance error: %w", err)
	}
	response := &generatedPayments.GetDailyLikeBalanceResponse{
		Balance: int32(balance),
	}
	return response, nil
}

func (h *GRPCHandler) GetPurchasedLikeBalance(ctx context.Context,
	in *generatedPayments.GetPurchasedLikeBalanceRequest) (*generatedPayments.GetPurchasedLikeBalanceResponse, error) {
	userID := int(in.UserID)
	balance, err := h.uc.GetPurchasedLikesCount(ctx, userID)
	if err != nil {
		h.logger.Error("grpc get balance error", zap.Error(err))
		return nil, fmt.Errorf("grpc get balance error: %w", err)
	}
	response := &generatedPayments.GetPurchasedLikeBalanceResponse{
		Balance: int32(balance),
	}
	return response, nil
}

func (h *GRPCHandler) RefreshDailyLikeBalance(ctx context.Context,
	in *generatedPayments.RefreshDailyLikeBalanceRequest) (*generatedPayments.RefreshDailyLikeBalanceResponse, error) {
	err := h.uc.SetDailyLikeCountToAll(ctx, consts.DailyLikeLimit)
	if err != nil {
		h.logger.Error("grpc set daily like like count error", zap.Error(err))
		return nil, fmt.Errorf("grpc set daily like like count error: %w", err)
	}
	return &generatedPayments.RefreshDailyLikeBalanceResponse{}, nil
}

func (h *GRPCHandler) ChangeBalance(ctx context.Context,
	in *generatedPayments.ChangeBalanceRequest) (*generatedPayments.ChangeBalanceResponse, error) {
	userID := int(in.UserID)
	amount := int(in.Amount)

	err := h.uc.ChangeBalance(ctx, userID, amount)
	if err != nil {
		h.logger.Error("grpc change balance error", zap.Error(err))
		return nil, fmt.Errorf("grpc change balance error: %w", err)
	}
	return &generatedPayments.ChangeBalanceResponse{}, nil
}

func (h *GRPCHandler) CheckAndSpendLike(ctx context.Context,
	in *generatedPayments.CheckAndSpendLikeRequest) (*generatedPayments.CheckAndSpendLikeResponse, error) {
	userID := int(in.UserID)

	dailyLikes, err := h.uc.GetDailyLikesCount(ctx, userID)
	if err != nil {
		h.logger.Error("grpc get daily likes count error", zap.Error(err))
		return nil, fmt.Errorf("grpc get daily likes count error: %w", err)
	}
	if dailyLikes >= 1 {
		err = h.uc.ChangeDailyLikeCount(ctx, userID, -1)
		if err != nil {
			h.logger.Error("grpc change balance error", zap.Error(err))
			return nil, fmt.Errorf("grpc change balance error: %w", err)
		}
	} else {
		balance, err := h.uc.GetPurchasedLikesCount(ctx, userID)
		if err != nil {
			h.logger.Error("grpc get balance error", zap.Error(err))
			return nil, fmt.Errorf("grpc get balance error: %w", err)
		}
		if balance < 1 {
			return nil, fmt.Errorf("dont have likes: %w", err)
		}
		err = h.uc.ChangePurchasedLikeCount(ctx, userID, -1)
	}
	return &generatedPayments.CheckAndSpendLikeResponse{}, nil
}

func (h *GRPCHandler) ChangePurchasedLikesBalance(ctx context.Context,
	in *generatedPayments.ChangePurchasedLikesBalanceRequest) (*generatedPayments.ChangePurchasedLikesBalanceResponse, error) {
	userID := int(in.UserID)
	amount := int(in.Amount)

	err := h.uc.ChangePurchasedLikeCount(ctx, userID, amount)
	if err != nil {
		h.logger.Error("grpc change balance error", zap.Error(err))
		return nil, fmt.Errorf("grpc change balance error: %w", err)
	}
	return &generatedPayments.ChangePurchasedLikesBalanceResponse{}, nil
}

func (h *GRPCHandler) GetAllBalance(ctx context.Context,
	in *generatedPayments.GetAllBalanceRequest) (*generatedPayments.GetAllBalanceResponse, error) {
	userID := int(in.UserID)
	dailyLikes, err := h.uc.GetDailyLikesCount(ctx, userID)
	if err != nil {
		h.logger.Error("grpc get balance error", zap.Error(err))
		return nil, fmt.Errorf("grpc get balance error: %w", err)
	}
	purchasedLikes, err := h.uc.GetPurchasedLikesCount(ctx, userID)
	if err != nil {
		h.logger.Error("grpc get balance error", zap.Error(err))
		return nil, fmt.Errorf("grpc get balance error: %w", err)
	}
	moneyBalance, err := h.uc.GetBalance(ctx, userID)
	if err != nil {
		h.logger.Error("grpc get balance error", zap.Error(err))
		return nil, fmt.Errorf("grpc get balance error: %w", err)
	}
	response := &generatedPayments.GetAllBalanceResponse{
		DailyLikeBalance:     int32(dailyLikes),
		PurchasedLikeBalance: int32(purchasedLikes),
		MoneyBalance:         int32(moneyBalance),
	}
	return response, nil
}

func (h *GRPCHandler) CreateBalances(ctx context.Context,
	in *generatedPayments.CreateBalancesRequest) (*generatedPayments.CreateBalancesResponse, error) {
	userID := int(in.UserID)
	moneyAmount := int(in.MoneyAmount)
	dailyAmount := int(in.DailyAmount)
	purchasedAmount := int(in.PurchasedAmount)

	err := h.uc.AddBalance(ctx, userID, moneyAmount)
	if err != nil {
		return nil, fmt.Errorf("bad add balance error: %w", err)
	}
	err = h.uc.AddDailyLikesCount(ctx, userID, dailyAmount)
	if err != nil {
		return nil, fmt.Errorf("bad daily likes count error: %w", err)
	}
	err = h.uc.AddPurchasedLikesCount(ctx, userID, purchasedAmount)
	if err != nil {
		return nil, fmt.Errorf("bad purchase count error: %w", err)
	}
	return &generatedPayments.CreateBalancesResponse{}, nil
}

func (h *GRPCHandler) BuyLikes(ctx context.Context,
	in *generatedPayments.BuyLikesRequest) (*generatedPayments.BuyLikesResponse, error) {
	title := in.Title
	userID := int(in.UserID)

	product, err := h.uc.GetProduct(ctx, title)
	if err != nil {
		h.logger.Error("grpc get balance error", zap.Error(err))
		return nil, fmt.Errorf("grpc get balance error: %w", err)
	}
	h.logger.Info("product", zap.Any("product", product))
	err = h.uc.CheckBalance(ctx, userID, product.Price)
	if err != nil {
		h.logger.Error("grpc check balance error", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "Недостаточно средств")
	}
	spend := product.Count * (-1)
	if product.Count < 1 {
		h.logger.Info("grpc count < 1")
		return nil, status.Error(codes.InvalidArgument, "Суммы не хватает даже на один лайк")
	}
	err = h.uc.ChangeBalance(ctx, userID, spend)
	if err != nil {
		h.logger.Error("grpc change balance error", zap.Error(err))
		return nil, fmt.Errorf("grpc change balance error: %w", err)
	}
	err = h.uc.ChangePurchasedLikeCount(ctx, userID, product.Count)
	if err != nil {
		h.logger.Error("grpc change balance error", zap.Error(err))
		return nil, fmt.Errorf("grpc change balance error: %w", err)
	}
	response := &generatedPayments.BuyLikesResponse{}
	return response, nil
}

func (h *GRPCHandler) CreateProduct(ctx context.Context,
	in *generatedPayments.CreateProductRequest) (*generatedPayments.CreateProductResponse, error) {
	product := models.Product{
		Title:       in.Product.Title,
		Description: in.Product.Description,
		ImageLink:   in.Product.ImageLink,
		Price:       int(in.Product.Price),
		Count:       int(in.Product.Count),
	}
	id, err := h.uc.CreateProduct(ctx, product)
	if err != nil {
		h.logger.Error("grpc create product error", zap.Error(err))
		return nil, fmt.Errorf("grpc create product error: %w", err)
	}
	response := &generatedPayments.CreateProductResponse{ID: int32(id)}
	return response, nil
}

func (h *GRPCHandler) GetProducts(ctx context.Context,
	in *generatedPayments.GetProductsRequest) (*generatedPayments.GetProductsResponse, error) {
	products, err := h.uc.GetProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("grpc get products error: %w", err)
	}
	var grpcProducts []*generatedPayments.Product
	for _, product := range products {
		grpcProducts = append(grpcProducts, &generatedPayments.Product{
			Title:       product.Title,
			Description: product.Description,
			ImageLink:   product.ImageLink,
			Price:       int32(product.Price),
			Count:       int32(product.Count),
		})
	}
	response := &generatedPayments.GetProductsResponse{
		Products: grpcProducts,
	}
	return response, nil
}

func (h *GRPCHandler) AddAward(ctx context.Context, in *generatedPayments.AddAwardRequest) (*generatedPayments.AddAwardResponse, error) {
	award := models.Award{
		DayNumber: int(in.Award.DayNumber),
		Type:      in.Award.Type,
		Count:     int(in.Award.Count),
	}
	err := h.uc.AddAward(ctx, award)
	if err != nil {
		return nil, fmt.Errorf("grpc add award error: %w", err)
	}
	return &generatedPayments.AddAwardResponse{}, nil
}

func (h *GRPCHandler) GetAwards(ctx context.Context,
	in *generatedPayments.GetAwardsRequest) (*generatedPayments.GetAwardsResponse, error) {
	awards, err := h.uc.GetAwards(ctx)
	if err != nil {
		return nil, fmt.Errorf("grpc get awards error: %w", err)
	}
	var grpcAwards []*generatedPayments.Award
	for _, award := range awards {
		grpcAwards = append(grpcAwards, &generatedPayments.Award{
			DayNumber: int32(award.DayNumber),
			Type:      award.Type,
			Count:     int32(award.Count),
		})
	}
	response := &generatedPayments.GetAwardsResponse{
		Awards: grpcAwards,
	}
	return response, nil
}

func (h *GRPCHandler) UpdateActivity(ctx context.Context,
	in *generatedPayments.UpdateActivityRequest) (*generatedPayments.UpdateActivityResponse, error) {
	userID := int(in.UserID)

	answer, err := h.uc.UpdateActivity(ctx, userID)
	if err != nil {
		h.logger.Error("grpc update activity error", zap.Error(err))
		return nil, fmt.Errorf("grpc update activity error: %w", err)
	}
	response := &generatedPayments.UpdateActivityResponse{
		Answer: answer,
	}
	return response, nil
}

func (h *GRPCHandler) CreateActivity(ctx context.Context,
	in *generatedPayments.CreateActivityRequest) (*generatedPayments.CreateActivityResponse, error) {
	userID := int(in.UserID)
	err := h.uc.AddActivity(ctx, userID)
	if err != nil {
		h.logger.Error("grpc add activity error", zap.Error(err))
		return nil, fmt.Errorf("grpc add activity error: %w", err)
	}
	return &generatedPayments.CreateActivityResponse{}, nil
}

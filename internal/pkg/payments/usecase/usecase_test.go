package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/usecase/mocks"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"testing"
)

func TestUseCase(t *testing.T) {
	type args struct {
		ctx       context.Context
		userID    int
		amount    int
		title     string
		product   models.Product
		needMoney int
	}

	logger := zap.NewNop()

	ctx := context.Background()

	tests := []struct {
		name       string
		method     string
		args       args
		mockSetup  func(repo *mocks.MockRepository)
		wantErr    bool
		wantErrMsg string
		wantValInt int
		wantProds  []models.Product
		wantProd   models.Product
	}{
		{
			name:   "AddBalance success",
			method: "AddBalance",
			args:   args{ctx: ctx, userID: 1, amount: 100},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().AddBalance(ctx, 1, 100).Return(nil)
			},
			wantErr: false,
		},
		{
			name:       "AddBalance error",
			method:     "AddBalance",
			args:       args{ctx: ctx, userID: 1, amount: 100},
			mockSetup:  func(repo *mocks.MockRepository) { repo.EXPECT().AddBalance(ctx, 1, 100).Return(errors.New("db error")) },
			wantErr:    true,
			wantErrMsg: "failed to change balance: db error",
		},

		{
			name:   "AddDailyLikesCount success",
			method: "AddDailyLikesCount",
			args:   args{ctx: ctx, userID: 2, amount: 10},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().AddDailyLikeCount(ctx, 2, 10).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "AddDailyLikesCount error",
			method: "AddDailyLikesCount",
			args:   args{ctx: ctx, userID: 2, amount: 10},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().AddDailyLikeCount(ctx, 2, 10).Return(errors.New("db error"))
			},
			wantErr:    true,
			wantErrMsg: "failed to change balance: db error",
		},

		{
			name:   "AddPurchasedLikesCount success",
			method: "AddPurchasedLikesCount",
			args:   args{ctx: ctx, userID: 3, amount: 5},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().AddPurchasedLikeCount(ctx, 3, 5).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "AddPurchasedLikesCount error",
			method: "AddPurchasedLikesCount",
			args:   args{ctx: ctx, userID: 3, amount: 5},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().AddPurchasedLikeCount(ctx, 3, 5).Return(errors.New("db error"))
			},
			wantErr:    true,
			wantErrMsg: "failed to change balance: db error",
		},

		{
			name:   "ChangeBalance success",
			method: "ChangeBalance",
			args:   args{ctx: ctx, userID: 4, amount: 50},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().ChangeBalance(ctx, 4, 50).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "ChangeBalance error",
			method: "ChangeBalance",
			args:   args{ctx: ctx, userID: 4, amount: 50},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().ChangeBalance(ctx, 4, 50).Return(errors.New("db error"))
			},
			wantErr:    true,
			wantErrMsg: "failed to change balance: db error",
		},

		{
			name:   "ChangeDailyLikeCount success",
			method: "ChangeDailyLikeCount",
			args:   args{ctx: ctx, userID: 5, amount: 20},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().ChangeDailyLikeCount(ctx, 5, 20).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "ChangeDailyLikeCount error",
			method: "ChangeDailyLikeCount",
			args:   args{ctx: ctx, userID: 5, amount: 20},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().ChangeDailyLikeCount(ctx, 5, 20).Return(errors.New("db error"))
			},
			wantErr:    true,
			wantErrMsg: "failed to change daily like count: db error",
		},

		{
			name:   "ChangePurchasedLikeCount success",
			method: "ChangePurchasedLikeCount",
			args:   args{ctx: ctx, userID: 6, amount: 30},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().ChangePurchasedLikeCount(ctx, 6, 30).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "ChangePurchasedLikeCount error",
			method: "ChangePurchasedLikeCount",
			args:   args{ctx: ctx, userID: 6, amount: 30},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().ChangePurchasedLikeCount(ctx, 6, 30).Return(errors.New("db error"))
			},
			wantErr:    true,
			wantErrMsg: "failed to change purchased like count: db error",
		},

		{
			name:   "SetBalance success",
			method: "SetBalance",
			args:   args{ctx: ctx, userID: 7, amount: 500},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().SetBalance(ctx, 7, 500).Return(nil)
			},
			wantErr: false,
		},
		{
			name:       "SetBalance error",
			method:     "SetBalance",
			args:       args{ctx: ctx, userID: 7, amount: 500},
			mockSetup:  func(repo *mocks.MockRepository) { repo.EXPECT().SetBalance(ctx, 7, 500).Return(errors.New("db error")) },
			wantErr:    true,
			wantErrMsg: "failed to change purchased like count: db error",
		},

		{
			name:   "SetDailyLikeCount success",
			method: "SetDailyLikeCount",
			args:   args{ctx: ctx, userID: 8, amount: 5},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().SetDailyLikesCount(ctx, 8, 5).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "SetDailyLikeCount error",
			method: "SetDailyLikeCount",
			args:   args{ctx: ctx, userID: 8, amount: 5},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().SetDailyLikesCount(ctx, 8, 5).Return(errors.New("db error"))
			},
			wantErr:    true,
			wantErrMsg: "failed to change purchased like count: db error",
		},

		{
			name:   "SetDailyLikeCountToAll success",
			method: "SetDailyLikeCountToAll",
			args:   args{ctx: ctx, amount: 5},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().SetDailyLikesCountToAll(ctx, 5).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "SetDailyLikeCountToAll error",
			method: "SetDailyLikeCountToAll",
			args:   args{ctx: ctx, amount: 5},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().SetDailyLikesCountToAll(ctx, 5).Return(errors.New("db error"))
			},
			wantErr:    true,
			wantErrMsg: "failed to change purchased like count: db error",
		},

		{
			name:   "SetPurchasedLikeCount success",
			method: "SetPurchasedLikeCount",
			args:   args{ctx: ctx, userID: 9, amount: 15},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().SetPurchasedLikesCount(ctx, 9, 15).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "SetPurchasedLikeCount error",
			method: "SetPurchasedLikeCount",
			args:   args{ctx: ctx, userID: 9, amount: 15},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().SetPurchasedLikesCount(ctx, 9, 15).Return(errors.New("db error"))
			},
			wantErr:    true,
			wantErrMsg: "failed to change purchased like count: db error",
		},

		{
			name:   "GetBalance success",
			method: "GetBalance",
			args:   args{ctx: ctx, userID: 10},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().GetBalance(ctx, 10).Return(1000, nil)
			},
			wantErr:    false,
			wantValInt: 1000,
		},
		{
			name:       "GetBalance error",
			method:     "GetBalance",
			args:       args{ctx: ctx, userID: 10},
			mockSetup:  func(repo *mocks.MockRepository) { repo.EXPECT().GetBalance(ctx, 10).Return(-1, errors.New("db error")) },
			wantErr:    true,
			wantErrMsg: "failed to get balance: db error",
			wantValInt: -1,
		},

		{
			name:   "GetDailyLikesCount success",
			method: "GetDailyLikesCount",
			args:   args{ctx: ctx, userID: 11},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().GetDailyLikesCount(ctx, 11).Return(50, nil)
			},
			wantValInt: 50,
		},
		{
			name:   "GetDailyLikesCount error",
			method: "GetDailyLikesCount",
			args:   args{ctx: ctx, userID: 11},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().GetDailyLikesCount(ctx, 11).Return(-1, errors.New("db error"))
			},
			wantErr:    true,
			wantErrMsg: "failed to get balance: db error",
			wantValInt: -1,
		},

		{
			name:   "GetPurchasedLikesCount success",
			method: "GetPurchasedLikesCount",
			args:   args{ctx: ctx, userID: 12},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().GetPurchasedLikesCount(ctx, 12).Return(25, nil)
			},
			wantValInt: 25,
		},
		{
			name:   "GetPurchasedLikesCount error",
			method: "GetPurchasedLikesCount",
			args:   args{ctx: ctx, userID: 12},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().GetPurchasedLikesCount(ctx, 12).Return(-1, errors.New("db error"))
			},
			wantErr:    true,
			wantErrMsg: "failed to get balance: db error",
			wantValInt: -1,
		},

		//{
		//	name:   "CreateProduct success",
		//	method: "CreateProduct",
		//	args:   args{ctx: ctx, product: models.Product{Title: "prod", Price: 100}},
		//	mockSetup: func(repo *mocks.MockRepository) {
		//		repo.EXPECT().CreateProduct(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, p models.Product) (int, error) {
		//			if p.Title != "prod" || p.Price != 100 || p.ImageLink != p.Title+".png" {
		//				return -1, fmt.Errorf("unexpected product")
		//			}
		//			return 1, nil
		//		})
		//	},
		//	wantValInt: 1,
		//},

		{
			name:   "GetProduct success",
			method: "GetProduct",
			args:   args{ctx: ctx, title: "prod"},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().GetProduct(ctx, "prod").Return(models.Product{Title: "prod", Price: 100}, nil)
			},
			wantProd: models.Product{Title: "prod", Price: 100},
		},
		{
			name:   "GetProduct error",
			method: "GetProduct",
			args:   args{ctx: ctx, title: "prod"},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().GetProduct(ctx, "prod").Return(models.Product{}, errors.New("db error"))
			},
			wantErr:    true,
			wantErrMsg: "failed to create product: db error",
		},

		{
			name:   "UpdateProduct success",
			method: "UpdateProduct",
			args:   args{ctx: ctx, title: "prod", product: models.Product{Title: "prod", Price: 50}},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().UpdateProduct(ctx, "prod", gomock.Any()).DoAndReturn(func(ctx context.Context, title string, p models.Product) error {
					if p.Price != 50 {
						return fmt.Errorf("unexpected price")
					}
					return nil
				})
			},
		},
		{
			name:       "UpdateProduct bad price",
			method:     "UpdateProduct",
			args:       args{ctx: ctx, title: "prod", product: models.Product{Title: "prod", Price: -20}},
			mockSetup:  func(repo *mocks.MockRepository) {},
			wantErr:    true,
			wantErrMsg: "invalid price: -20",
		},
		{
			name:   "UpdateProduct error",
			method: "UpdateProduct",
			args:   args{ctx: ctx, title: "prod", product: models.Product{Title: "prod", Price: 10}},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().UpdateProduct(ctx, "prod", gomock.Any()).Return(errors.New("db error"))
			},
			wantErr:    true,
			wantErrMsg: "failed to create product: db error",
		},

		{
			name:   "CheckBalance success",
			method: "CheckBalance",
			args:   args{ctx: ctx, userID: 13, needMoney: 500},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().GetBalance(ctx, 13).Return(1000, nil)
			},
			wantErr: false,
		},
		{
			name:       "CheckBalance get error",
			method:     "CheckBalance",
			args:       args{ctx: ctx, userID: 13, needMoney: 500},
			mockSetup:  func(repo *mocks.MockRepository) { repo.EXPECT().GetBalance(ctx, 13).Return(-1, errors.New("db error")) },
			wantErr:    true,
			wantErrMsg: "failed to get balance: db error",
		},
		{
			name:       "CheckBalance insufficient funds",
			method:     "CheckBalance",
			args:       args{ctx: ctx, userID: 13, needMoney: 2000},
			mockSetup:  func(repo *mocks.MockRepository) { repo.EXPECT().GetBalance(ctx, 13).Return(1000, nil) },
			wantErr:    true,
			wantErrMsg: "Недостаточно средств",
		},

		{
			name:   "GetProducts success",
			method: "GetProducts",
			args:   args{ctx: ctx},
			mockSetup: func(repo *mocks.MockRepository) {
				repo.EXPECT().GetProducts(ctx).Return([]models.Product{
					{Title: "p1", Price: 10},
					{Title: "p2", Price: 20},
				}, nil)
			},
			wantProds: []models.Product{
				{Title: "p1", Price: 10},
				{Title: "p2", Price: 20},
			},
		},
		{
			name:       "GetProducts error",
			method:     "GetProducts",
			args:       args{ctx: ctx},
			mockSetup:  func(repo *mocks.MockRepository) { repo.EXPECT().GetProducts(ctx).Return(nil, errors.New("db error")) },
			wantErr:    true,
			wantErrMsg: "failed to create product: db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			repo := mocks.NewMockRepository(mockCtrl)
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}
			uc := New(repo, logger)

			var (
				err    error
				valInt int
				prod   models.Product
				prods  []models.Product
			)

			switch tt.method {
			case "AddBalance":
				err = uc.AddBalance(tt.args.ctx, tt.args.userID, tt.args.amount)
			case "AddDailyLikesCount":
				err = uc.AddDailyLikesCount(tt.args.ctx, tt.args.userID, tt.args.amount)
			case "AddPurchasedLikesCount":
				err = uc.AddPurchasedLikesCount(tt.args.ctx, tt.args.userID, tt.args.amount)
			case "ChangeBalance":
				err = uc.ChangeBalance(tt.args.ctx, tt.args.userID, tt.args.amount)
			case "ChangeDailyLikeCount":
				err = uc.ChangeDailyLikeCount(tt.args.ctx, tt.args.userID, tt.args.amount)
			case "ChangePurchasedLikeCount":
				err = uc.ChangePurchasedLikeCount(tt.args.ctx, tt.args.userID, tt.args.amount)
			case "SetBalance":
				err = uc.SetBalance(tt.args.ctx, tt.args.userID, tt.args.amount)
			case "SetDailyLikeCount":
				err = uc.SetDailyLikeCount(tt.args.ctx, tt.args.userID, tt.args.amount)
			case "SetDailyLikeCountToAll":
				err = uc.SetDailyLikeCountToAll(tt.args.ctx, tt.args.amount)
			case "SetPurchasedLikeCount":
				err = uc.SetPurchasedLikeCount(tt.args.ctx, tt.args.userID, tt.args.amount)
			case "GetBalance":
				valInt, err = uc.GetBalance(tt.args.ctx, tt.args.userID)
			case "GetDailyLikesCount":
				valInt, err = uc.GetDailyLikesCount(tt.args.ctx, tt.args.userID)
			case "GetPurchasedLikesCount":
				valInt, err = uc.GetPurchasedLikesCount(tt.args.ctx, tt.args.userID)
			case "CreateProduct":
				valInt, err = uc.CreateProduct(tt.args.ctx, tt.args.product)
			case "GetProduct":
				prod, err = uc.GetProduct(tt.args.ctx, tt.args.title)
			case "UpdateProduct":
				err = uc.UpdateProduct(tt.args.ctx, tt.args.title, tt.args.product)
			case "CheckBalance":
				err = uc.CheckBalance(tt.args.ctx, tt.args.userID, tt.args.needMoney)
			case "GetProducts":
				prods, err = uc.GetProducts(tt.args.ctx)
			default:
				t.Fatalf("unknown method: %s", tt.method)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error mismatch: got err=%v, wantErr=%v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && tt.wantErrMsg != "" {
				if err.Error() != tt.wantErrMsg {
					t.Errorf("error message mismatch: got %v, want %v", err.Error(), tt.wantErrMsg)
				}
			}

			if tt.method == "GetBalance" ||
				tt.method == "GetDailyLikesCount" ||
				tt.method == "GetPurchasedLikesCount" ||
				tt.method == "CreateProduct" {
				if valInt != tt.wantValInt {
					t.Errorf("returned int mismatch: got %v, want %v", valInt, tt.wantValInt)
				}
			}

			if tt.method == "GetProduct" {
				if prod != tt.wantProd {
					t.Errorf("product mismatch: got %+v, want %+v", prod, tt.wantProd)
				}
			}

			if tt.method == "GetProducts" {
				if len(prods) != len(tt.wantProds) {
					t.Errorf("products length mismatch: got %d, want %d", len(prods), len(tt.wantProds))
				} else {
					for i := range prods {
						if prods[i] != tt.wantProds[i] {
							t.Errorf("product at %d mismatch: got %+v, want %+v", i, prods[i], tt.wantProds[i])
						}
					}
				}
			}
		})
	}
}

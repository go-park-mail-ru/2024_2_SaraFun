package getcurrentprofile

import (
	"context"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedPayments "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc/gen"
	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
)

//go:generate easyjson -all handler.go

//go:generate mockgen -destination=./mocks/mock_ImageService.go -package=sign_up_mocks . ImageService
type ImageService interface {
	GetImageLinksByUserId(ctx context.Context, id int) ([]models.Image, error)
}

//go:generate mockgen -destination=./mocks/mock_ProfileService.go -package=sign_up_mocks . ProfileService

//go:generate mockgen -destination=./mocks/mock_UserService.go -package=sign_up_mocks . UserService

type PersonalitiesClient interface {
	GetProfile(ctx context.Context,
		in *generatedPersonalities.GetProfileRequest) (*generatedPersonalities.GetProfileResponse, error)
	GetProfileIdByUserId(ctx context.Context,
		in *generatedPersonalities.GetProfileIDByUserIDRequest) (*generatedPersonalities.GetProfileIDByUserIDResponse, error)
}

//go:generate mockgen -destination=./mocks/mock_SessionService.go -package=sign_up_mocks . SessionService

type SessionClient interface {
	GetUserIDBySessionID(ctx context.Context, in *generatedAuth.GetUserIDBySessionIDRequest) (*generatedAuth.GetUserIDBYSessionIDResponse, error)
}

type Response struct {
	Username              string         `json:"username"`
	Profile               models.Profile `json:"profile"`
	Images                []models.Image `json:"images"`
	MoneyBalance          int            `json:"money_balance"`
	DailyLikesBalance     int            `json:"daily_likes_balance"`
	PurchasedLikesBalance int            `json:"purchased_likes_balance"`
}

//easyjson:skip
type Handler struct {
	imageService        ImageService
	personalitiesClient generatedPersonalities.PersonalitiesClient
	sessionClient       generatedAuth.AuthClient
	paymentsClient      generatedPayments.PaymentClient
	logger              *zap.Logger
}

func NewHandler(imageService ImageService, personalitiesClient generatedPersonalities.PersonalitiesClient,
	sessionClient generatedAuth.AuthClient, paymentsClient generatedPayments.PaymentClient, logger *zap.Logger) *Handler {
	return &Handler{
		imageService:        imageService,
		personalitiesClient: personalitiesClient,
		sessionClient:       sessionClient,
		paymentsClient:      paymentsClient,
		logger:              logger}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	req_id := ctx.Value(consts.RequestIDKey).(string)
	h.logger.Info("Handling request", zap.String("request_id", req_id))

	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		h.logger.Error("error getting session cookie", zap.Error(err))
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	getUserRequest := &generatedAuth.GetUserIDBySessionIDRequest{
		SessionID: cookie.Value,
	}

	userId, err := h.sessionClient.GetUserIDBySessionID(ctx, getUserRequest)
	h.logger.Info("GetUserByCookie", zap.Int("userid", int(userId.UserId)))
	if err != nil {
		h.logger.Error("error getting user id", zap.Error(err))
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}
	getUsernameRequest := &generatedPersonalities.GetUsernameByUserIDRequest{UserID: userId.UserId}

	username, err := h.personalitiesClient.GetUsernameByUserID(ctx, getUsernameRequest)
	if err != nil {
		h.logger.Error("error getting username", zap.Error(err))
		http.Error(w, "user username not found", http.StatusUnauthorized)
		return
	}

	getProfileByUserRequest := &generatedPersonalities.GetProfileIDByUserIDRequest{UserID: userId.UserId}
	profileId, err := h.personalitiesClient.GetProfileIDByUserID(ctx, getProfileByUserRequest)
	h.logger.Info("GetProfileByUser", zap.Int("profileid", int(profileId.ProfileID)))
	if err != nil {
		h.logger.Error("getprofileidbyuserid error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var links []models.Image
	links, err = h.imageService.GetImageLinksByUserId(ctx, int(userId.UserId))
	if err != nil {
		h.logger.Error("getimagelinkbyuserid error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	getProfileRequest := &generatedPersonalities.GetProfileRequest{Id: profileId.ProfileID}
	profile, err := h.personalitiesClient.GetProfile(ctx, getProfileRequest)
	if err != nil {
		h.logger.Error("getprofileerror", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	getBalancesReq := &generatedPayments.GetAllBalanceRequest{UserID: userId.UserId}
	balance, err := h.paymentsClient.GetAllBalance(ctx, getBalancesReq)
	if err != nil {
		h.logger.Error("getbalanceserror", zap.Error(err))
		http.Error(w, "не удалось получить баланс", http.StatusInternalServerError)
		return
	}

	profileResponse := models.Profile{
		ID:           int(profile.Profile.ID),
		FirstName:    profile.Profile.FirstName,
		LastName:     profile.Profile.LastName,
		Age:          int(profile.Profile.Age),
		Gender:       profile.Profile.Gender,
		Target:       profile.Profile.Target,
		About:        profile.Profile.About,
		BirthdayDate: profile.Profile.BirthDate,
	}
	response := Response{
		Username:              username.Username,
		Profile:               profileResponse,
		Images:                links,
		MoneyBalance:          int(balance.MoneyBalance),
		DailyLikesBalance:     int(balance.DailyLikeBalance),
		PurchasedLikesBalance: int(balance.PurchasedLikeBalance),
	}
	jsonData, err := easyjson.Marshal(response)
	if err != nil {
		h.logger.Error("json marshal error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		h.logger.Error("write jsonData error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Info("getprofile success")

}

package addreaction

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedCommunications "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen"
	generatedPayments "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc/gen"
	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
)

//go:generate mockgen -destination=./mocks/mock_ReactionService.go -package=sign_up_mocks . ReactionService

//go:generate mockgen -destination=./mocks/mock_SessionService.go -package=sign_up_mocks . SessionService

//go:generate easyjson -all handler.go

//go:generate mockgen -destination=./mocks/mock_ImageService.go -package=sign_up_mocks . ImageService
type ImageService interface {
	GetFirstImage(ctx context.Context, userID int) (models.Image, error)
}

type WebSocketService interface {
	SendNotification(ctx context.Context, receiverID int, authorUsername string, authorImageLink string) error
}

type Request struct {
	Receiver string `json:"receiver"`
	Type     bool   `json:"type"`
}

//easyjson:skip
type Handler struct {
	reactionClient       generatedCommunications.CommunicationsClient
	SessionClient        generatedAuth.AuthClient
	personalitiesClient  generatedPersonalities.PersonalitiesClient
	communicationsClient generatedCommunications.CommunicationsClient
	paymentsClient       generatedPayments.PaymentClient
	imageService         ImageService
	wsService            WebSocketService
	logger               *zap.Logger
}

func NewHandler(reactionClient generatedCommunications.CommunicationsClient,
	sessionClient generatedAuth.AuthClient, personalitiesClient generatedPersonalities.PersonalitiesClient,
	communicationsClient generatedCommunications.CommunicationsClient,
	paymentsClient generatedPayments.PaymentClient, imageService ImageService,
	wsService WebSocketService, logger *zap.Logger) *Handler {
	return &Handler{
		reactionClient:       reactionClient,
		SessionClient:        sessionClient,
		personalitiesClient:  personalitiesClient,
		communicationsClient: communicationsClient,
		paymentsClient:       paymentsClient,
		imageService:         imageService,
		wsService:            wsService,
		logger:               logger,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req_id := ctx.Value(consts.RequestIDKey).(string)
	h.logger.Info("Handling request", zap.String("request_id", req_id))
	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		h.logger.Error("AddReaction Handler: bad getting cookie ", zap.Error(err))
		http.Error(w, "Вы не авторизованы", http.StatusUnauthorized)
		return
	}
	getUserIdRequest := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: cookie.Value}
	userId, err := h.SessionClient.GetUserIDBySessionID(ctx, getUserIdRequest)
	if err != nil {
		h.logger.Error("AddReaction Handler: bad getting user id ", zap.Error(err))
		http.Error(w, "Вы не авторизованы", http.StatusUnauthorized)
		return
	}
	var reaction models.Reaction
	var request Request
	//err = json.NewDecoder(r.Body).Decode(&reaction)
	err = easyjson.UnmarshalFromReader(r.Body, &request)
	if err != nil {
		h.logger.Error("AddReaction Handler: bad unmarshal", zap.Error(err))
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	getUsernameReq := &generatedPersonalities.GetUserIDByUsernameRequest{Username: request.Receiver}
	receiverID, err := h.personalitiesClient.GetUserIDByUsername(ctx, getUsernameReq)
	if err != nil {
		h.logger.Error("getting receiver ID error", zap.Error(err))
		http.Error(w, "Получатель не найден", http.StatusBadRequest)
		return
	}

	if userId.UserId == receiverID.UserID {
		h.logger.Error("самолайк", zap.Error(errors.New("самолайк")))
		http.Error(w, "Себя лайкать нельзя!", http.StatusBadRequest)
		return
	}

	reaction.Author = int(userId.UserId)
	reaction.Receiver = int(receiverID.UserID)
	reaction.Type = request.Type
	react := &generatedCommunications.Reaction{
		ID:       int32(reaction.Id),
		Author:   int32(reaction.Author),
		Receiver: int32(reaction.Receiver),
		Type:     reaction.Type,
	}
	checkAndSpendReq := &generatedPayments.CheckAndSpendLikeRequest{UserID: userId.UserId}
	_, err = h.paymentsClient.CheckAndSpendLike(ctx, checkAndSpendReq)
	if err != nil {
		h.logger.Error("AddReaction Handler: bad checking and spend like", zap.Error(err))
		http.Error(w, "у вас нет лайков", http.StatusBadRequest)
		return
	}

	addReactionRequest := &generatedCommunications.AddReactionRequest{Reaction: react}
	_, err = h.reactionClient.AddReaction(ctx, addReactionRequest)
	if err != nil {
		h.logger.Error("AddReaction Handler: error adding reaction", zap.Error(err))
		http.Error(w, "Что-то пошло не так :(", http.StatusInternalServerError)
		return
	}

	checkMatchExistsRequest := &generatedCommunications.CheckMatchExistsRequest{
		FirstUser:  int32(reaction.Author),
		SecondUser: receiverID.UserID,
	}

	checkMatchExistsResponse, err := h.communicationsClient.CheckMatchExists(ctx, checkMatchExistsRequest)
	if err != nil {
		h.logger.Error("AddReaction Handler: error checking match exists", zap.Error(err))
		http.Error(w, "Что-то пошло не так :(", http.StatusInternalServerError)
		return
	}
	if checkMatchExistsResponse.Exists {
		firstReq := &generatedPersonalities.GetUsernameByUserIDRequest{UserID: int32(reaction.Author)}
		secondReq := &generatedPersonalities.GetUsernameByUserIDRequest{UserID: receiverID.UserID}

		firstUsername, err := h.personalitiesClient.GetUsernameByUserID(ctx, firstReq)
		if err != nil {
			h.logger.Error("AddReaction Handler: error getting first username", zap.Error(err))
			http.Error(w, "Что-то пошло не так :(", http.StatusInternalServerError)
			return
		}
		secondUsername, err := h.personalitiesClient.GetUsernameByUserID(ctx, secondReq)
		if err != nil {
			h.logger.Error("AddReaction Handler: error getting second username", zap.Error(err))
			http.Error(w, "Что-то пошло не так :(", http.StatusInternalServerError)
			return
		}

		firstUserImage, err := h.imageService.GetFirstImage(ctx, reaction.Author)
		if err != nil {
			h.logger.Error("AddReaction Handler: error getting first image", zap.Error(err))
			firstUserImage.Link = ""
		}
		secondUserImage, err := h.imageService.GetFirstImage(ctx, int(receiverID.UserID))
		if err != nil {
			h.logger.Error("AddReaction Handler: error getting second image", zap.Error(err))
			secondUserImage.Link = ""
		}

		err = h.wsService.SendNotification(ctx, int(receiverID.UserID), firstUsername.Username, firstUserImage.Link)
		if err != nil {
			h.logger.Error("AddReaction Handler: error sending notification", zap.Error(err))
			http.Error(w, "Что-то пошло не так :(", http.StatusInternalServerError)
			return
		}
		err = h.wsService.SendNotification(ctx, reaction.Author, secondUsername.Username, secondUserImage.Link)
		if err != nil {
			h.logger.Error("AddReaction Handler: error sending notification", zap.Error(err))
			http.Error(w, "Что-то пошло не так :(", http.StatusInternalServerError)
			return
		}

	}

	h.logger.Info("AddReaction Handler: added reaction", zap.Any("reaction", reaction))
	fmt.Fprintf(w, "Вы успешно лайкнули пользователя!")
}

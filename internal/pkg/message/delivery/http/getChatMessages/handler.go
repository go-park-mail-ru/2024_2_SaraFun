package getChatMessages

import (
	"context"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedMessage "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/message/delivery/grpc/gen"
	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

//go:generate easyjson -all handler.go

type ResponseMessage struct {
	Body string `json:"body"`
	Self bool   `json:"self"`
	Time string `json:"time"`
}

//go:generate mockgen -destination=./mocks/mock_ImageService.go -package=sign_up_mocks . ImageService
type ImageService interface {
	GetImageLinksByUserId(ctx context.Context, id int) ([]models.Image, error)
}

type Response struct {
	Username string            `json:"username"`
	Profile  models.Profile    `json:"profile"`
	Messages []ResponseMessage `json:"messages"`
	Images   []models.Image    `json:"images"`
}

//easyjson:skip
type Handler struct {
	authClient          generatedAuth.AuthClient
	messageClient       generatedMessage.MessageClient
	personalitiesClient generatedPersonalities.PersonalitiesClient
	imageService        ImageService
	logger              *zap.Logger
}

func NewHandler(authClient generatedAuth.AuthClient, messageClient generatedMessage.MessageClient,
	personalitiesClient generatedPersonalities.PersonalitiesClient, imageService ImageService, logger *zap.Logger) *Handler {
	return &Handler{
		authClient:          authClient,
		messageClient:       messageClient,
		personalitiesClient: personalitiesClient,
		imageService:        imageService,
		logger:              logger,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var resp Response
	req_id := ctx.Value(consts.RequestIDKey)
	h.logger.Info("Handle request", zap.String("request_id", req_id.(string)))

	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		h.logger.Error("bad cookie", zap.Error(err))
		http.Error(w, "Вы не авторизованы", http.StatusUnauthorized)
		return
	}
	getUserIDRequest := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: cookie.Value}
	userId, err := h.authClient.GetUserIDBySessionID(ctx, getUserIDRequest)
	if err != nil {
		h.logger.Error("dont get user by session id", zap.Error(err))
		http.Error(w, "Вы не авторизованы", http.StatusUnauthorized)
		return
	}

	firstUserID := userId.UserId
	secondUserID, err := strconv.Atoi(r.URL.Query().Get("userID"))
	if err != nil {
		h.logger.Error("dont get user id", zap.Error(err))
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	getUsernameRequest := &generatedPersonalities.GetUsernameByUserIDRequest{UserID: int32(secondUserID)}
	username, err := h.personalitiesClient.GetUsernameByUserID(ctx, getUsernameRequest)
	if err != nil {
		h.logger.Error("dont get username by userID", zap.Error(err))
		http.Error(w, "Что-то пошло не так :(", http.StatusInternalServerError)
		return
	}
	resp.Username = username.Username

	getProfileRequestID := &generatedPersonalities.GetProfileIDByUserIDRequest{UserID: int32(secondUserID)}
	secondProfileID, err := h.personalitiesClient.GetProfileIDByUserID(ctx, getProfileRequestID)
	if err != nil {
		h.logger.Error("dont get user profile", zap.Error(err))
		http.Error(w, "Что-то пошло не так :(", http.StatusInternalServerError)
		return
	}
	getProfileRequest := &generatedPersonalities.GetProfileRequest{Id: secondProfileID.ProfileID}
	secondProfile, err := h.personalitiesClient.GetProfile(ctx, getProfileRequest)
	if err != nil {
		h.logger.Error("dont get user profile", zap.Error(err))
		http.Error(w, "Что-то пошло не так :(", http.StatusInternalServerError)
		return
	}
	respProfile := models.Profile{
		ID:        int(secondProfileID.ProfileID),
		FirstName: secondProfile.Profile.FirstName,
		LastName:  secondProfile.Profile.LastName,
		Age:       int(secondProfile.Profile.Age),
		Gender:    secondProfile.Profile.Gender,
		Target:    secondProfile.Profile.Target,
		About:     secondProfile.Profile.About,
	}

	resp.Profile = respProfile

	getChatMessagesRequest := &generatedMessage.GetChatMessagesRequest{
		FirstUserID:  firstUserID,
		SecondUserID: int32(secondUserID),
	}
	msgs, err := h.messageClient.GetChatMessages(ctx, getChatMessagesRequest)
	if err != nil {
		h.logger.Error("dont get chat messages", zap.Error(err))
		http.Error(w, "Что-то пошло не так :(", http.StatusInternalServerError)
		return
	}
	var responseMessages []ResponseMessage

	for _, msg := range msgs.Messages {
		responseMessage := ResponseMessage{
			Body: msg.Body,
			Time: msg.Time,
			Self: true,
		}
		if msg.Author != userId.UserId {
			responseMessage.Self = false
		}
		responseMessages = append(responseMessages, responseMessage)
	}

	resp.Messages = responseMessages

	var links []models.Image
	links, err = h.imageService.GetImageLinksByUserId(ctx, secondUserID)
	if err != nil {
		h.logger.Error("getimagelinkbyuserid error", zap.Error(err))
		http.Error(w, "Что-то пошло не так :(", http.StatusInternalServerError)
		return
	}

	resp.Images = links

	w.Header().Set("Content-Type", "application/json")
	jsonData, err := easyjson.Marshal(resp)
	if err != nil {
		h.logger.Error("dont marshal response", zap.Error(err))
		http.Error(w, "Что-то пошло не так :(", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonData)
	if err != nil {
		h.logger.Error("dont write response", zap.Error(err))
		http.Error(w, "Что-то пошло не так :(", http.StatusInternalServerError)
		return
	}
	h.logger.Info("getChatMessages success")
}

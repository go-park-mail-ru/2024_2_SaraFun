package getChatMessages

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedMessage "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/message/delivery/grpc/gen"
	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type ResponseMessage struct {
	Body string `json:"body"`
	Self bool   `json:"self"`
	Time string `json:"time"`
}

type Response struct {
	Username string            `json:"username"`
	Profile  models.Profile    `json:"profile"`
	Messages []ResponseMessage `json:"messages"`
}

type Handler struct {
	authClient          generatedAuth.AuthClient
	messageClient       generatedMessage.MessageClient
	personalitiesClient generatedPersonalities.PersonalitiesClient
	logger              *zap.Logger
}

func NewHandler(authClient generatedAuth.AuthClient, messageClient generatedMessage.MessageClient, personalitiesClient generatedPersonalities.PersonalitiesClient, logger *zap.Logger) *Handler {
	return &Handler{
		authClient:          authClient,
		messageClient:       messageClient,
		personalitiesClient: personalitiesClient,
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
		http.Error(w, "bad cookie", http.StatusUnauthorized)
		return
	}
	getUserIDRequest := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: cookie.Value}
	userId, err := h.authClient.GetUserIDBySessionID(ctx, getUserIDRequest)
	if err != nil {
		h.logger.Error("dont get user by session id", zap.Error(err))
		http.Error(w, "dont get user by session id", http.StatusUnauthorized)
	}

	firstUserID := userId.UserId
	//err = json.NewDecoder(r.Body).Decode(&secondUserID)
	//if err != nil {
	//	h.logger.Error("dont decode secondUser id", zap.Error(err))
	//	http.Error(w, "dont decode secondUser id", http.StatusBadRequest)
	//}
	secondUserID, err := strconv.Atoi(r.URL.Query().Get("userID"))
	if err != nil {
		h.logger.Error("dont get user id", zap.Error(err))
		http.Error(w, "dont get user id", http.StatusBadRequest)
		return
	}

	getUsernameRequest := &generatedPersonalities.GetUsernameByUserIDRequest{UserID: int32(secondUserID)}
	username, err := h.personalitiesClient.GetUsernameByUserID(ctx, getUsernameRequest)
	if err != nil {
		h.logger.Error("dont get username by userID", zap.Error(err))
		http.Error(w, "dont get username by user id", http.StatusInternalServerError)
		return
	}
	resp.Username = username.Username

	getProfileRequestID := &generatedPersonalities.GetProfileIDByUserIDRequest{UserID: int32(secondUserID)}
	secondProfileID, err := h.personalitiesClient.GetProfileIDByUserID(ctx, getProfileRequestID)
	if err != nil {
		h.logger.Error("dont get user profile", zap.Error(err))
		http.Error(w, "dont get user profile", http.StatusInternalServerError)
		return
	}
	getProfileRequest := &generatedPersonalities.GetProfileRequest{Id: secondProfileID.ProfileID}
	secondProfile, err := h.personalitiesClient.GetProfile(ctx, getProfileRequest)
	if err != nil {
		h.logger.Error("dont get user profile", zap.Error(err))
		http.Error(w, "dont get user profile", http.StatusInternalServerError)
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
		http.Error(w, "dont get chat messages", http.StatusBadRequest)
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

	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(resp)
	if err != nil {
		h.logger.Error("dont marshal response", zap.Error(err))
		http.Error(w, "dont marshal response", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonData)
	if err != nil {
		h.logger.Error("dont write response", zap.Error(err))
		http.Error(w, "dont write response", http.StatusInternalServerError)
		return
	}
	h.logger.Info("getChatMessages success")
}

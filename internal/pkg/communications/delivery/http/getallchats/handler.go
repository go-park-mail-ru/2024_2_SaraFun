package getallchats

import (
	"context"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedCommunications "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen"
	generatedMessage "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/message/delivery/grpc/gen"
	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
	"sort"
	"time"
)

//go:generate mockgen -destination=./mocks/mock_ReactionService.go -package=sign_up_mocks . ReactionService

//go:generate easyjson -all handler.go

type CommunicationsClient interface {
	GetMatchList(ctx context.Context,
		in *generatedCommunications.GetMatchListRequest) (*generatedCommunications.GetMatchListResponse, error)
}

//go:generate mockgen -destination=./mocks/mock_SessionService.go -package=sign_up_mocks . SessionService

type SessionClient interface {
	GetUserIDBySessionID(ctx context.Context, in *generatedAuth.GetUserIDBySessionIDRequest) (*generatedAuth.GetUserIDBYSessionIDResponse, error)
}

//go:generate mockgen -destination=./mocks/mock_ProfileService.go -package=sign_up_mocks . ProfileService

//go:generate mockgen -destination=./mocks/mock_UserService.go -package=sign_up_mocks . UserService

type PersonalitiesClient interface {
	GetUsernameByUserId(ctx context.Context,
		in *generatedPersonalities.GetUsernameByUserIDRequest) (*generatedPersonalities.GetUsernameByUserIDResponse, error)
	GetProfile(ctx context.Context,
		in *generatedPersonalities.GetProfileRequest) (*generatedPersonalities.GetProfileResponse, error)
}

//go:generate mockgen -destination=./mocks/mock_ImageService.go -package=sign_up_mocks . ImageService
type ImageService interface {
	GetImageLinksByUserId(ctx context.Context, id int) ([]models.Image, error)
}

//easyjson:skip
type Handler struct {
	communicationsClient generatedCommunications.CommunicationsClient
	sessionClient        generatedAuth.AuthClient
	personalitiesClient  generatedPersonalities.PersonalitiesClient
	messageClient        generatedMessage.MessageClient
	imageService         ImageService
	logger               *zap.Logger
}

type Response struct {
	ID          int            `json:"id"`
	Username    string         `json:"username"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Images      []models.Image `json:"images"`
	LastMessage string         `json:"last_message"`
	Self        bool           `json:"self"`
	Time        string         `json:"time"`
}

type Responses struct {
	Responses []Response `json:"responses"`
}

func NewHandler(communicationsClient generatedCommunications.CommunicationsClient, sessionClient generatedAuth.AuthClient,
	personalitiesClient generatedPersonalities.PersonalitiesClient, imageService ImageService, messageClient generatedMessage.MessageClient, logger *zap.Logger) *Handler {
	return &Handler{
		communicationsClient: communicationsClient,
		sessionClient:        sessionClient,
		personalitiesClient:  personalitiesClient,
		imageService:         imageService,
		messageClient:        messageClient,
		logger:               logger,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req_id := ctx.Value(consts.RequestIDKey).(string)
	h.logger.Info("Handling request", zap.String("request_id", req_id))
	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		h.logger.Error("GetMatches Handler: bad getting cookie ", zap.Error(err))
		http.Error(w, "Вы не авторизованы", http.StatusUnauthorized)
		return
	}
	getUserIdRequest := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: cookie.Value}
	userId, err := h.sessionClient.GetUserIDBySessionID(ctx, getUserIdRequest)
	if err != nil {
		h.logger.Error("GetMatches Handler: bad getting user id ", zap.Error(err))
		http.Error(w, "Вы не авторизованы", http.StatusUnauthorized)
		return
	}
	getMatchListRequest := &generatedCommunications.GetMatchListRequest{UserID: userId.UserId}
	authors, err := h.communicationsClient.GetMatchList(ctx, getMatchListRequest)
	if err != nil {
		h.logger.Error("GetMatches Handler: bad getting authors ", zap.Error(err))
		http.Error(w, "Что-то пошло не так :(", http.StatusInternalServerError)
		return
	}

	var chats []Response
	for _, author := range authors.Authors {
		var chatter Response
		getProfileRequest := &generatedPersonalities.GetProfileRequest{Id: author}
		profile, err := h.personalitiesClient.GetProfile(ctx, getProfileRequest)
		if err != nil {
			h.logger.Error("GetMatches Handler: bad getting profile ", zap.Error(err))
			http.Error(w, "Что-то пошло не так :(", http.StatusInternalServerError)
			return
		}
		chatter.FirstName = profile.Profile.FirstName
		chatter.LastName = profile.Profile.LastName
		getUsernameRequest := &generatedPersonalities.GetUsernameByUserIDRequest{UserID: author}
		username, err := h.personalitiesClient.GetUsernameByUserID(ctx, getUsernameRequest)
		if err != nil {
			h.logger.Error("GetMatches Handler: bad getting username ", zap.Error(err))
			http.Error(w, "Что-то пошло не так :(", http.StatusInternalServerError)
			return
		}
		chatter.Username = username.Username
		var links []models.Image
		links, err = h.imageService.GetImageLinksByUserId(ctx, int(author))
		if err != nil {
			h.logger.Error("getimagelinkbyuserid error", zap.Error(err))
			http.Error(w, "Что-то пошло не так :(", http.StatusInternalServerError)
			return
		}
		chatter.Images = links
		chatter.ID = int(author)

		getLastRequest := &generatedMessage.GetLastMessageRequest{AuthorID: userId.UserId, ReceiverID: author}
		msg, err := h.messageClient.GetLastMessage(ctx, getLastRequest)
		if err != nil {
			h.logger.Error("getlastmessage error", zap.Error(err))
			http.Error(w, "Что-то пошло не так :(", http.StatusInternalServerError)
			return
		}
		if msg.Message == "" {
			getMatchRequest := &generatedCommunications.GetMatchTimeRequest{
				FirstUser:  userId.UserId,
				SecondUser: author,
			}
			time, err := h.communicationsClient.GetMatchTime(ctx, getMatchRequest)
			if err != nil {
				h.logger.Error("getmatchtime error", zap.Error(err))
				http.Error(w, "Что-то пошло не так :(", http.StatusInternalServerError)
				return
			}
			chatter.Time = time.Time
		} else {
			chatter.LastMessage = msg.Message
			chatter.Self = msg.Self
			chatter.Time = msg.Time
		}
		chats = append(chats, chatter)
	}

	sort.Slice(chats, func(i, j int) bool {
		a, err := time.Parse("RFC3339", chats[i].Time)
		if err != nil {
			return false
		}
		b, err := time.Parse("RFC3339", chats[j].Time)
		if err != nil {
			return false
		}
		return a.Before(b)
	})
	responses := Responses{Responses: chats}
	w.Header().Set("Content-Type", "application/json")

	jsonData, err := easyjson.Marshal(responses)
	if err != nil {
		h.logger.Error("json marshal error", zap.Error(err))
		http.Error(w, "json marshal error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonData)
	if err != nil {
		h.logger.Error("GetMatches Handler: error writing response", zap.Error(err))
		http.Error(w, "Что-то пошло не так :(", http.StatusInternalServerError)
		return
	}
	h.logger.Info("GetMatches Handler: success")

}

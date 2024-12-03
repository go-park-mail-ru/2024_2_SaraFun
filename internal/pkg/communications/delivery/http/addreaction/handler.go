package addreaction

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedCommunications "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
)

//go:generate mockgen -destination=./mocks/mock_ReactionService.go -package=sign_up_mocks . ReactionService

//go:generate mockgen -destination=./mocks/mock_SessionService.go -package=sign_up_mocks . SessionService

//go:generate easyjson -all handler.go

type SessionClient interface {
	GetUserIDBySessionID(ctx context.Context, in *generatedAuth.GetUserIDBySessionIDRequest) (*generatedAuth.GetUserIDBYSessionIDResponse, error)
}

type ReactionClient interface {
	AddReaction(ctx context.Context,
		in *generatedCommunications.AddReactionRequest) (*generatedCommunications.AddReactionResponse, error)
}

//easyjson:skip
type Handler struct {
	reactionClient generatedCommunications.CommunicationsClient
	SessionClient  generatedAuth.AuthClient
	logger         *zap.Logger
}

func NewHandler(reactionClient generatedCommunications.CommunicationsClient, sessionClient generatedAuth.AuthClient, logger *zap.Logger) *Handler {
	return &Handler{reactionClient: reactionClient, SessionClient: sessionClient, logger: logger}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req_id := ctx.Value(consts.RequestIDKey).(string)
	h.logger.Info("Handling request", zap.String("request_id", req_id))
	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		h.logger.Error("AddReaction Handler: bad getting cookie ", zap.Error(err))
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}
	getUserIdRequest := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: cookie.Value}
	userId, err := h.SessionClient.GetUserIDBySessionID(ctx, getUserIdRequest)
	if err != nil {
		h.logger.Error("AddReaction Handler: bad getting user id ", zap.Error(err))
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}
	var reaction models.Reaction
	//err = json.NewDecoder(r.Body).Decode(&reaction)
	err = easyjson.UnmarshalFromReader(r.Body, &reaction)
	if err != nil {
		h.logger.Error("AddReaction Handler: bad unmarshal", zap.Error(err))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	reaction.Author = int(userId.UserId)
	if err != nil {
		h.logger.Error("AddReaction Handler: bad decoding ", zap.Error(err))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	react := &generatedCommunications.Reaction{
		ID:       int32(reaction.Id),
		Author:   int32(reaction.Author),
		Receiver: int32(reaction.Receiver),
		Type:     reaction.Type,
	}
	addReactionRequest := &generatedCommunications.AddReactionRequest{Reaction: react}
	_, err = h.reactionClient.AddReaction(ctx, addReactionRequest)
	if err != nil {
		h.logger.Error("AddReaction Handler: error adding reaction", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	h.logger.Info("AddReaction Handler: added reaction", zap.Any("reaction", reaction))
	fmt.Fprintf(w, "ok")
}

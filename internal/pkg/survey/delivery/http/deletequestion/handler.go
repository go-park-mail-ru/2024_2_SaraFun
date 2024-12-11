package deletequestion

import (
	"fmt"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedSurvey "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	authClient   generatedAuth.AuthClient
	surveyClient generatedSurvey.SurveyClient
	logger       *zap.Logger
}

func NewHandler(authClient generatedAuth.AuthClient, surveyClient generatedSurvey.SurveyClient, logger *zap.Logger) *Handler {
	return &Handler{
		authClient:   authClient,
		surveyClient: surveyClient,
		logger:       logger,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cookie, err := r.Cookie(consts.SessionCookie)
	if err != nil {
		h.logger.Error("bad cookie", zap.Error(err))
		http.Error(w, "bad cookie", http.StatusUnauthorized)
		return
	}
	getUserIDReq := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: cookie.Value}
	_, err = h.authClient.GetUserIDBySessionID(ctx, getUserIDReq)
	if err != nil {
		h.logger.Error("get user id by session id", zap.Error(err))
		http.Error(w, "get user id by session id", http.StatusUnauthorized)
		return
	}
	//content := r.URL.Query().Get("content")
	vars := mux.Vars(r)
	content := vars["content"]
	h.logger.Info("delete question", zap.String("content", content))
	deleteQuestionReq := &generatedSurvey.DeleteQuestionRequest{
		Content: content,
	}
	h.logger.Info("delete question", zap.Any("deleteQuest", deleteQuestionReq))
	_, err = h.surveyClient.DeleteQuestion(ctx, deleteQuestionReq)
	if err != nil {
		h.logger.Error("delete question", zap.Error(err))
		http.Error(w, "delete question", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "ok")
}

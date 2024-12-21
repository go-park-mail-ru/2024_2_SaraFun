package getquestions

import (
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedSurvey "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
)

//go:generate easyjson -all handler.go

//easyjson:skip
type Handler struct {
	authClient   generatedAuth.AuthClient
	surveyClient generatedSurvey.SurveyClient
	logger       *zap.Logger
}

type Response struct {
	Questions []models.AdminQuestion `json:"questions"`
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

	getQuestionsRequest := &generatedSurvey.GetQuestionsRequest{}
	questions, err := h.surveyClient.GetQuestions(ctx, getQuestionsRequest)
	if err != nil {
		h.logger.Error("get questions", zap.Error(err))
		http.Error(w, "get questions", http.StatusUnauthorized)
		return
	}
	var respQuestions []models.AdminQuestion
	for _, question := range questions.Questions {
		respQuestions = append(respQuestions, models.AdminQuestion{
			Content: question.Content,
			Grade:   int(question.Grade),
		})
	}
	response := Response{Questions: respQuestions}
	jsonData, err := easyjson.Marshal(response)
	if err != nil {
		h.logger.Error("error marshaling survey stats", zap.Error(err))
		http.Error(w, "survey bad json marshaling", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		h.logger.Error("error writing response", zap.Error(err))
		http.Error(w, "survey bad response writing", http.StatusInternalServerError)
		return
	}
	h.logger.Info("get questions success")
}

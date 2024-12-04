package addquestion

import (
	"encoding/json"
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
	authCLient   generatedAuth.AuthClient
	surveyClient generatedSurvey.SurveyClient
	logger       *zap.Logger
}

type Response struct {
	ID int32 `json:"id"`
}

func NewHandler(authClient generatedAuth.AuthClient, surveyClient generatedSurvey.SurveyClient, logger *zap.Logger) *Handler {
	return &Handler{
		authCLient:   authClient,
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
	_, err = h.authCLient.GetUserIDBySessionID(ctx, getUserIDReq)
	if err != nil {
		h.logger.Error("get user id by session id", zap.Error(err))
		http.Error(w, "get user id by session id", http.StatusUnauthorized)
		return
	}

	var question models.AdminQuestion
	err = json.NewDecoder(r.Body).Decode(&question)
	if err != nil {
		h.logger.Error("decode question", zap.Error(err))
		http.Error(w, "json decode question error", http.StatusBadRequest)
		return
	}
	reqQuestion := &generatedSurvey.AdminQuestion{
		Content: question.Content,
		Grade:   int32(question.Grade),
	}
	addQuestionRequest := generatedSurvey.AddQuestionRequest{
		Question: reqQuestion,
	}
	questionID, err := h.surveyClient.AddQuestion(ctx, &addQuestionRequest)
	if err != nil {
		h.logger.Error("add question", zap.Error(err))
		http.Error(w, "add question error", http.StatusInternalServerError)
		return
	}
	response := Response{ID: questionID.QuestionID}
	jsonData, err := easyjson.Marshal(response)
	if err != nil {
		h.logger.Error("encode question", zap.Error(err))
		http.Error(w, "encode question error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		h.logger.Error("write response", zap.Error(err))
		http.Error(w, "write response error", http.StatusInternalServerError)
		return
	}
}

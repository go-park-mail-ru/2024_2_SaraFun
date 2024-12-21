package addsurvey

import (
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedSurvey "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	surveyClient  generatedSurvey.SurveyClient
	sessionClient generatedAuth.AuthClient
	logger        *zap.Logger
}

func NewHandler(surveyClient generatedSurvey.SurveyClient, sessionClient generatedAuth.AuthClient, logger *zap.Logger) *Handler {
	return &Handler{
		surveyClient:  surveyClient,
		sessionClient: sessionClient,
		logger:        logger,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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

	var survey models.Survey
	err = easyjson.UnmarshalFromReader(r.Body, &survey)
	if err != nil {
		h.logger.Error("error decoding survey", zap.Error(err))
		http.Error(w, "error decoding survey", http.StatusBadRequest)
		return
	}
	survey.Author = int(userId.UserId)

	reqSurvey := &generatedSurvey.SSurvey{
		Author:   int32(survey.Author),
		Question: survey.Question,
		Comment:  survey.Comment,
		Rating:   int32(survey.Rating),
		Grade:    int32(survey.Grade),
	}
	addSurveyRequest := &generatedSurvey.AddSurveyRequest{Survey: reqSurvey}
	_, err = h.surveyClient.AddSurvey(ctx, addSurveyRequest)
	if err != nil {
		h.logger.Error("error adding survey", zap.Error(err))
		http.Error(w, "error adding survey", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "ok")
}

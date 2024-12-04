package getsurveyinfo

import (
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedSurvey "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
)

//go:generate easyjson -all handler.go

type Response struct {
	Question      string
	AverageRating float32
	Grade         int
}

type Responses struct {
	Responses []Response
}

//easyjson:skip
type Handler struct {
	sessionClient generatedAuth.AuthClient
	surveyClient  generatedSurvey.SurveyClient
	logger        *zap.Logger
}

func NewHandler(authClient generatedAuth.AuthClient, surveyClient generatedSurvey.SurveyClient, logger *zap.Logger) *Handler {
	return &Handler{
		sessionClient: authClient,
		surveyClient:  surveyClient,
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

	getSurveyReq := &generatedSurvey.GetSurveyInfoRequest{}
	stats, err := h.surveyClient.GetSurveyInfo(ctx, getSurveyReq)
	if err != nil {
		h.logger.Error("error getting survey info", zap.Error(err))
		http.Error(w, "survey not found", http.StatusInternalServerError)
		return
	}
	var respStats []Response
	for _, st := range stats.Stats {
		h.logger.Info("stats", zap.Any("stats", st))
		stat := Response{
			Question:      st.Question,
			AverageRating: st.AvgRating,
			Grade:         int(st.Grade),
		}
		respStats = append(respStats, stat)
	}
	responses := Responses{Responses: respStats}
	jsonData, err := easyjson.Marshal(responses)
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
	h.logger.Info("get survey info success")
}

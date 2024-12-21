package getsurveyinfo

import (
	"context"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	authmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen/mocks"
	generatedSurvey "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/delivery/grpc/gen"
	surveymocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/delivery/grpc/gen/mocks"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	authClient := authmocks.NewMockAuthClient(mockCtrl)
	surveyClient := surveymocks.NewMockSurveyClient(mockCtrl)
	handler := NewHandler(authClient, surveyClient, logger)

	successStats := []*generatedSurvey.Stat{
		{
			Question:  "тестовый вопрос один",
			AvgRating: 3,
			Grade:     5,
		},
	}

	tests := []struct {
		name            string
		method          string
		path            string
		cookieValue     string
		stats           []*generatedSurvey.Stat
		authReturn      int
		authError       error
		authTimes       int
		surveyReturn    int
		surveyError     error
		surveyTimes     int
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:            "succesfull test",
			method:          "GET",
			path:            "/api/survey/getquestions",
			cookieValue:     "sparkit",
			stats:           successStats,
			authReturn:      1,
			authError:       nil,
			authTimes:       1,
			surveyReturn:    1,
			surveyError:     nil,
			surveyTimes:     1,
			expectedStatus:  http.StatusOK,
			expectedMessage: `{"Responses":[{"Question":"тестовый вопрос один","AverageRating":3,"Grade":5}]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getUserIDReq := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: tt.cookieValue}
			getUserIDResponse := &generatedAuth.GetUserIDBYSessionIDResponse{UserId: int32(tt.authReturn)}
			authClient.EXPECT().GetUserIDBySessionID(ctx, getUserIDReq).Return(getUserIDResponse, tt.authError).Times(tt.authTimes)

			getSurveyReq := &generatedSurvey.GetSurveyInfoRequest{}
			getSurveyResponse := &generatedSurvey.GetSurveyInfoResponse{Stats: successStats}
			surveyClient.EXPECT().GetSurveyInfo(ctx, getSurveyReq).Return(getSurveyResponse, tt.surveyError).Times(tt.surveyReturn)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			req = req.WithContext(ctx)
			cookie := &http.Cookie{
				Name:  consts.SessionCookie,
				Value: tt.cookieValue,
			}
			req.AddCookie(cookie)
			w := httptest.NewRecorder()
			handler.Handle(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", w.Code, tt.expectedStatus)
			}
			if w.Body.String() != tt.expectedMessage {
				t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), tt.expectedMessage)
			}
		})
	}
}

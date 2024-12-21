package addsurvey

import (
	"bytes"
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
	surveyClient := surveymocks.NewMockSurveyClient(mockCtrl)
	authClient := authmocks.NewMockAuthClient(mockCtrl)
	handler := NewHandler(surveyClient, authClient, logger)

	tests := []struct {
		name              string
		method            string
		path              string
		cookieValue       string
		authReturnID      int
		authReturnError   error
		authReturnCount   int
		surveyReturnID    int
		surveyReturnError error
		surveyReturnCount int
		question          string
		comment           string
		rating            int
		grade             int
		body              []byte
		expectedCode      int
		expectedMessage   string
	}{
		{
			name:        "good test",
			method:      "POST",
			path:        "/survey/sendsurvey",
			cookieValue: "sparkit",
			body: []byte(`{
    "question" : "тест",
    "comment": "норм",
    "rating": 0,
    "grade": 100
}`),
			question:          "тест",
			comment:           "норм",
			rating:            0,
			grade:             100,
			authReturnID:      1,
			authReturnError:   nil,
			authReturnCount:   1,
			surveyReturnID:    1,
			surveyReturnError: nil,
			surveyReturnCount: 1,
			expectedCode:      http.StatusOK,
			expectedMessage:   "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getUserReq := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: tt.cookieValue}
			getUserResponse := &generatedAuth.GetUserIDBYSessionIDResponse{UserId: int32(tt.authReturnID)}
			authClient.EXPECT().GetUserIDBySessionID(ctx, getUserReq).Return(getUserResponse, tt.authReturnError).
				Times(tt.authReturnCount)
			survey := &generatedSurvey.SSurvey{
				Author:   int32(tt.authReturnID),
				Question: tt.question,
				Comment:  tt.comment,
				Rating:   int32(tt.rating),
				Grade:    int32(tt.grade),
			}
			addSurveyReq := &generatedSurvey.AddSurveyRequest{
				Survey: survey,
			}
			addSurveyResponse := &generatedSurvey.AddSurveyResponse{
				SurveyID: int32(tt.surveyReturnID),
			}
			surveyClient.EXPECT().AddSurvey(ctx, addSurveyReq).Return(addSurveyResponse, tt.surveyReturnError).
				Times(tt.surveyReturnCount)

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(tt.body))
			req = req.WithContext(ctx)
			cookie := &http.Cookie{Name: consts.SessionCookie, Value: tt.cookieValue}
			req.AddCookie(cookie)
			w := httptest.NewRecorder()
			handler.Handle(w, req)
			if w.Code != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", w.Code, tt.expectedCode)
			}
			if w.Body.String() != tt.expectedMessage {
				t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), tt.expectedMessage)
			}

		})
	}
}

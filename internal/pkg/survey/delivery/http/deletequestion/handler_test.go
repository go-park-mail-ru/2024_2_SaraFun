package deletequestion

import (
	"context"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	authmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen/mocks"
	generatedSurvey "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/delivery/grpc/gen"
	surveymocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/delivery/grpc/gen/mocks"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
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

	tests := []struct {
		name            string
		method          string
		path            string
		cookieValue     string
		content         string
		authReturn      int
		authError       error
		authTimes       int
		surveyError     error
		surveyTimes     int
		expectedCode    int
		expectedMessage string
	}{
		{
			name:            "good test",
			method:          "DELETE",
			path:            "/api/survey/question/тест",
			cookieValue:     "sparkit",
			content:         "тест",
			authTimes:       1,
			surveyTimes:     1,
			expectedCode:    http.StatusOK,
			expectedMessage: "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authRequest := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: tt.cookieValue}
			authResponse := &generatedAuth.GetUserIDBYSessionIDResponse{UserId: int32(tt.authReturn)}
			authClient.EXPECT().GetUserIDBySessionID(gomock.Any(), authRequest).Return(authResponse, tt.authError)
			surveyRequest := &generatedSurvey.DeleteQuestionRequest{Content: tt.content}
			surveyResponse := &generatedSurvey.DeleteQuestionResponse{}
			surveyClient.EXPECT().DeleteQuestion(gomock.Any(), surveyRequest).Return(surveyResponse, tt.surveyError).Times(tt.surveyTimes)
			cookie := &http.Cookie{
				Name:  consts.SessionCookie,
				Value: tt.cookieValue,
			}
			t.Log(tt.path)
			req := httptest.NewRequest(tt.method, tt.path, nil)
			req = req.WithContext(ctx)
			req.AddCookie(cookie)
			req = mux.SetURLVars(req, map[string]string{"content": tt.content})
			vars := mux.Vars(req)
			t.Log(vars)
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

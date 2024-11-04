package getuserlist

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sparkit/internal/models"
	"testing"

	"github.com/golang/mock/gomock"
	getuserlist_mocks "sparkit/internal/handlers/getuserlist/mocks"
)

func TestGetUserListHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		name                string
		method              string
		path                string
		getUserListResponse []models.User
		getUserListError    error
		expectedStatus      int
		expectedResponse    string
		getUserListCalled   bool
	}{
		{
			name:   "successful user list retrieval",
			method: http.MethodGet,
			path:   "/userlist",
			getUserListResponse: []models.User{
				{ID: 0, Username: "user1"},
				{ID: 0, Username: "user2"},
			},
			expectedStatus:    http.StatusOK,
			expectedResponse:  `[{"id":0,"username":"user1"},{"id":0,"username":"user2"}]`,
			getUserListCalled: true,
		},
		{
			name:              "wrong method",
			method:            http.MethodPost,
			path:              "/userlist",
			expectedStatus:    http.StatusMethodNotAllowed,
			expectedResponse:  "Method not allowed\n",
			getUserListCalled: false,
		},
		{
			name:              "error fetching user list",
			method:            http.MethodGet,
			path:              "/userlist",
			getUserListError:  errors.New("database error"),
			expectedStatus:    http.StatusInternalServerError,
			expectedResponse:  "ошибка в получении списка пользователей\n",
			getUserListCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := getuserlist_mocks.NewMockUserUsecase(mockCtrl)
			handler := NewHandler(usecase)

			if tt.getUserListCalled {
				usecase.EXPECT().GetUserList(gomock.Any()).Return(tt.getUserListResponse, tt.getUserListError).Times(1)
			} else {
				usecase.EXPECT().GetUserList(gomock.Any()).Times(0)
			}

			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			handler.Handle(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", w.Code, tt.expectedStatus)
			}

			if w.Body.String() != tt.expectedResponse {
				t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), tt.expectedResponse)
			}
		})
	}
}

//nolint:golint
package signup

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"

	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	authmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen/mocks"
	generatedPayments "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc/gen"
	paymentsmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/payments/delivery/grpc/gen/mocks"
	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
	personalitiesmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen/mocks"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/mailru/easyjson"
)

//nolint:all
func TestHandler(t *testing.T) {
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	personalitiesClient := personalitiesmocks.NewMockPersonalitiesClient(mockCtrl)
	sessionClient := authmocks.NewMockAuthClient(mockCtrl)
	paymentsClient := paymentsmocks.NewMockPaymentClient(mockCtrl)

	handler := NewHandler(personalitiesClient, sessionClient, paymentsClient, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "test_req_id")

	validReq := &Request{
		Username:  "john_doe",
		Password:  "secret",
		FirstName: "John",
		LastName:  "Doe",
		Age:       30,
		BirthDate: "1990-01-01",
		Gender:    "male",
	}
	validBody, _ := easyjson.Marshal(validReq)

	tests := []struct {
		name                     string
		method                   string
		body                     []byte
		checkUsernameError       error
		usernameExists           bool
		createProfileError       error
		passwordIsBad            bool // Если true, хэширование пароля упадёт
		registerUserError        error
		createBalancesError      error
		createSessionError       error
		expectedStatus           int
		expectedResponseContains string
	}{
		{
			name:                     "not POST method",
			method:                   http.MethodGet,
			body:                     validBody,
			expectedStatus:           http.StatusMethodNotAllowed,
			expectedResponseContains: "Method not allowed",
		},
		{
			name:                     "bad json",
			method:                   http.MethodPost,
			body:                     []byte(`{bad json`),
			expectedStatus:           http.StatusBadRequest,
			expectedResponseContains: "Неверный формат данных",
		},
		{
			name:                     "check username error",
			method:                   http.MethodPost,
			body:                     validBody,
			checkUsernameError:       errors.New("check error"),
			expectedStatus:           http.StatusInternalServerError,
			expectedResponseContains: "Неудачная проверка на никнейм",
		},
		{
			name:                     "username exists",
			method:                   http.MethodPost,
			body:                     validBody,
			usernameExists:           true,
			expectedStatus:           http.StatusBadRequest,
			expectedResponseContains: "Пользователь с таким никнеймом уже существует",
		},
		{
			name:                     "create profile error",
			method:                   http.MethodPost,
			body:                     validBody,
			createProfileError:       errors.New("profile error"),
			expectedStatus:           http.StatusInternalServerError,
			expectedResponseContains: "profile error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/signup", bytes.NewBuffer(tt.body))
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			if tt.method == http.MethodPost && isEasyJSONValidRequest(tt.body) {

				var r Request
				_ = easyjson.Unmarshal(tt.body, &r)

				chReq := &generatedPersonalities.CheckUsernameExistsRequest{Username: r.Username}
				if tt.checkUsernameError == nil {
					ceResp := &generatedPersonalities.CheckUsernameExistsResponse{Exists: tt.usernameExists}
					personalitiesClient.EXPECT().CheckUsernameExists(gomock.Any(), chReq).
						Return(ceResp, nil).Times(1)
					if !tt.usernameExists {

						cpReq := &generatedPersonalities.CreateProfileRequest{Profile: &generatedPersonalities.Profile{
							FirstName: r.FirstName,
							LastName:  r.LastName,
							Age:       int32(r.Age),
							Gender:    r.Gender,
							BirthDate: r.BirthDate,
						}}
						if tt.createProfileError == nil {
							cpResp := &generatedPersonalities.CreateProfileResponse{ProfileId: 100}
							personalitiesClient.EXPECT().CreateProfile(gomock.Any(), cpReq).
								Return(cpResp, nil).Times(1)

							if tt.passwordIsBad && r.Password == "errorpass" {

							} else if !tt.passwordIsBad {

								rgReq := &generatedPersonalities.RegisterUserRequest{
									User: &generatedPersonalities.User{
										Username: r.Username,
										Password: "hashed_" + r.Password,
										Email:    "",
									},
								}
								if tt.registerUserError == nil {
									rgResp := &generatedPersonalities.RegisterUserResponse{UserId: 10}
									personalitiesClient.EXPECT().RegisterUser(gomock.Any(), rgReq).
										Return(rgResp, nil).Times(1)

									// CreateBalances
									cbReq := &generatedPayments.CreateBalancesRequest{
										UserID:          10,
										MoneyAmount:     0,
										DailyAmount:     consts.DailyLikeLimit,
										PurchasedAmount: 5,
									}
									if tt.createBalancesError == nil {
										cbResp := &generatedPayments.CreateBalancesResponse{}
										paymentsClient.EXPECT().CreateBalances(gomock.Any(), cbReq).
											Return(cbResp, nil).Times(1)

										csReq := &generatedAuth.CreateSessionRequest{
											User: &generatedAuth.User{
												ID:       10,
												Username: r.Username,
												Password: "hashed_" + r.Password,
												Email:    "",
												Profile:  100,
											},
										}
										if tt.createSessionError == nil {
											csResp := &generatedAuth.CreateSessionResponse{
												Session: &generatedAuth.Session{SessionID: "valid_session_id"},
											}
											sessionClient.EXPECT().CreateSession(gomock.Any(), csReq).
												Return(csResp, nil).Times(1)
										} else {
											sessionClient.EXPECT().CreateSession(gomock.Any(), csReq).
												Return(nil, tt.createSessionError).Times(1)
										}

									} else {
										paymentsClient.EXPECT().CreateBalances(gomock.Any(), cbReq).
											Return(nil, tt.createBalancesError).Times(1)
									}

								} else {
									personalitiesClient.EXPECT().RegisterUser(gomock.Any(), rgReq).
										Return(nil, tt.registerUserError).Times(1)
								}

							}

						} else {
							personalitiesClient.EXPECT().CreateProfile(gomock.Any(), cpReq).
								Return(nil, tt.createProfileError).Times(1)
						}

					}
				} else {
					personalitiesClient.EXPECT().CheckUsernameExists(gomock.Any(), chReq).
						Return(nil, tt.checkUsernameError).Times(1)
				}
			}

			handler.Handle(w, req)
			if w.Code != tt.expectedStatus {
				t.Errorf("%s: wrong status code: got %v, want %v", tt.name, w.Code, tt.expectedStatus)
			}
			if tt.expectedResponseContains != "" && !contains(w.Body.String(), tt.expectedResponseContains) {
				t.Errorf("%s: wrong body: got %v, want substring %v", tt.name, w.Body.String(), tt.expectedResponseContains)
			}

			if tt.expectedStatus == http.StatusOK && tt.createSessionError == nil &&
				!tt.passwordIsBad && tt.registerUserError == nil &&
				tt.createBalancesError == nil && tt.method == http.MethodPost &&
				isEasyJSONValidRequest(tt.body) && tt.checkUsernameError == nil && !tt.usernameExists && tt.createProfileError == nil {
				resp := w.Result()
				defer resp.Body.Close()
				cookies := resp.Cookies()
				found := false
				for _, c := range cookies {
					if c.Name == consts.SessionCookie && c.Value == "valid_session_id" {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("%s: expected session cookie to be set", tt.name)
				}
			}
		})
	}
}

func isEasyJSONValidRequest(b []byte) bool {
	var r Request
	return easyjson.Unmarshal(b, &r) == nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(substr) == 0 ||
			(len(s) > 0 && len(substr) > 0 && s[0:len(substr)] == substr) ||
			(len(s) > len(substr) && s[len(s)-len(substr):] == substr) ||
			(len(substr) > 0 && len(s) > len(substr) && findInString(s, substr)))
}

func findInString(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

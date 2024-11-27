package getmatches

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	generatedCommunications "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen"
	getmatches_mocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/http/getmatches/mocks"
	generatedPersonalities "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=./mocks/mock_CommunicationsClient.go -package=getmatches_mocks github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen CommunicationsClient
//go:generate mockgen -destination=./mocks/mock_AuthClient.go -package=getmatches_mocks github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen AuthClient
//go:generate mockgen -destination=./mocks/mock_PersonalitiesClient.go -package=getmatches_mocks github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/personalities/delivery/grpc/gen PersonalitiesClient
//go:generate mockgen -destination=./mocks/mock_ImageService.go -package=getmatches_mocks github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/http/getmatches ImageService

func TestHandler_Handle(t *testing.T) {
	tests := []struct {
		name                     string
		setupRequest             func() *http.Request
		mockAuthClient           func(mock *getmatches_mocks.MockAuthClient)
		mockCommunicationsClient func(mock *getmatches_mocks.MockCommunicationsClient)
		mockPersonalitiesClient  func(mock *getmatches_mocks.MockPersonalitiesClient)
		mockImageService         func(mock *getmatches_mocks.MockImageService)
		expectedStatus           int
		expectedResponseContains string
	}{
		{
			name: "Successful Response",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.AddCookie(&http.Cookie{
					Name:  consts.SessionCookie,
					Value: "valid-session-id",
				})
				ctx := context.WithValue(req.Context(), consts.RequestIDKey, "test-request-id")
				return req.WithContext(ctx)
			},
			mockAuthClient: func(mock *getmatches_mocks.MockAuthClient) {
				mock.EXPECT().
					GetUserIDBySessionID(gomock.Any(), &generatedAuth.GetUserIDBySessionIDRequest{
						SessionID: "valid-session-id",
					}).
					Return(&generatedAuth.GetUserIDBYSessionIDResponse{
						UserId: 1,
					}, nil)
			},
			mockCommunicationsClient: func(mock *getmatches_mocks.MockCommunicationsClient) {
				mock.EXPECT().
					GetMatchList(gomock.Any(), &generatedCommunications.GetMatchListRequest{
						UserID: 1,
					}).
					Return(&generatedCommunications.GetMatchListResponse{
						Authors: []int32{2, 3},
					}, nil)
			},
			mockPersonalitiesClient: func(mock *getmatches_mocks.MockPersonalitiesClient) {
				// Mock GetProfile for author 2
				mock.EXPECT().
					GetProfile(gomock.Any(), &generatedPersonalities.GetProfileRequest{
						Id: 2,
					}).
					Return(&generatedPersonalities.GetProfileResponse{
						Profile: &generatedPersonalities.Profile{
							ID:        2,
							FirstName: "John",
							LastName:  "Doe",
							Age:       30,
							Gender:    "Male",
							Target:    "Friendship",
							About:     "About John",
						},
					}, nil)
				// Mock GetUsernameByUserID for author 2
				mock.EXPECT().
					GetUsernameByUserID(gomock.Any(), &generatedPersonalities.GetUsernameByUserIDRequest{
						UserID: 2,
					}).
					Return(&generatedPersonalities.GetUsernameByUserIDResponse{
						Username: "johndoe",
					}, nil)
				// Mock GetProfile for author 3
				mock.EXPECT().
					GetProfile(gomock.Any(), &generatedPersonalities.GetProfileRequest{
						Id: 3,
					}).
					Return(&generatedPersonalities.GetProfileResponse{
						Profile: &generatedPersonalities.Profile{
							ID:        3,
							FirstName: "Jane",
							LastName:  "Smith",
							Age:       28,
							Gender:    "Female",
							Target:    "Dating",
							About:     "About Jane",
						},
					}, nil)
				// Mock GetUsernameByUserID for author 3
				mock.EXPECT().
					GetUsernameByUserID(gomock.Any(), &generatedPersonalities.GetUsernameByUserIDRequest{
						UserID: 3,
					}).
					Return(&generatedPersonalities.GetUsernameByUserIDResponse{
						Username: "janesmith",
					}, nil)
			},
			mockImageService: func(mock *getmatches_mocks.MockImageService) {
				// Mock GetImageLinksByUserId for author 2
				mock.EXPECT().
					GetImageLinksByUserId(gomock.Any(), 2).
					Return([]models.Image{
						{Id: 1, Link: "http://example.com/image1.jpg"},
					}, nil)
				// Mock GetImageLinksByUserId for author 3
				mock.EXPECT().
					GetImageLinksByUserId(gomock.Any(), 3).
					Return([]models.Image{
						{Id: 2, Link: "http://example.com/image2.jpg"},
					}, nil)
			},
			expectedStatus:           http.StatusOK,
			expectedResponseContains: `"username":"johndoe"`,
		},
		{
			name: "Missing Session Cookie",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				ctx := context.WithValue(req.Context(), consts.RequestIDKey, "test-request-id")
				return req.WithContext(ctx)
			},
			mockAuthClient:           func(mock *getmatches_mocks.MockAuthClient) {},
			mockCommunicationsClient: func(mock *getmatches_mocks.MockCommunicationsClient) {},
			mockPersonalitiesClient:  func(mock *getmatches_mocks.MockPersonalitiesClient) {},
			mockImageService:         func(mock *getmatches_mocks.MockImageService) {},
			expectedStatus:           http.StatusUnauthorized,
			expectedResponseContains: "session not found",
		},
		{
			name: "Failed to Get User ID",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.AddCookie(&http.Cookie{
					Name:  consts.SessionCookie,
					Value: "invalid-session-id",
				})
				ctx := context.WithValue(req.Context(), consts.RequestIDKey, "test-request-id")
				return req.WithContext(ctx)
			},
			mockAuthClient: func(mock *getmatches_mocks.MockAuthClient) {
				mock.EXPECT().
					GetUserIDBySessionID(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("session not found"))
			},
			mockCommunicationsClient: func(mock *getmatches_mocks.MockCommunicationsClient) {},
			mockPersonalitiesClient:  func(mock *getmatches_mocks.MockPersonalitiesClient) {},
			mockImageService:         func(mock *getmatches_mocks.MockImageService) {},
			expectedStatus:           http.StatusUnauthorized,
			expectedResponseContains: "session not found",
		},
		{
			name: "Failed to Get Match List",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.AddCookie(&http.Cookie{
					Name:  consts.SessionCookie,
					Value: "valid-session-id",
				})
				ctx := context.WithValue(req.Context(), consts.RequestIDKey, "test-request-id")
				return req.WithContext(ctx)
			},
			mockAuthClient: func(mock *getmatches_mocks.MockAuthClient) {
				mock.EXPECT().
					GetUserIDBySessionID(gomock.Any(), gomock.Any()).
					Return(&generatedAuth.GetUserIDBYSessionIDResponse{
						UserId: 1,
					}, nil)
			},
			mockCommunicationsClient: func(mock *getmatches_mocks.MockCommunicationsClient) {
				mock.EXPECT().
					GetMatchList(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("failed to get matches"))
			},
			mockPersonalitiesClient:  func(mock *getmatches_mocks.MockPersonalitiesClient) {},
			mockImageService:         func(mock *getmatches_mocks.MockImageService) {},
			expectedStatus:           http.StatusUnauthorized,
			expectedResponseContains: "session not found",
		},
		{
			name: "Failed to Get Profile",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.AddCookie(&http.Cookie{
					Name:  consts.SessionCookie,
					Value: "valid-session-id",
				})
				ctx := context.WithValue(req.Context(), consts.RequestIDKey, "test-request-id")
				return req.WithContext(ctx)
			},
			mockAuthClient: func(mock *getmatches_mocks.MockAuthClient) {
				mock.EXPECT().
					GetUserIDBySessionID(gomock.Any(), gomock.Any()).
					Return(&generatedAuth.GetUserIDBYSessionIDResponse{
						UserId: 1,
					}, nil)
			},
			mockCommunicationsClient: func(mock *getmatches_mocks.MockCommunicationsClient) {
				mock.EXPECT().
					GetMatchList(gomock.Any(), gomock.Any()).
					Return(&generatedCommunications.GetMatchListResponse{
						Authors: []int32{2},
					}, nil)
			},
			mockPersonalitiesClient: func(mock *getmatches_mocks.MockPersonalitiesClient) {
				mock.EXPECT().
					GetProfile(gomock.Any(), gomock.Any()).
					Return(&generatedPersonalities.GetProfileResponse{
						Profile: &generatedPersonalities.Profile{
							ID:        0,
							FirstName: "",
							LastName:  "",
							Age:       0,
							Gender:    "",
							Target:    "",
							About:     "",
						},
					}, errors.New("failed to get profile"))
			},
			mockImageService:         func(mock *getmatches_mocks.MockImageService) {},
			expectedStatus:           http.StatusInternalServerError,
			expectedResponseContains: "bad get profile",
		},

		{
			name: "Failed to Get Image Links",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.AddCookie(&http.Cookie{
					Name:  consts.SessionCookie,
					Value: "valid-session-id",
				})
				ctx := context.WithValue(req.Context(), consts.RequestIDKey, "test-request-id")
				return req.WithContext(ctx)
			},
			mockAuthClient: func(mock *getmatches_mocks.MockAuthClient) {
				mock.EXPECT().
					GetUserIDBySessionID(gomock.Any(), gomock.Any()).
					Return(&generatedAuth.GetUserIDBYSessionIDResponse{
						UserId: 1,
					}, nil)
			},
			mockCommunicationsClient: func(mock *getmatches_mocks.MockCommunicationsClient) {
				mock.EXPECT().
					GetMatchList(gomock.Any(), gomock.Any()).
					Return(&generatedCommunications.GetMatchListResponse{
						Authors: []int32{2},
					}, nil)
			},
			mockPersonalitiesClient: func(mock *getmatches_mocks.MockPersonalitiesClient) {
				mock.EXPECT().
					GetProfile(gomock.Any(), gomock.Any()).
					Return(&generatedPersonalities.GetProfileResponse{
						Profile: &generatedPersonalities.Profile{
							ID:        2,
							FirstName: "John",
							LastName:  "Doe",
							Age:       30,
							Gender:    "Male",
							Target:    "Friendship",
							About:     "About John",
						},
					}, nil)
			},
			mockImageService: func(mock *getmatches_mocks.MockImageService) {
				mock.EXPECT().
					GetImageLinksByUserId(gomock.Any(), 2).
					Return(nil, errors.New("failed to get images"))
			},
			expectedStatus:           http.StatusInternalServerError,
			expectedResponseContains: "failed to get images",
		},
		{
			name: "Failed to Get Username",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.AddCookie(&http.Cookie{
					Name:  consts.SessionCookie,
					Value: "valid-session-id",
				})
				ctx := context.WithValue(req.Context(), consts.RequestIDKey, "test-request-id")
				return req.WithContext(ctx)
			},
			mockAuthClient: func(mock *getmatches_mocks.MockAuthClient) {
				mock.EXPECT().
					GetUserIDBySessionID(gomock.Any(), gomock.Any()).
					Return(&generatedAuth.GetUserIDBYSessionIDResponse{
						UserId: 1,
					}, nil)
			},
			mockCommunicationsClient: func(mock *getmatches_mocks.MockCommunicationsClient) {
				mock.EXPECT().
					GetMatchList(gomock.Any(), gomock.Any()).
					Return(&generatedCommunications.GetMatchListResponse{
						Authors: []int32{2},
					}, nil)
			},
			mockPersonalitiesClient: func(mock *getmatches_mocks.MockPersonalitiesClient) {
				mock.EXPECT().
					GetProfile(gomock.Any(), gomock.Any()).
					Return(&generatedPersonalities.GetProfileResponse{
						Profile: &generatedPersonalities.Profile{
							ID:        2,
							FirstName: "John",
							LastName:  "Doe",
							Age:       30,
							Gender:    "Male",
							Target:    "Friendship",
							About:     "About John",
						},
					}, nil)
				mock.EXPECT().
					GetUsernameByUserID(gomock.Any(), gomock.Any()).
					Return(&generatedPersonalities.GetUsernameByUserIDResponse{
						Username: "",
					}, errors.New("failed to get username"))
			},
			mockImageService: func(mock *getmatches_mocks.MockImageService) {
				mock.EXPECT().
					GetImageLinksByUserId(gomock.Any(), 2).
					Return([]models.Image{
						{Id: 1, Link: "http://example.com/image1.jpg"},
					}, nil)
			},
			expectedStatus:           http.StatusInternalServerError,
			expectedResponseContains: "bad get username",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create gomock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create mock clients and services
			mockAuthClient := getmatches_mocks.NewMockAuthClient(ctrl)
			mockCommunicationsClient := getmatches_mocks.NewMockCommunicationsClient(ctrl)
			mockPersonalitiesClient := getmatches_mocks.NewMockPersonalitiesClient(ctrl)
			mockImageService := getmatches_mocks.NewMockImageService(ctrl)

			// Setup mocks
			tt.mockAuthClient(mockAuthClient)
			tt.mockCommunicationsClient(mockCommunicationsClient)
			tt.mockPersonalitiesClient(mockPersonalitiesClient)
			tt.mockImageService(mockImageService)

			// Create logger
			logger := zap.NewNop()

			// Create handler
			handler := NewHandler(mockCommunicationsClient, mockAuthClient, mockPersonalitiesClient, mockImageService, logger)

			// Create request
			req := tt.setupRequest()

			// Create ResponseRecorder
			rr := httptest.NewRecorder()

			// Call handler
			handler.Handle(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("expected status code %d, got %d", tt.expectedStatus, status)
			}

			// Check response body
			if !strings.Contains(rr.Body.String(), tt.expectedResponseContains) {
				t.Errorf("expected response body to contain %q, got %q", tt.expectedResponseContains, rr.Body.String())
			}
		})
	}
}

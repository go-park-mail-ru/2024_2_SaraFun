package uploadimage

//
//import (
//	"bytes"
//	"context"
//	"errors"
//	"io"
//	"mime/multipart"
//	"net/http"
//	"net/http/httptest"
//	"path/filepath"
//	"strconv"
//	"testing"
//	"time"
//
//	"github.com/golang/mock/gomock"
//	"go.uber.org/zap"
//
//	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
//	authmocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen/mocks"
//	uploadimage_mocks "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/image/delivery/uploadimage/mocks"
//	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
//)
//
////nolint:all
//func TestHandler(t *testing.T) {
//	logger := zap.NewNop()
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//
//	imageService := uploadimage_mocks.NewMockImageService(mockCtrl)
//	sessionService := authmocks.NewMockAuthClient(mockCtrl)
//
//	handler := NewHandler(imageService, sessionService, logger)
//
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//	ctx = context.WithValue(ctx, consts.RequestIDKey, "test_req_id")
//
//	buildMultipartRequest := func(method, filename, fieldName, fileContent string, number string, cookieValue string) (*http.Request, error) {
//		var body bytes.Buffer
//		writer := multipart.NewWriter(&body)
//		if filename != "" && fieldName != "" {
//			part, err := writer.CreateFormFile(fieldName, filename)
//			if err != nil {
//				return nil, err
//			}
//			_, err = io.WriteString(part, fileContent)
//			if err != nil {
//				return nil, err
//			}
//		}
//		if number != "" {
//			err := writer.WriteField("number", number)
//			if err != nil {
//				return nil, err
//			}
//		}
//		writer.Close()
//
//		req := httptest.NewRequest(method, "/uploadimage", &body)
//		req.Header.Set("Content-Type", writer.FormDataContentType())
//		req = req.WithContext(ctx)
//		if cookieValue != "" {
//			req.AddCookie(&http.Cookie{
//				Name:  consts.SessionCookie,
//				Value: cookieValue,
//			})
//		}
//		return req, nil
//	}
//
//	tests := []struct {
//		name                     string
//		method                   string
//		filename                 string
//		fieldName                string
//		fileContent              string
//		number                   string
//		cookieValue              string
//		getUserError             error
//		saveImageError           error
//		expectedStatus           int
//		expectedResponseContains string
//	}{
//		{
//			name:                     "parse multipart form error (no multipart)",
//			method:                   http.MethodPost,
//			expectedStatus:           http.StatusBadRequest,
//			expectedResponseContains: "bad image",
//		},
//		{
//			name:                     "bad image header",
//			method:                   http.MethodPost,
//			filename:                 "test.txt",
//			fieldName:                "image",
//			fileContent:              "",
//			number:                   "1",
//			cookieValue:              "valid_session",
//			expectedStatus:           http.StatusBadRequest,
//			expectedResponseContains: "bad image header",
//		},
//		{
//			name:                     "invalid image format",
//			method:                   http.MethodPost,
//			filename:                 "image.txt",
//			fieldName:                "image",
//			fileContent:              "not an image content",
//			number:                   "1",
//			cookieValue:              "valid_session",
//			expectedStatus:           http.StatusBadRequest,
//			expectedResponseContains: "invalid image format",
//		},
//		{
//			name:                     "no cookie",
//			method:                   http.MethodPost,
//			filename:                 "image.jpg",
//			fieldName:                "image",
//			fileContent:              "\xFF\xD8\xFF\xE0", // JPEG
//			number:                   "1",
//			expectedStatus:           http.StatusUnauthorized,
//			expectedResponseContains: "session not found",
//		},
//		{
//			name:                     "save image error",
//			method:                   http.MethodPost,
//			filename:                 "image.jpg",
//			fieldName:                "image",
//			fileContent:              "\xFF\xD8\xFF\xE0", // JPEG
//			number:                   "1",
//			cookieValue:              "valid_session",
//			saveImageError:           errors.New("save error"),
//			expectedStatus:           http.StatusInternalServerError,
//			expectedResponseContains: "save image err",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			var req *http.Request
//			var err error
//			if tt.method == http.MethodPost {
//				if tt.filename != "" && tt.fieldName != "" {
//					req, err = buildMultipartRequest(tt.method, tt.filename, tt.fieldName, tt.fileContent, tt.number, tt.cookieValue)
//					if err != nil {
//						t.Fatalf("failed to build request: %v", err)
//					}
//				} else if tt.filename == "" && tt.fieldName == "" {
//					req = httptest.NewRequest(tt.method, "/uploadimage", nil)
//					req = req.WithContext(ctx)
//					if tt.cookieValue != "" {
//						req.AddCookie(&http.Cookie{Name: consts.SessionCookie, Value: tt.cookieValue})
//					}
//				} else {
//					req = httptest.NewRequest(tt.method, "/uploadimage", nil)
//					req = req.WithContext(ctx)
//				}
//			} else {
//				req = httptest.NewRequest(tt.method, "/uploadimage", nil)
//				req = req.WithContext(ctx)
//			}
//
//			w := httptest.NewRecorder()
//
//			if tt.method != http.MethodPost {
//				handler.Handle(w, req)
//				checkResponse(t, w, tt.expectedStatus, tt.expectedResponseContains)
//				return
//			}
//
//			if tt.filename == "" && tt.fieldName == "" {
//				// Нет multipart
//				handler.Handle(w, req)
//				checkResponse(t, w, tt.expectedStatus, tt.expectedResponseContains)
//				return
//			}
//
//			if tt.fieldName != "image" {
//				handler.Handle(w, req)
//				checkResponse(t, w, tt.expectedStatus, tt.expectedResponseContains)
//				return
//			}
//
//			if tt.filename == "test.txt" && tt.fileContent == "" {
//				handler.Handle(w, req)
//				checkResponse(t, w, tt.expectedStatus, tt.expectedResponseContains)
//				return
//			}
//
//			if tt.filename == "image.txt" {
//				handler.Handle(w, req)
//				checkResponse(t, w, tt.expectedStatus, tt.expectedResponseContains)
//				return
//			}
//
//			if tt.cookieValue == "" {
//				handler.Handle(w, req)
//				checkResponse(t, w, tt.expectedStatus, tt.expectedResponseContains)
//				return
//			}
//
//			getUserReq := &generatedAuth.GetUserIDBySessionIDRequest{SessionID: tt.cookieValue}
//			if tt.getUserError == nil {
//				userResp := &generatedAuth.GetUserIDBYSessionIDResponse{UserId: 10}
//				sessionService.EXPECT().GetUserIDBySessionID(gomock.Any(), getUserReq).
//					Return(userResp, nil).Times(1)
//			} else {
//				sessionService.EXPECT().GetUserIDBySessionID(gomock.Any(), getUserReq).
//					Return(nil, tt.getUserError).Times(1)
//				handler.Handle(w, req)
//				checkResponse(t, w, tt.expectedStatus, tt.expectedResponseContains)
//				return
//			}
//
//			if _, err2 := strconv.Atoi(tt.number); err2 != nil {
//				handler.Handle(w, req)
//				checkResponse(t, w, tt.expectedStatus, tt.expectedResponseContains)
//				return
//			}
//
//			fileExt := filepath.Ext(tt.filename)
//			if tt.saveImageError == nil {
//				imageService.EXPECT().SaveImage(gomock.Any(), gomock.Any(), fileExt, 10, gomock.Any()).
//					Return(100, nil).Times(1)
//			} else {
//				imageService.EXPECT().SaveImage(gomock.Any(), gomock.Any(), fileExt, 10, gomock.Any()).
//					Return(0, tt.saveImageError).Times(1)
//			}
//
//			handler.Handle(w, req)
//			checkResponse(t, w, tt.expectedStatus, tt.expectedResponseContains)
//		})
//	}
//}
//
//func checkResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedContains string) {
//	if w.Code != expectedStatus {
//		t.Errorf("wrong status code: got %v, want %v", w.Code, expectedStatus)
//	}
//	if expectedContains != "" && !contains(w.Body.String(), expectedContains) {
//		t.Errorf("wrong body: got %v, want substring %v", w.Body.String(), expectedContains)
//	}
//}
//
//func contains(s, substr string) bool {
//	return len(s) >= len(substr) &&
//		(s == substr ||
//			len(substr) == 0 ||
//			(len(s) > 0 && len(substr) > 0 && string(s[0:len(substr)]) == substr) ||
//			(len(s) > len(substr) && string(s[len(s)-len(substr):]) == substr) ||
//			(len(substr) > 0 && len(s) > len(substr) && findInString(s, substr)))
//}
//
//func findInString(s, substr string) bool {
//	for i := 0; i+len(substr) <= len(s); i++ {
//		if s[i:i+len(substr)] == substr {
//			return true
//		}
//	}
//	return false
//}

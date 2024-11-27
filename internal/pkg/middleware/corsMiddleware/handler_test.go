package corsMiddleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/stretchr/testify/require"
)

func nextHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Next handler called"))
}

type testCase struct {
	name           string
	method         string
	url            string
	origin         string
	expectedStatus int
	expectedBody   string
	expectLog      bool
	expectNext     bool
}

func TestMiddleware(t *testing.T) {
	tests := []testCase{
		{
			name:           "GET request with Origin",
			method:         "GET",
			url:            "/test",
			origin:         "http://example.com",
			expectedStatus: http.StatusOK,
			expectedBody:   "Next handler called",
			expectLog:      true,
			expectNext:     true,
		},
		{
			name:           "OPTIONS request with Origin",
			method:         "OPTIONS",
			url:            "/preflight",
			origin:         "https://preflight.com",
			expectedStatus: http.StatusOK,
			expectedBody:   "",
			expectLog:      true,
			expectNext:     false,
		},
		{
			name:           "GET request without Origin",
			method:         "GET",
			url:            "/no-origin",
			origin:         "",
			expectedStatus: http.StatusOK,
			expectedBody:   "Next handler called",
			expectLog:      true,
			expectNext:     true,
		},
		{
			name:           "OPTIONS request without Origin",
			method:         "OPTIONS",
			url:            "/preflight-no-origin",
			origin:         "",
			expectedStatus: http.StatusOK,
			expectedBody:   "",
			expectLog:      true,
			expectNext:     false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			core, logs := observer.New(zapcore.InfoLevel)
			logger := zap.New(core)

			mw := New(logger)
			handler := mw.Middleware(http.HandlerFunc(nextHandler))

			req := httptest.NewRequest(tt.method, tt.url, nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code, "Статус код должен совпадать")
			require.Equal(t, tt.expectedBody, rr.Body.String(), "Тело ответа должно совпадать")

			allowMethods := rr.Header().Get("Access-Control-Allow-Methods")
			require.Equal(t, "POST,PUT,DELETE,GET", allowMethods, "Access-Control-Allow-Methods не совпадает")

			allowHeaders := rr.Header().Get("Access-Control-Allow-Headers")
			require.Equal(t, "Content-Type", allowHeaders, "Access-Control-Allow-Headers не совпадает")

			allowCredentials := rr.Header().Get("Access-Control-Allow-Credentials")
			require.Equal(t, "true", allowCredentials, "Access-Control-Allow-Credentials не совпадает")

			allowOrigin := rr.Header().Get("Access-Control-Allow-Origin")
			require.Equal(t, tt.origin, allowOrigin, "Access-Control-Allow-Origin не совпадает")

			if tt.expectLog {
				logEntries := logs.All()
				if tt.method == http.MethodOptions {
					require.Len(t, logEntries, 1, "Ожидается ровно одна запись в логах для preflight-запроса")
					require.Equal(t, "Handling request preflight", logEntries[0].Message, "Сообщение лога не совпадает для preflight-запроса")
					require.Equal(t, tt.method, logEntries[0].ContextMap()["method"], "HTTP метод в логе не совпадает")
					require.Equal(t, tt.url, logEntries[0].ContextMap()["url"], "URL в логе не совпадает")
				} else {
					require.Len(t, logEntries, 1, "Ожидается ровно одна запись в логах для обычного запроса")
					require.Equal(t, "Handling request", logEntries[0].Message, "Сообщение лога не совпадает для обычного запроса")
					require.Equal(t, tt.method, logEntries[0].ContextMap()["method"], "HTTP метод в логе не совпадает")
					require.Equal(t, tt.url, logEntries[0].ContextMap()["url"], "URL в логе не совпадает")
				}
			} else {
				logEntries := logs.All()
				require.Len(t, logEntries, 0, "Не ожидались записи в логах")
			}

			if tt.expectNext {
				require.Equal(t, tt.expectedBody, rr.Body.String(), "Тело ответа должно совпадать с ответом nextHandler")
			} else {
				require.Empty(t, rr.Body.String(), "Тело ответа должно быть пустым для preflight-запросов")
			}
		})
	}
}

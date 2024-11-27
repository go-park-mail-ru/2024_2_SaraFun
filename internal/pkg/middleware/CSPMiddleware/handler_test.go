// internal/pkg/middleware/CSPMiddleware/middleware_test.go
package CSPMiddleware

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
	expectedStatus int
	expectedBody   string
	expectedCSP    string
	expectLog      bool
}

func TestMiddleware(t *testing.T) {
	tests := []testCase{
		{
			name:           "GET request",
			method:         "GET",
			url:            "/test",
			expectedStatus: http.StatusOK,
			expectedBody:   "Next handler called",
			expectedCSP:    "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self';base-uri 'self';form-action 'self'",
			expectLog:      true,
		},
		{
			name:           "POST request",
			method:         "POST",
			url:            "/submit",
			expectedStatus: http.StatusOK,
			expectedBody:   "Next handler called",
			expectedCSP:    "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self';base-uri 'self';form-action 'self'",
			expectLog:      true,
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

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedStatus, rr.Code, "Статус код должен совпадать")

			require.Equal(t, tt.expectedBody, rr.Body.String(), "Тело ответа должно совпадать")

			csp := rr.Header().Get("Content-Security-Policy")
			require.Equal(t, tt.expectedCSP, csp, "Заголовок Content-Security-Policy не совпадает")

			if tt.expectLog {
				logEntries := logs.All()
				require.Len(t, logEntries, 1, "Ожидается ровно одна запись в логах")

				entry := logEntries[0]
				require.Equal(t, "Handling request", entry.Message, "Сообщение лога не совпадает")
				require.Equal(t, tt.method, entry.ContextMap()["method"], "HTTP метод в логе не совпадает")
				require.Equal(t, tt.url, entry.ContextMap()["url"], "URL в логе не совпадает")
			} else {
				logEntries := logs.All()
				require.Len(t, logEntries, 0, "Не ожидались записи в логах")
			}
		})
	}
}

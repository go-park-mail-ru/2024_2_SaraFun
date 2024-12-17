package CSPMiddleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/middleware/CSPMiddleware"
	"go.uber.org/zap"
)

func TestCSPMiddleware(t *testing.T) {
	logger := zap.NewNop()
	mw := CSPMiddleware.New(logger)

	nextCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusTeapot) // для проверки что мы дошли до хендлера
	})

	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	rr := httptest.NewRecorder()

	mw.Middleware(nextHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusTeapot {
		t.Errorf("expected status code from next handler to be 418, got %d", rr.Code)
	}

	if !nextCalled {
		t.Errorf("expected next handler to be called")
	}

	csp := rr.Header().Get("Content-Security-Policy")
	expectedCSP := "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self';base-uri 'self';form-action 'self'"
	if csp != expectedCSP {
		t.Errorf("unexpected Content-Security-Policy header: got %q, want %q", csp, expectedCSP)
	}
}

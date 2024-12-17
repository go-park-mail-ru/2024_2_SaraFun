package corsMiddleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/middleware/corsMiddleware"
	"go.uber.org/zap"
)

func TestMiddleware(t *testing.T) {
	logger := zap.NewNop()
	mw := corsMiddleware.New(logger)

	t.Run("OPTIONS request", func(t *testing.T) {
		nextHandlerCalled := false
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextHandlerCalled = true
		})

		req := httptest.NewRequest(http.MethodOptions, "http://example.com", nil)
		rr := httptest.NewRecorder()

		mw.Middleware(nextHandler).ServeHTTP(rr, req)

		if nextHandlerCalled {
			t.Errorf("expected next handler not to be called on OPTIONS request")
		}

		checkCommonHeaders(t, rr)
		if rr.Code != http.StatusOK {
			t.Errorf("expected status OK for OPTIONS, got %d", rr.Code)
		}
	})

	t.Run("GET request", func(t *testing.T) {
		nextHandlerCalled := false
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextHandlerCalled = true
			w.WriteHeader(http.StatusTeapot) // для проверки, что хендлер был вызван
		})

		req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		req.Header.Set("Origin", "http://example-origin.com")
		rr := httptest.NewRecorder()

		mw.Middleware(nextHandler).ServeHTTP(rr, req)

		if !nextHandlerCalled {
			t.Errorf("expected next handler to be called on non-OPTIONS request")
		}

		checkCommonHeaders(t, rr)
		if got, want := rr.Header().Get("Access-Control-Allow-Origin"), "http://example-origin.com"; got != want {
			t.Errorf("unexpected Access-Control-Allow-Origin header: got %s, want %s", got, want)
		}

		if rr.Code != http.StatusTeapot {
			t.Errorf("expected status code from next handler to be 418, got %d", rr.Code)
		}
	})
}

func checkCommonHeaders(t *testing.T, rr *httptest.ResponseRecorder) {
	t.Helper()
	resp := rr.Result()
	defer resp.Body.Close()
	headers := resp.Header
	if got, want := headers.Get("Access-Control-Allow-Methods"), "POST,PUT,DELETE,GET"; got != want {
		t.Errorf("unexpected Access-Control-Allow-Methods: got %s, want %s", got, want)
	}
	if got, want := headers.Get("Access-Control-Allow-Headers"), "Content-Type"; got != want {
		t.Errorf("unexpected Access-Control-Allow-Headers: got %s, want %s", got, want)
	}
	if got, want := headers.Get("Access-Control-Allow-Credentials"), "true"; got != want {
		t.Errorf("unexpected Access-Control-Allow-Credentials: got %s, want %s", got, want)
	}

}

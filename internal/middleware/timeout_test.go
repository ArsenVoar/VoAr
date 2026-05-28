package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTimeoutMiddleware(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deadline, ok := r.Context().Deadline()

		if !ok {
			t.Errorf("deadline not attached to context")
		}

		_ = deadline
	})

	handler := TimeOutMiddleware(testHandler)

	request := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)
}

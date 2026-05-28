package middleware

import (
	"VoAr/internal/contextkeys"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestIDMiddleware(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := contextkeys.GetRequestID(r.Context())

		if requestID == "" {
			t.Errorf("expected request_id, got empty string")
		}
	})

	handler := RequestIDMiddleware(testHandler)

	request := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)
}

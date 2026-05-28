package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()

		if r.Context().Err() == nil {
			t.Errorf("expected context cancellation")
		}

		done <- struct{}{}
	})

	request := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	).WithContext(ctx)

	recorder := httptest.NewRecorder()

	go testHandler.ServeHTTP(recorder, request)

	cancel()
	<-done
}

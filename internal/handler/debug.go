package handler

import (
	"VoAr/internal/logger"
	"net/http"
	"time"
)

// SlowHandler is used only for timeout middleware testing.
func (h *Handler) SlowHandler(w http.ResponseWriter, r *http.Request) {
	select {
	case <-time.After(20 * time.Second):
		if _, err := w.Write([]byte("slow response")); err != nil {
			logger.Error(err.Error())
		}

	case <-r.Context().Done():
		http.Error(
			w,
			r.Context().Err().Error(),
			http.StatusRequestTimeout,
		)
		return
	}
}

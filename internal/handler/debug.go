package handler

import (
	"net/http"
	"time"
)

func (h *Handler) SlowHandler(w http.ResponseWriter, r *http.Request) {
	select {
	case <-time.After(20 * time.Second):
		w.Write([]byte("slow response"))

	case <-r.Context().Done():
		http.Error(
			w,
			r.Context().Err().Error(),
			http.StatusRequestTimeout,
		)
		return
	}
}

package mw

import (
	"log/slog"
	"net/http"

	"github.com/Nemagu/dnd/internal/logger/sl"
	"github.com/google/uuid"
)

func LogRequestID(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			requestID := uuid.New()
			ctx := sl.WithLogRequestID(r.Context(), requestID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

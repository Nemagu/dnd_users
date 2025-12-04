package mw

import (
	"log/slog"
	"net/http"
	"time"

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

func LogRequest(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now().UTC()
			defer func() {
				logger.Info("request",
					"method", r.Method,
					"path", r.URL.Path,
					"remote_addr", r.RemoteAddr,
					"user_agent", r.UserAgent(),
					"duration", time.Since(start),
					"start", start.Format(time.RFC3339),
					"end", time.Now().UTC().Format(time.RFC3339),
				)
			}()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

package mw

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Nemagu/dnd/internal/logger/sl"
	weberror "github.com/Nemagu/dnd/internal/port/http/web/error"
	webservice "github.com/Nemagu/dnd/internal/port/http/web/service"
)

type userIDKey int

const UserIDKey = userIDKey(0)

type ErrorParser interface {
	Parse(ctx context.Context, err error) *weberror.ResponseError
}

type ResponseEncoder interface {
	Encode(ctx context.Context, w http.ResponseWriter, statusCode int, response any) error
}

type JWTAuth struct {
	logger          *slog.Logger
	provider        *webservice.JWTProvider
	errorParser     ErrorParser
	responseEncoder ResponseEncoder
}

func MustNewJWTAuth(
	logger *slog.Logger,
	provider *webservice.JWTProvider,
	errorParser ErrorParser,
	responseEncoder ResponseEncoder,
) *JWTAuth {
	if logger == nil {
		panic("jwt auth middleware does not get logger")
	}
	if provider == nil {
		panic("jwt auth middleware does not get service")
	}
	if errorParser == nil {
		panic("jwt auth middleware does not get error parser")
	}
	if responseEncoder == nil {
		panic("jwt auth middleware does not get response encoder")
	}
	return &JWTAuth{
		logger:          logger,
		provider:        provider,
		errorParser:     errorParser,
		responseEncoder: responseEncoder,
	}
}

func (m *JWTAuth) Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		if len(authHeader) < 10 || authHeader[:7] != "Bearer " {
			m.responseEncoder.Encode(r.Context(), w, http.StatusUnauthorized, &weberror.ResponseError{
				Detail: "invalid token",
			})
			return
		}
		token := authHeader[7:]
		claims, err := m.provider.ValidateToken(token)
		if err != nil {
			rerr := m.errorParser.Parse(r.Context(), err)
			m.responseEncoder.Encode(r.Context(), w, rerr.StatusCode, rerr)
			return
		}
		userID, err := m.provider.UserID(claims)
		if err != nil {
			rerr := m.errorParser.Parse(r.Context(), err)
			m.responseEncoder.Encode(r.Context(), w, rerr.StatusCode, rerr)
			return
		}
		ctx = context.WithValue(ctx, UserIDKey, userID)
		ctx = sl.WithLogUserID(ctx, userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

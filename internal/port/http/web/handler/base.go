package handler

import (
	"context"
	"log/slog"
	"net/http"

	weberror "github.com/Nemagu/dnd/internal/port/http/web/error"
	"github.com/Nemagu/dnd/internal/port/http/web/mw"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ErrorParser interface {
	Parse(ctx context.Context, err error) *weberror.ResponseError
}

type ResponseEncoder interface {
	Encode(ctx context.Context, w http.ResponseWriter, statusCode int, response any)
}

type RequestDecoder interface {
	Decode(ctx context.Context, r *http.Request, request any) error
}

type BaseHandler struct {
	logger          *slog.Logger
	errorParser     ErrorParser
	responseEncoder ResponseEncoder
	requestDecoder  RequestDecoder
}

func MustNewBaseHandler(
	logger *slog.Logger,
	errorParser ErrorParser,
	responseEncoder ResponseEncoder,
	requestDecoder RequestDecoder,
) *BaseHandler {
	if logger == nil {
		panic("base handler does not get logger")
	}
	if errorParser == nil {
		panic("base handler does not get error parser")
	}
	if responseEncoder == nil {
		panic("base handler does not get response encoder")
	}
	if requestDecoder == nil {
		panic("base handler does not get request decoder")
	}
	return &BaseHandler{
		logger:          logger,
		errorParser:     errorParser,
		responseEncoder: responseEncoder,
		requestDecoder:  requestDecoder,
	}
}

func (h *BaseHandler) extractAuthUserID(ctx context.Context) (uuid.UUID, error) {
	ctxUserID := ctx.Value(mw.UserIDKey)
	userID, ok := ctxUserID.(uuid.UUID)
	if !ok {
		h.logger.ErrorContext(ctx, "can not extract user id from context")
		return uuid.Nil, &weberror.ResponseError{
			StatusCode: http.StatusInternalServerError,
			Detail:     "can not extract user id from context",
		}
	}
	return userID, nil
}

func (h *BaseHandler) extractPathUserID(ctx context.Context, r *http.Request) (uuid.UUID, error) {
	userIDStr := chi.URLParam(r, "userID")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.logger.InfoContext(
			ctx,
			"can not parse user id from url",
			"parse_error",
			err,
			"user_id_tried",
			userIDStr,
		)
		return userID, &weberror.ResponseError{
			StatusCode: http.StatusBadRequest,
			Detail:     "не корректный id пользователя",
		}
	}

	return userID, nil
}

func (h *BaseHandler) handleError(ctx context.Context, w http.ResponseWriter, err error) {
	h.logger.InfoContext(ctx, "handle error", "error", err)
	rerr := h.errorParser.Parse(ctx, err)
	w.WriteHeader(rerr.StatusCode)
	h.responseEncoder.Encode(ctx, w, rerr.StatusCode, rerr)
}

package handler

import (
	"context"
	"log/slog"
	"net/http"

	weberror "github.com/Nemagu/dnd/internal/port/http/web/error"
)

type ErrorParser interface {
	Parse(ctx context.Context, err error) *weberror.ResponseError
}

type ResponseEncoder interface {
	Encode(ctx context.Context, w http.ResponseWriter, statusCode int, response any) error
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

func (h *BaseHandler) errorHandle(ctx context.Context, w http.ResponseWriter, err error) {
	rerr := h.errorParser.Parse(ctx, err)
	w.WriteHeader(rerr.StatusCode)
	h.responseEncoder.Encode(ctx, w, rerr.StatusCode, rerr)
}

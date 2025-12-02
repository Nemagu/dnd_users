package webservice

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Nemagu/dnd/internal/application"
	weberror "github.com/Nemagu/dnd/internal/port/http/web/error"
)

type ErrorParser struct {
	logger *slog.Logger
}

func MustNewErrorParser(logger *slog.Logger) *ErrorParser {
	if logger == nil {
		panic("error parser does not get logger")
	}
	return &ErrorParser{
		logger: logger,
	}
}

func (rw *ErrorParser) Parse(ctx context.Context, err error) *weberror.ResponseError {
	var responseError *weberror.ResponseError
	var status int
	switch {
	case errors.As(err, &responseError):
		return responseError
	case errors.Is(
		err, application.ErrValidation,
	) || errors.Is(
		err, application.ErrAlreadyExists,
	) || errors.Is(err, application.ErrVersionConflict):
		status = http.StatusBadRequest
	case errors.Is(err, application.ErrCredential):
		status = http.StatusUnauthorized
	case errors.Is(err, application.ErrNotAllowed):
		status = http.StatusForbidden
	case errors.Is(err, application.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, application.ErrInternal):
		status = http.StatusInternalServerError
	default:
		return rw.parseError(ctx, err)
	}
	return rw.parseApplicationError(ctx, status, err)
}

func (rw *ErrorParser) parseApplicationError(
	ctx context.Context, statusCode int, err error,
) *weberror.ResponseError {
	rerr := &weberror.ResponseError{
		StatusCode: statusCode,
	}
	if statusCode >= 500 {
		rw.logger.ErrorContext(
			ctx,
			"internal server error",
			"error", err,
		)
		return rerr
	}
	rerr.Detail = err.Error()
	return rerr
}

func (rw *ErrorParser) parseError(ctx context.Context, err error) *weberror.ResponseError {
	rw.logger.ErrorContext(
		ctx,
		"internal server error",
		"error",
		err,
	)
	return &weberror.ResponseError{
		StatusCode: http.StatusInternalServerError,
	}
}

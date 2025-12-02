package webservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	weberror "github.com/Nemagu/dnd/internal/port/http/web/error"
)

type JSONRequestDecoder struct {
	logger *slog.Logger
}

func MustNewJSONRequestDecoder(logger *slog.Logger) *JSONRequestDecoder {
	if logger == nil {
		panic("json decoder does not get logger")
	}
	return &JSONRequestDecoder{
		logger: logger,
	}
}

func (d *JSONRequestDecoder) Decode(ctx context.Context, r *http.Request, request any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(request); err != nil {
		d.logger.InfoContext(ctx, "decode json body error", "decode_error", err)
		return d.parseJSONError(err)
	}
	return nil
}

func (d *JSONRequestDecoder) parseJSONError(err error) *weberror.ResponseError {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	var unknownFieldError interface{ Field() string }
	switch {
	case errors.As(err, &syntaxError):
		return &weberror.ResponseError{
			StatusCode: http.StatusBadRequest,
			Detail:     "the JSON syntax is invalid",
		}
	case errors.As(err, &unmarshalTypeError):
		return &weberror.ResponseError{
			StatusCode: http.StatusBadRequest,
			Detail:     "one or more fields have invalid types",
			Errors: []weberror.ValidationError{
				{
					Field: unmarshalTypeError.Field,
					Message: fmt.Sprintf(
						"expected type %s",
						unmarshalTypeError.Type,
					),
				},
			},
		}
	case errors.As(err, &unknownFieldError):
		return &weberror.ResponseError{
			StatusCode: http.StatusBadRequest,
			Detail:     "request contains unknown fields",
			Errors: []weberror.ValidationError{
				{
					Field:   unknownFieldError.Field(),
					Message: "this field is not allowed",
				},
			},
		}
	case errors.Is(err, io.EOF):
		return &weberror.ResponseError{
			StatusCode: http.StatusBadRequest,
			Detail:     "request body is empty",
		}
	default:
		return &weberror.ResponseError{
			StatusCode: http.StatusBadRequest,
			Detail:     "request body is invalid",
		}
	}
}

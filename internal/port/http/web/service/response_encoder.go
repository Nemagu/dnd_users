package webservice

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type JSONResponseEncoder struct {
	logger *slog.Logger
}

func MustNewJSONResponseEncoder(logger *slog.Logger) *JSONResponseEncoder {
	if logger == nil {
		panic("json response encoder does not get logger")
	}
	return &JSONResponseEncoder{
		logger: logger,
	}
}

func (e *JSONResponseEncoder) Encode(
	ctx context.Context, w http.ResponseWriter, statusCode int, response any,
) {
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		e.logger.ErrorContext(
			ctx, "encode json error",
			"encode_error", err,
		)
	}
}

package sl

import (
	"context"
	"log/slog"
	"os"

	"github.com/google/uuid"
)

func MustNewJSONLogger(level slog.Leveler) *slog.Logger {
	if level == nil {
		panic("logger does not get level")
	}
	handler := slog.Handler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	}))
	handler = newHandlerMiddleware(handler)
	return slog.New(handler)
}

func MustNewTextLogger(level slog.Leveler) *slog.Logger {
	if level == nil {
		panic("logger does not get level")
	}
	handler := slog.Handler(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	}))
	handler = newHandlerMiddleware(handler)
	return slog.New(handler)
}

func WithLogUserID(ctx context.Context, userID uuid.UUID) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.UserID = userID
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{UserID: userID})
}

func WithLogRequestID(ctx context.Context, requestID uuid.UUID) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.RequestID = requestID
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{RequestID: requestID})
}

type handlerMiddleware struct {
	next slog.Handler
}

func newHandlerMiddleware(next slog.Handler) *handlerMiddleware {
	return &handlerMiddleware{
		next: next,
	}
}

func (h *handlerMiddleware) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *handlerMiddleware) Handle(ctx context.Context, rec slog.Record) error {
	if c, ok := ctx.Value(key).(logCtx); ok {
		if c.UserID != uuid.Nil {
			rec.Add("user_id", c.UserID)
		}
		if c.RequestID != uuid.Nil {
			rec.Add("request_id", c.RequestID)
		}
	}
	return h.next.Handle(ctx, rec)
}

func (h *handlerMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handlerMiddleware{next: h.next.WithAttrs(attrs)}
}

func (h *handlerMiddleware) WithGroup(name string) slog.Handler {
	return &handlerMiddleware{next: h.next.WithGroup(name)}
}

type keyType int

const key = keyType(0)

type logCtx struct {
	UserID    uuid.UUID
	RequestID uuid.UUID
}

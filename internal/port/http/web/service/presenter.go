package webservice

import (
	"context"
	"log/slog"

	appdto "github.com/Nemagu/dnd/internal/application/dto"
	webschema "github.com/Nemagu/dnd/internal/port/http/web/schema"
)

type UserPresenter struct {
	logger *slog.Logger
}

func MustNewUserPresenter(logger *slog.Logger) *UserPresenter {
	return &UserPresenter{
		logger: logger,
	}
}

func (p *UserPresenter) Present(ctx context.Context, user *appdto.User) *webschema.UserResponse {
	p.logger.DebugContext(ctx, "present user", "user", user)
	return &webschema.UserResponse{
		UserID: user.UserID,
		Email:  user.Email,
		State:  user.State,
		Status: user.Status,
	}
}

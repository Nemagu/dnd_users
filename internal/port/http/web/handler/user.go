package handler

import (
	"context"
	"net/http"

	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/Nemagu/dnd/internal/application/usecase"
	webschema "github.com/Nemagu/dnd/internal/port/http/web/schema"
)

type UserPresenter interface {
	Present(ctx context.Context, user *appdto.User) *webschema.UserResponse
}

type GetUserHandler struct {
	BaseHandler
	useCase   *usecase.UserUseCase
	presenter UserPresenter
}

func MustNewGetUserHandler(
	base BaseHandler, useCase *usecase.UserUseCase, presenter UserPresenter,
) *GetUserHandler {
	if useCase == nil {
		panic("get user handler did not get use case")
	}
	if presenter == nil {
		panic("get user handler did not get presenter")
	}
	return &GetUserHandler{
		BaseHandler: base,
		useCase:     useCase,
		presenter:   presenter,
	}
}

func (h *GetUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID, err := h.BaseHandler.extractPathUserID(r.Context(), r)
	if err != nil {
		h.BaseHandler.handleError(r.Context(), w, err)
		return
	}

	initiatorID, err := h.BaseHandler.extractAuthUserID(r.Context())
	if err != nil {
		h.BaseHandler.handleError(r.Context(), w, err)
		return
	}

	input := &appdto.UserQuery{
		InitiatorID: initiatorID,
		UserID:      userID,
	}
	appUser, err := h.useCase.Execute(r.Context(), input)
	if err != nil {
		h.BaseHandler.handleError(r.Context(), w, err)
		return
	}

	h.responseEncoder.Encode(
		r.Context(),
		w,
		http.StatusOK,
		h.presenter.Present(r.Context(), appUser),
	)
}

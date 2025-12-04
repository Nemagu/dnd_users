package handler

import (
	"net/http"

	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/Nemagu/dnd/internal/application/usecase"
	webschema "github.com/Nemagu/dnd/internal/port/http/web/schema"
)

type ChangeUserHandler struct {
	BaseHandler
	useCase *usecase.ChangeUserUseCase
}

func MustNewChangeUserHandler(
	base BaseHandler,
	useCase *usecase.ChangeUserUseCase,
) *ChangeUserHandler {
	if useCase == nil {
		panic("change user handler does not get use case")
	}
	return &ChangeUserHandler{
		BaseHandler: base,
		useCase:     useCase,
	}
}

func (h *ChangeUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body *webschema.ChangeUserRequest
	if err := h.BaseHandler.requestDecoder.Decode(r.Context(), r, body); err != nil {
		h.BaseHandler.handleError(r.Context(), w, err)
		return
	}

	initiatorID, err := h.BaseHandler.extractAuthUserID(r.Context())
	if err != nil {
		h.BaseHandler.handleError(r.Context(), w, err)
		return
	}

	userID, err := h.BaseHandler.extractPathUserID(r.Context(), r)
	if err != nil {
		h.BaseHandler.handleError(r.Context(), w, err)
		return
	}

	input := &appdto.ChangeUserCommand{
		InitiatorID: initiatorID,
		UserID:      userID,
		Email:       body.Email,
		State:       body.State,
		Status:      body.Status,
		Password:    body.Password,
	}
	if err = h.useCase.Execute(r.Context(), input); err != nil {
		h.BaseHandler.handleError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

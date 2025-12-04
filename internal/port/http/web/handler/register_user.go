package handler

import (
	"net/http"

	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/Nemagu/dnd/internal/application/usecase"
	webschema "github.com/Nemagu/dnd/internal/port/http/web/schema"
	"github.com/google/uuid"
)

type RegisterUserHandler struct {
	BaseHandler
	useCase *usecase.RegisterUserUseCase
}

func MustNewRegisterUserHandler(
	base BaseHandler, useCase *usecase.RegisterUserUseCase,
) *RegisterUserHandler {
	if useCase == nil {
		panic("confirm email handler does not get use case")
	}
	return &RegisterUserHandler{
		BaseHandler: base,
		useCase:     useCase,
	}
}

func (h *RegisterUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body *webschema.RegisterUserRequest

	if err := h.BaseHandler.requestDecoder.Decode(r.Context(), r, &body); err != nil {
		h.BaseHandler.handleError(r.Context(), w, err)
		return
	}
	input := appdto.RegisterUserCommand{
		Token:    body.Token,
		Password: body.Password,
	}
	id, err := h.useCase.Execute(r.Context(), &input)
	if err != nil {
		h.BaseHandler.handleError(r.Context(), w, err)
		return
	}
	response := struct {
		UserID uuid.UUID `json:"user_id"`
	}{
		UserID: id,
	}
	h.BaseHandler.responseEncoder.Encode(r.Context(), w, http.StatusOK, response)
}

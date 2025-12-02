package handler

import (
	"net/http"

	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/Nemagu/dnd/internal/application/usecase"
)

type ConfirmEmailHandler struct {
	BaseHandler
	useCase *usecase.ConfirmEmailUseCase
}

func MustNewConfirmEmailHandler(
	base BaseHandler, useCase *usecase.ConfirmEmailUseCase,
) *ConfirmEmailHandler {
	if useCase == nil {
		panic("confirm email handler does not get use case")
	}
	return &ConfirmEmailHandler{
		BaseHandler: base,
		useCase:     useCase,
	}
}

func (h *ConfirmEmailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string `json:"email"`
	}

	if err := h.BaseHandler.requestDecoder.Decode(r.Context(), r, &body); err != nil {
		h.BaseHandler.errorHandle(r.Context(), w, err)
		return
	}
	input := appdto.ConfirmEmailCommand{
		Email: body.Email,
	}
	if err := h.useCase.Execute(r.Context(), &input); err != nil {
		h.BaseHandler.errorHandle(r.Context(), w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

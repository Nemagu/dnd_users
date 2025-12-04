package handler

import (
	"net/http"

	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/Nemagu/dnd/internal/application/usecase"
)

type ConfirmResetPasswordHandler struct {
	BaseHandler
	useCase *usecase.ConfirmResetPasswordUseCase
}

func MustNewConfirmResetPasswordHandler(
	base BaseHandler, useCase *usecase.ConfirmResetPasswordUseCase,
) *ConfirmResetPasswordHandler {
	if useCase == nil {
		panic("confirm reset password handler does not get use case")
	}
	return &ConfirmResetPasswordHandler{
		BaseHandler: base,
		useCase:     useCase,
	}
}

func (h *ConfirmResetPasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string `json:"email"`
	}

	if err := h.BaseHandler.requestDecoder.Decode(r.Context(), r, &body); err != nil {
		h.BaseHandler.handleError(r.Context(), w, err)
		return
	}

	input := appdto.ConfirmResetPasswordCommand{
		Email: body.Email,
	}

	if err := h.useCase.Execute(r.Context(), input); err != nil {
		h.BaseHandler.handleError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

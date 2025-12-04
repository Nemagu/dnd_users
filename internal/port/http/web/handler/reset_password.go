package handler

import (
	"net/http"

	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/Nemagu/dnd/internal/application/usecase"
)

type ResetPasswordHandler struct {
	BaseHandler
	useCase *usecase.ResetPasswordUseCase
}

func MustNewResetPasswordHandler(
	base BaseHandler, useCase *usecase.ResetPasswordUseCase,
) *ResetPasswordHandler {
	if useCase == nil {
		panic("reset password handler does not get use case")
	}
	return &ResetPasswordHandler{
		BaseHandler: base,
		useCase:     useCase,
	}
}

func (h *ResetPasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := h.BaseHandler.requestDecoder.Decode(r.Context(), r, &body); err != nil {
		h.BaseHandler.handleError(r.Context(), w, err)
		return
	}

	input := &appdto.ResetPasswordCommand{
		Token:       body.Token,
		NewPassword: body.NewPassword,
	}

	if err := h.useCase.Execute(r.Context(), input); err != nil {
		h.BaseHandler.handleError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

package handler

import (
	"net/http"

	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/Nemagu/dnd/internal/application/usecase"
	webschema "github.com/Nemagu/dnd/internal/port/http/web/schema"
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
	var body *webschema.ResetPasswordRequest
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

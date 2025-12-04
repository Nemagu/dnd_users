package handler

import (
	"net/http"

	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/Nemagu/dnd/internal/application/usecase"
	webschema "github.com/Nemagu/dnd/internal/port/http/web/schema"
)

type ConfirmNewEmailHandler struct {
	BaseHandler
	useCase *usecase.ConfirmNewEmailUseCase
}

func MustNewConfirmNewEmailHandler(
	base BaseHandler, useCase *usecase.ConfirmNewEmailUseCase,
) *ConfirmNewEmailHandler {
	if useCase == nil {
		panic("confirm new email handler does not get use case")
	}
	return &ConfirmNewEmailHandler{
		BaseHandler: base,
		useCase:     useCase,
	}
}

func (h *ConfirmNewEmailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body *webschema.ConfirmNewEmailRequest
	if err := h.BaseHandler.requestDecoder.Decode(r.Context(), r, &body); err != nil {
		h.BaseHandler.handleError(r.Context(), w, err)
		return
	}

	userID, err := h.BaseHandler.extractAuthUserID(r.Context())
	if err != nil {
		h.BaseHandler.handleError(r.Context(), w, err)
		return
	}

	input := appdto.ConfirmNewEmailCommand{
		UserID:   userID,
		Email:    body.Email,
		Password: body.Password,
	}
	err = h.useCase.Execute(r.Context(), &input)
	if err != nil {
		h.BaseHandler.handleError(r.Context(), w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

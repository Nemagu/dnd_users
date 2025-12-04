package handler

import (
	"net/http"

	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/Nemagu/dnd/internal/application/usecase"
	webschema "github.com/Nemagu/dnd/internal/port/http/web/schema"
)

type ChangePasswordHandler struct {
	BaseHandler
	useCase *usecase.ChangePasswordUseCase
}

func MustNewChangePasswordHandler(
	base BaseHandler,
	useCase *usecase.ChangePasswordUseCase,
) *ChangePasswordHandler {
	if useCase == nil {
		panic("change password handler does not get use case")
	}
	return &ChangePasswordHandler{
		BaseHandler: base,
		useCase:     useCase,
	}
}

func (h *ChangePasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body *webschema.ChangePasswordRequest
	if err := h.BaseHandler.requestDecoder.Decode(r.Context(), r, &body); err != nil {
		h.BaseHandler.handleError(r.Context(), w, err)
		return
	}

	userID, err := h.BaseHandler.extractAuthUserID(r.Context())
	if err != nil {
		h.handleError(r.Context(), w, err)
		return
	}

	input := &appdto.ChangePasswordCommand{
		UserID:      userID,
		OldPassword: body.OldPassword,
		NewPassword: body.NewPassword,
	}

	if err = h.useCase.Execute(r.Context(), input); err != nil {
		h.BaseHandler.handleError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

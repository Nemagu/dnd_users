package handler

import (
	"net/http"

	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/Nemagu/dnd/internal/application/usecase"
	"github.com/Nemagu/dnd/internal/port/http/web/mw"
	"github.com/google/uuid"
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
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := h.BaseHandler.requestDecoder.Decode(r.Context(), r, &body); err != nil {
		h.BaseHandler.errorHandle(r.Context(), w, err)
		return
	}
	ctxUserID := r.Context().Value(mw.UserIDKey)
	userID, ok := ctxUserID.(uuid.UUID)
	if !ok {
		h.BaseHandler.logger.ErrorContext(r.Context(), "register handler did not get user id from context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	input := appdto.ConfirmNewEmailCommand{
		UserID:   userID,
		Email:    body.Email,
		Password: body.Password,
	}
	err := h.useCase.Execute(r.Context(), &input)
	if err != nil {
		h.BaseHandler.errorHandle(r.Context(), w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

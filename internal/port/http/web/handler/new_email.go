package handler

import (
	"net/http"

	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/Nemagu/dnd/internal/application/usecase"
	webschema "github.com/Nemagu/dnd/internal/port/http/web/schema"
)

type NewEmailHandler struct {
	BaseHandler
	useCase *usecase.NewEmailUseCase
}

func MustNewEmailHandler(base BaseHandler, useCase *usecase.NewEmailUseCase) *NewEmailHandler {
	if useCase == nil {
		panic("new email handler did not get use case")
	}
	return &NewEmailHandler{
		BaseHandler: base,
		useCase:     useCase,
	}
}

func (h *NewEmailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body *webschema.NewEmailRequest
	if err := h.requestDecoder.Decode(r.Context(), r, &body); err != nil {
		h.handleError(r.Context(), w, err)
		return
	}

	input := &appdto.NewEmailCommand{
		Token:    body.Token,
		Password: body.Password,
	}
	if err := h.useCase.Execute(r.Context(), input); err != nil {
		h.handleError(r.Context(), w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

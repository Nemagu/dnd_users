package handler

import (
	"net/http"

	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/Nemagu/dnd/internal/application/usecase"
)

type GetMeHandler struct {
	BaseHandler
	useCase   *usecase.UserUseCase
	presenter UserPresenter
}

func MustNewGetMeHandler(
	base BaseHandler,
	useCase *usecase.UserUseCase,
	presenter UserPresenter,
) *GetMeHandler {
	if useCase == nil {
		panic("get me handler did not get use case")
	}
	if presenter == nil {
		panic("get me handler did not get presenter")
	}
	return &GetMeHandler{
		BaseHandler: base,
		useCase:     useCase,
		presenter:   presenter,
	}
}

func (h *GetMeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	initiatorID, err := h.BaseHandler.extractAuthUserID(r.Context())
	if err != nil {
		h.BaseHandler.handleError(r.Context(), w, err)
		return
	}

	input := &appdto.UserQuery{
		InitiatorID: initiatorID,
		UserID:      initiatorID,
	}
	appUser, err := h.useCase.Execute(r.Context(), input)
	if err != nil {
		h.BaseHandler.handleError(r.Context(), w, err)
		return
	}

	h.responseEncoder.Encode(
		r.Context(),
		w,
		http.StatusOK,
		h.presenter.Present(r.Context(), appUser),
	)
}

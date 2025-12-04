package handler

import (
	"net/http"
	"strconv"

	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/Nemagu/dnd/internal/application/usecase"
	webschema "github.com/Nemagu/dnd/internal/port/http/web/schema"
)

type GetUsersHandler struct {
	BaseHandler
	useCase   *usecase.UsersUseCase
	presenter UserPresenter
}

func MustNewGetUsersHandler(
	base BaseHandler, useCase *usecase.UsersUseCase, presenter UserPresenter,
) *GetUsersHandler {
	if useCase == nil {
		panic("get users handler did not get use case")
	}
	if presenter == nil {
		panic("get users handler did not get presenter")
	}
	return &GetUsersHandler{
		BaseHandler: base,
		useCase:     useCase,
		presenter:   presenter,
	}
}

func (h *GetUsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	initiatorID, err := h.BaseHandler.extractAuthUserID(r.Context())
	if err != nil {
		h.BaseHandler.handleError(r.Context(), w, err)
		return
	}

	input := h.extractQueryParams(r)
	input.InitiatorID = initiatorID

	appUsers, err := h.useCase.Execute(r.Context(), input)
	if err != nil {
		h.BaseHandler.handleError(r.Context(), w, err)
		return
	}

	response := make([]*webschema.UserResponse, 0, len(appUsers))
	for _, u := range appUsers {
		response = append(response, h.presenter.Present(r.Context(), u))
	}

	h.responseEncoder.Encode(r.Context(), w, http.StatusOK, response)
}

func (h *GetUsersHandler) extractQueryParams(r *http.Request) *appdto.UsersQuery {
	query := r.URL.Query()

	email := query.Get("email")
	h.BaseHandler.logger.DebugContext(r.Context(), "email in query", "email", email)

	states := query["state"]
	h.BaseHandler.logger.DebugContext(r.Context(), "states in query", "states", states)

	statuses := query["status"]
	h.BaseHandler.logger.DebugContext(r.Context(), "statuses in query", "statuses", statuses)

	limitStr := query.Get("limit")
	h.logger.DebugContext(r.Context(), "limit in query", "limit", limitStr)
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		h.logger.WarnContext(r.Context(), "can not parse limit", "limit", limitStr)
	}

	offsetStr := query.Get("offset")
	h.logger.DebugContext(r.Context(), "offset in query", "offset", offsetStr)
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		h.logger.WarnContext(r.Context(), "can not parse offset", "offset", offsetStr)
	}

	return &appdto.UsersQuery{
		SearchByEmail:  email,
		FilterByState:  states,
		FilterByStatus: statuses,
		Limit:          limit,
		Offset:         offset,
	}
}

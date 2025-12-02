package handler

import (
	"net/http"

	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/Nemagu/dnd/internal/application/usecase"
	"github.com/google/uuid"
)

type JWTTokensProvider interface {
	GenerateTokens(userID uuid.UUID) (any, error)
}

type JWTAuthHandler struct {
	BaseHandler
	useCase     *usecase.AuthenticateUseCase
	jwtProvider JWTTokensProvider
}

func MustNewJWTAuthHandler(
	base BaseHandler,
	useCase *usecase.AuthenticateUseCase,
	jwtProvider JWTTokensProvider,
) *JWTAuthHandler {
	if useCase == nil {
		panic("auth handler does not get use case")
	}
	if jwtProvider == nil {
		panic("auth handler does not get jwt provider")
	}
	return &JWTAuthHandler{
		BaseHandler: base,
		useCase:     useCase,
		jwtProvider: jwtProvider,
	}
}

func (h *JWTAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := h.BaseHandler.requestDecoder.Decode(r.Context(), r, &body); err != nil {
		h.BaseHandler.errorHandle(r.Context(), w, err)
		return
	}
	input := appdto.AuthenticateCommand{
		Email:    body.Email,
		Password: body.Password,
	}
	userID, err := h.useCase.Execute(r.Context(), input)
	if err != nil {
		h.BaseHandler.errorHandle(r.Context(), w, err)
		return
	}
	tokens, err := h.jwtProvider.GenerateTokens(userID)
	if err != nil {
		h.BaseHandler.errorHandle(r.Context(), w, err)
		return
	}
	err = h.BaseHandler.responseEncoder.Encode(r.Context(), w, http.StatusOK, tokens)
	if err != nil {
		h.BaseHandler.errorHandle(r.Context(), w, err)
		return
	}
}

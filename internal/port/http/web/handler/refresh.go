package handler

import (
	"net/http"
)

type JWTRefreshProvider interface {
	RefreshToken(tokenString string) (any, error)
}

type JWTRefreshHandler struct {
	BaseHandler
	jwtProvider JWTRefreshProvider
}

func MustNewJWTRefreshHandler(
	base BaseHandler, jwtProvider JWTRefreshProvider,
) *JWTRefreshHandler {
	if jwtProvider == nil {
		panic("jwt refresh handler does not get jwt provider")
	}
	return &JWTRefreshHandler{
		BaseHandler: base,
		jwtProvider: jwtProvider,
	}
}

func (h *JWTRefreshHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := h.BaseHandler.requestDecoder.Decode(r.Context(), r, &body); err != nil {
		h.BaseHandler.errorHandle(r.Context(), w, err)
		return
	}
	tokens, err := h.jwtProvider.RefreshToken(body.RefreshToken)
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

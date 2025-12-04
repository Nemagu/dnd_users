package webservice

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Nemagu/dnd/internal/config"
	weberror "github.com/Nemagu/dnd/internal/port/http/web/error"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

type JWTToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type JWTProvider struct {
	logger          *slog.Logger
	secret          []byte
	accessLifetime  time.Duration
	refreshLifetime time.Duration
}

func MustNewJWTProvider(
	logger *slog.Logger,
	cfg *config.WebConfig,
) *JWTProvider {
	return &JWTProvider{
		logger:          logger,
		secret:          []byte(cfg.JWTSecretKey),
		accessLifetime:  cfg.JWTAccessLifetime,
		refreshLifetime: cfg.JWTRefreshLifetime,
	}
}

func (p *JWTProvider) GenerateTokens(userID uuid.UUID) (any, error) {
	accessToken, err := p.generateToken(p.generateAccessClaims(userID))
	if err != nil {
		return JWTToken{}, err
	}
	refreshToken, err := p.generateToken(p.generateRefreshClaims(userID))
	if err != nil {
		return JWTToken{}, err
	}
	return JWTToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (p *JWTProvider) RefreshToken(tokenString string) (any, error) {
	clm, err := p.validateToken(tokenString)
	if err != nil {
		return JWTToken{}, err
	}
	token, err := p.generateToken(p.generateAccessClaims(clm.UserID))
	if err != nil {
		return JWTToken{}, err
	}
	return JWTToken{
		AccessToken:  token,
		RefreshToken: tokenString,
	}, nil
}

func (p *JWTProvider) ValidateToken(tokenString string) (any, error) {
	return p.validateToken(tokenString)
}

func (p *JWTProvider) UserID(clm any) (uuid.UUID, error) {
	if claims, ok := clm.(JWTClaims); ok {
		return claims.UserID, nil
	} else {
		return uuid.UUID{}, &weberror.ResponseError{
			StatusCode: http.StatusUnauthorized,
			Detail:     "invalid jwt token",
		}
	}
}

func (p *JWTProvider) validateToken(tokenString string) (JWTClaims, error) {
	claims := JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return claims, &weberror.ResponseError{
				StatusCode: http.StatusUnauthorized,
				Detail:     "invalid jwt token",
			}
		}
		return p.secret, nil
	})
	if err != nil {
		return claims, &weberror.ResponseError{
			StatusCode: http.StatusUnauthorized,
			Detail:     err.Error(),
		}
	}
	if !token.Valid {
		return claims, &weberror.ResponseError{
			StatusCode: http.StatusUnauthorized,
			Detail:     "invalid jwt token",
		}
	}
	return claims, nil
}

func (p *JWTProvider) generateToken(claims JWTClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenString, err := token.SignedString(p.secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (p *JWTProvider) generateAccessClaims(userID uuid.UUID) JWTClaims {
	return JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(p.accessLifetime)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}
}

func (p *JWTProvider) generateRefreshClaims(userID uuid.UUID) JWTClaims {
	return JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(p.refreshLifetime)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}
}

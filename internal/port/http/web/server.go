package web

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Nemagu/dnd/internal/adapter/repository/inmemory"
	"github.com/Nemagu/dnd/internal/adapter/service/email"
	"github.com/Nemagu/dnd/internal/adapter/service/password"
	"github.com/Nemagu/dnd/internal/application/usecase"
	"github.com/Nemagu/dnd/internal/config"
	"github.com/Nemagu/dnd/internal/port/http/web/handler"
	"github.com/Nemagu/dnd/internal/port/http/web/mw"
	webservice "github.com/Nemagu/dnd/internal/port/http/web/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HTTPServer struct {
	logger *slog.Logger
	config *config.WebConfig
}

func MustNewHTTPServer(logger *slog.Logger, config *config.WebConfig) *HTTPServer {
	if logger == nil {
		panic("http server does not get logger")
	}
	if config == nil {
		panic("http server does not get config")
	}
	return &HTTPServer{
		logger: logger,
		config: config,
	}
}

func (s *HTTPServer) MustServe() {
	r := chi.NewRouter()
	r.Use(middleware.Timeout(s.config.HTTPTimeout))
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(mw.LogRequestID(s.logger))
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Use(middleware.Recoverer)
	apiRouter := s.createAPIRouter()
	r.Mount("/auth/v1", apiRouter)

	if err := http.ListenAndServe(
		fmt.Sprintf("%s:%d", s.config.HTTPHost, s.config.HTTPPort), r,
	); err != nil {
		panic(err)
	}
}

func (s *HTTPServer) createAPIRouter() *chi.Mux {
	// Repositories
	userRepo := inmemory.MustNewInMemoryUserRepository()

	// Email application services
	emailValidator := email.MustNewEmailValidator()
	emailProvider := email.MustNewFileEmailProvider(
		s.logger, s.config.EmailFolderPath, s.config.EmailTimeout,
	)
	emailCrypter := email.MustNewBcryptEmailCrypter(
		s.config.EmailSecretKey, s.config.EmailTokenLifetime,
	)

	// Password application services
	passwordValidator := password.MustNewPasswordValidator(
		8, 64, true, true, true, true,
	)
	passwordHasher := password.MustNewBcryptPasswordHasher(s.config.PasswordCost)

	// Use cases
	confirmEmailUC := usecase.MustNewConfirmEmailUseCase(
		userRepo, emailCrypter, emailProvider, emailValidator,
	)
	registerUserUC := usecase.MustNewRegisterUserUseCase(
		userRepo, passwordValidator, passwordHasher, emailCrypter, emailValidator,
	)
	confirmNewEmailUC := usecase.MustNewConfirmNewEmailUseCase(
		userRepo, passwordHasher, emailCrypter, emailValidator, emailProvider,
	)
	authUC := usecase.MustNewAuthenticateUseCase(
		userRepo, passwordHasher,
	)

	// Web services
	jwtProvider := webservice.MustNewJWTProvider(s.logger, s.config)
	errorParser := webservice.MustNewErrorParser(s.logger)
	requestDecoder := webservice.MustNewJSONRequestDecoder(s.logger)
	responseEncoder := webservice.MustNewJSONResponseEncoder(s.logger)

	// Handlers need to use without auth middleware
	authMW := mw.MustNewJWTAuth(s.logger, jwtProvider, errorParser, responseEncoder)
	baseHandler := handler.MustNewBaseHandler(s.logger, errorParser, responseEncoder, requestDecoder)
	confirmEmailHandler := handler.MustNewConfirmEmailHandler(*baseHandler, confirmEmailUC)
	registerUserHandler := handler.MustNewRegisterUserHandler(*baseHandler, registerUserUC)

	// Handlers need to use with auth middleware
	confirmNewEmailHandler := handler.MustNewConfirmNewEmailHandler(*baseHandler, confirmNewEmailUC)
	authHandler := handler.MustNewJWTAuthHandler(*baseHandler, authUC, jwtProvider)
	refreshHandler := handler.MustNewJWTRefreshHandler(*baseHandler, jwtProvider)

	r := chi.NewRouter()
	authR := chi.NewRouter()

	authR.Use(authMW.Middleware)

	authR.Post("/email/new", confirmNewEmailHandler.ServeHTTP)

	r.Post("/jwt/create", authHandler.ServeHTTP)
	r.Post("/jwt/refresh", refreshHandler.ServeHTTP)
	r.Post("/email/confirm", confirmEmailHandler.ServeHTTP)
	r.Post("/register", registerUserHandler.ServeHTTP)
	r.Mount("/", authR)
	return r
}

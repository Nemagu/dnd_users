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
	"github.com/Nemagu/dnd/internal/domain/duser"
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
	r.Use(mw.LogRequest(s.logger))
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
	s.logger.Info("init user repository")

	userRepo := inmemory.MustNewInMemoryUserRepository()

	s.logger.Info("init email application services")

	emailValidator := email.MustNewEmailValidator()
	emailProvider := email.MustNewFileEmailProvider(
		s.logger,
		s.config.EmailFolderPath,
		s.config.EmailTimeout,
	)
	emailCrypter := email.MustNewBcryptEmailCrypter(
		s.config.EmailSecretKey,
		s.config.EmailTokenLifetime,
	)

	s.logger.Info("init password application services")

	passwordValidator := password.MustNewPasswordValidator(8, 64, true, true, true, true)
	passwordHasher := password.MustNewBcryptPasswordHasher(s.config.PasswordCost)

	s.logger.Info("init domain services")

	policyService := duser.MustNewPolicyService()

	s.logger.Info("init use cases")

	authUC := usecase.MustNewAuthenticateUseCase(userRepo, passwordHasher)
	changePasswordUC := usecase.MustNewChangePasswordUseCase(
		userRepo,
		passwordValidator,
		passwordHasher,
		passwordHasher,
	)
	changeUserUC := usecase.MustNewChangeUserUseCase(
		userRepo,
		policyService,
		passwordHasher,
		passwordValidator,
	)
	confirmEmailUC := usecase.MustNewConfirmEmailUseCase(
		userRepo,
		emailCrypter,
		emailProvider,
		emailValidator,
	)
	confirmNewEmailUC := usecase.MustNewConfirmNewEmailUseCase(
		userRepo,
		passwordHasher,
		emailCrypter,
		emailValidator,
		emailProvider,
	)
	confirmResetPasswordUC := usecase.MustNewConfirmResetPasswordUseCase(
		userRepo,
		emailValidator,
		emailCrypter,
		emailProvider,
	)
	registerUserUC := usecase.MustNewRegisterUserUseCase(
		userRepo,
		passwordValidator,
		passwordHasher,
		emailCrypter,
		emailValidator,
	)
	resetPasswordUC := usecase.MustNewResetPasswordUseCase(
		userRepo,
		emailCrypter,
		emailValidator,
		passwordValidator,
		passwordHasher,
	)
	userUC := usecase.MustNewUserUseCase(userRepo, policyService)
	usersUC := usecase.MustNewUsersUseCase(userRepo, policyService)

	s.logger.Info("init web services")

	jwtProvider := webservice.MustNewJWTProvider(s.logger, s.config)
	errorParser := webservice.MustNewErrorParser(s.logger)
	requestDecoder := webservice.MustNewJSONRequestDecoder(s.logger)
	responseEncoder := webservice.MustNewJSONResponseEncoder(s.logger)
	userPresenter := webservice.MustNewUserPresenter(s.logger)

	s.logger.Info("init auth middleware")

	authMW := mw.MustNewJWTAuth(s.logger, jwtProvider, errorParser, responseEncoder)

	r := chi.NewRouter()
	authR := chi.NewRouter()
	authR.Use(authMW.Middleware)

	s.logger.Info("init handlers without auth middleware")

	baseHandler := handler.MustNewBaseHandler(
		s.logger,
		errorParser,
		responseEncoder,
		requestDecoder,
	)
	confirmEmailHandler := handler.MustNewConfirmEmailHandler(*baseHandler, confirmEmailUC)
	registerUserHandler := handler.MustNewRegisterUserHandler(*baseHandler, registerUserUC)
	authHandler := handler.MustNewJWTAuthHandler(*baseHandler, authUC, jwtProvider)
	refreshHandler := handler.MustNewJWTRefreshHandler(*baseHandler, jwtProvider)
	confirmResetPasswordHandler := handler.MustNewConfirmResetPasswordHandler(
		*baseHandler,
		confirmResetPasswordUC,
	)
	resetPasswordHandler := handler.MustNewResetPasswordHandler(*baseHandler, resetPasswordUC)

	r.Post("/jwt/create", authHandler.ServeHTTP)
	r.Post("/jwt/refresh", refreshHandler.ServeHTTP)
	r.Post("/email/confirm", confirmEmailHandler.ServeHTTP)
	r.Post("/register", registerUserHandler.ServeHTTP)
	r.Post("/password/reset/confirm", confirmResetPasswordHandler.ServeHTTP)
	r.Post("/password/reset", resetPasswordHandler.ServeHTTP)

	s.logger.Info("init handlers with auth middleware")

	confirmNewEmailHandler := handler.MustNewConfirmNewEmailHandler(*baseHandler, confirmNewEmailUC)
	changePasswordHandler := handler.MustNewChangePasswordHandler(*baseHandler, changePasswordUC)
	changeUserHandler := handler.MustNewChangeUserHandler(*baseHandler, changeUserUC)
	userHandler := handler.MustNewGetUserHandler(*baseHandler, userUC, userPresenter)
	usersHandler := handler.MustNewGetUsersHandler(*baseHandler, usersUC, userPresenter)
	meHandler := handler.MustNewGetMeHandler(*baseHandler, userUC, userPresenter)

	authR.Post("/email/new", confirmNewEmailHandler.ServeHTTP)
	authR.Post("/password/change", changePasswordHandler.ServeHTTP)
	authR.Put("/users/{userID}", changeUserHandler.ServeHTTP)
	authR.Get("/users/{userID}", userHandler.ServeHTTP)
	authR.Get("/users", usersHandler.ServeHTTP)
	authR.Get("/users/me", meHandler.ServeHTTP)

	r.Mount("/", authR)
	return r
}

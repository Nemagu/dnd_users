package usecase

import (
	"context"
	"fmt"

	"github.com/Nemagu/dnd/internal/application"
	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/Nemagu/dnd/internal/application/repository"
	"github.com/Nemagu/dnd/internal/application/service"
	"github.com/google/uuid"
)

type AuthenticateUseCase struct {
	userRepo       repository.UserRepository
	passwordHasher service.PasswordHasher
}

func MustNewAuthenticateUseCase(
	userRepo repository.UserRepository, passwordHasher service.PasswordHasher,
) *AuthenticateUseCase {
	if userRepo == nil {
		panic("auth use case does not get user repository")
	}
	if passwordHasher == nil {
		panic("auth use case does not get password hasher")
	}
	return &AuthenticateUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
	}
}

func (u *AuthenticateUseCase) Execute(
	ctx context.Context, input appdto.AuthenticateCommand,
) (uuid.UUID, error) {
	stdErr := fmt.Errorf(
		"%w: не верный логин или пароль", application.ErrCredential,
	)
	user, err := u.userRepo.ByEmail(ctx, input.Email)
	if err != nil {
		return uuid.Nil, stdErr
	}
	compare, err := u.passwordHasher.ComparePassword(input.Password, user.PasswordHash)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	if !compare {
		return uuid.Nil, stdErr
	}
	return user.UserID, nil
}

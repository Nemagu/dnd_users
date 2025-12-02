package usecase

import (
	"context"
	"fmt"

	"github.com/Nemagu/dnd/internal/application"
	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/google/uuid"
)

type AuthenticateUserRepository interface {
	ByEmail(ctx context.Context, email string) (*appdto.User, error)
}

type AuthenticateUseCase struct {
	userRepo         AuthenticateUserRepository
	passwordComparer PasswordComparer
}

func MustNewAuthenticateUseCase(
	userRepo AuthenticateUserRepository, passwordComparer PasswordComparer,
) *AuthenticateUseCase {
	if userRepo == nil {
		panic("auth use case does not get user repository")
	}
	if passwordComparer == nil {
		panic("auth use case does not get password hasher")
	}
	return &AuthenticateUseCase{
		userRepo:         userRepo,
		passwordComparer: passwordComparer,
	}
}

func (u *AuthenticateUseCase) Execute(
	ctx context.Context, input *appdto.AuthenticateCommand,
) (uuid.UUID, error) {
	stdErr := fmt.Errorf(
		"%w: не верный логин или пароль", application.ErrCredential,
	)
	user, err := u.userRepo.ByEmail(ctx, input.Email)
	if err != nil {
		return uuid.Nil, stdErr
	}
	compare, err := u.passwordComparer.ComparePassword(input.Password, user.PasswordHash)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	if !compare {
		return uuid.Nil, stdErr
	}
	return user.UserID, nil
}

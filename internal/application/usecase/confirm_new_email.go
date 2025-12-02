package usecase

import (
	"context"
	"fmt"

	"github.com/Nemagu/dnd/internal/application"
	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/Nemagu/dnd/internal/application/repository"
	"github.com/Nemagu/dnd/internal/application/service"
)

type ConfirmNewEmailUseCase struct {
	userRepo       repository.UserRepository
	passwordHasher service.PasswordHasher
	emailCrypter   service.EmailCrypter
	emailValidator service.EmailValidator
	emailProvider  service.EmailProvider
}

func MustNewConfirmNewEmailUseCase(
	userRepo repository.UserRepository,
	passwordHasher service.PasswordHasher,
	emailCrypter service.EmailCrypter,
	emailValidator service.EmailValidator,
	emailProvider service.EmailProvider,
) *ConfirmNewEmailUseCase {
	return &ConfirmNewEmailUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		emailCrypter:   emailCrypter,
		emailValidator: emailValidator,
		emailProvider:  emailProvider,
	}
}

func (u *ConfirmNewEmailUseCase) Execute(
	ctx context.Context,
	input appdto.ConfirmNewEmailCommand,
) error {
	user, err := u.userRepo.ByID(ctx, input.UserID)
	if err != nil {
		return err
	}

	compare, err := u.passwordHasher.ComparePassword(
		input.Password,
		user.PasswordHash,
	)
	if err != nil {
		return err
	}
	if !compare {
		return fmt.Errorf("%w: не верный логи или пароль", application.ErrValidation)
	}

	if err = u.emailValidator.Validate(user.Email); err != nil {
		return err
	}

	token, err := u.emailCrypter.Encrypt(user.Email)
	if err != nil {
		return err
	}

	go u.emailProvider.SendChangeEmail(
		appdto.Email{To: user.Email, Token: token},
	)

	return nil
}

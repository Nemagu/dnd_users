package usecase

import (
	"context"
	"fmt"

	"github.com/Nemagu/dnd/internal/application"
	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/google/uuid"
)

type ConfirmNewEmailUserRepository interface {
	ByID(
		ctx context.Context,
		id uuid.UUID,
	) (*appdto.User, error)
}

type PasswordComparer interface {
	ComparePassword(
		password string,
		hash string,
	) (bool, error)
}

type EmailCrypter interface {
	Encrypt(email string) (string, error)
}

type ConfirmNewEmailProvider interface {
	SendChangeEmail(
		message appdto.Email,
	)
}

type ConfirmNewEmailUseCase struct {
	userRepo         ConfirmNewEmailUserRepository
	passwordComparer PasswordComparer
	emailCrypter     EmailCrypter
	emailValidator   EmailValidator
	emailProvider    ConfirmNewEmailProvider
}

func MustNewConfirmNewEmailUseCase(
	userRepo ConfirmNewEmailUserRepository,
	passwordComparer PasswordComparer,
	emailCrypter EmailCrypter,
	emailValidator EmailValidator,
	emailProvider ConfirmNewEmailProvider,
) *ConfirmNewEmailUseCase {
	return &ConfirmNewEmailUseCase{
		userRepo:         userRepo,
		passwordComparer: passwordComparer,
		emailCrypter:     emailCrypter,
		emailValidator:   emailValidator,
		emailProvider:    emailProvider,
	}
}

func (u *ConfirmNewEmailUseCase) Execute(
	ctx context.Context,
	input *appdto.ConfirmNewEmailCommand,
) error {
	user, err := u.userRepo.ByID(ctx, input.UserID)
	if err != nil {
		return err
	}

	compare, err := u.passwordComparer.ComparePassword(
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

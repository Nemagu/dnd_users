package usecase

import (
	"context"
	"fmt"

	"github.com/Nemagu/dnd/internal/application"
	appdto "github.com/Nemagu/dnd/internal/application/dto"
)

type ConfirmResetPasswordUserRepository interface {
	EmailExists(
		ctx context.Context,
		email string,
	) (bool, error)
}

type ConfirmResetPasswordEmailProvider interface {
	SendResetPasswordEmail(
		message appdto.Email,
	)
}

type ConfirmResetPasswordUseCase struct {
	userRepo       ConfirmResetPasswordUserRepository
	emailValidator EmailValidator
	emailCrypter   EmailCrypter
	emailProvider  ConfirmResetPasswordEmailProvider
}

func (u *ConfirmResetPasswordUseCase) Execute(
	ctx context.Context, input appdto.ConfirmResetPasswordCommand,
) error {
	if err := u.emailValidator.Validate(input.Email); err != nil {
		return err
	}

	exists, err := u.userRepo.EmailExists(ctx, input.Email)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%w: такого email не существует", application.ErrValidation)
	}

	token, err := u.emailCrypter.Encrypt(input.Email)
	if err != nil {
		return err
	}

	go u.emailProvider.SendResetPasswordEmail(appdto.Email{
		To:    input.Email,
		Token: token,
	})

	return nil
}

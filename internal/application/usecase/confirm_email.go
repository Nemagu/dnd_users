package usecase

import (
	"context"
	"fmt"

	"github.com/Nemagu/dnd/internal/application"
	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/Nemagu/dnd/internal/application/repository"
	"github.com/Nemagu/dnd/internal/application/service"
	"github.com/Nemagu/dnd/internal/domain"
)

type ConfirmEmailUseCase struct {
	userRepo       repository.UserRepository
	emailCrypter   service.EmailCrypter
	emailProvider  service.EmailProvider
	emailValidator service.EmailValidator
}

func MustNewConfirmEmailUseCase(
	userRepo repository.UserRepository,
	emailCrypter service.EmailCrypter,
	emailProvider service.EmailProvider,
	emailValidator service.EmailValidator,
) *ConfirmEmailUseCase {
	return &ConfirmEmailUseCase{
		userRepo:       userRepo,
		emailCrypter:   emailCrypter,
		emailProvider:  emailProvider,
		emailValidator: emailValidator,
	}
}

func (u *ConfirmEmailUseCase) Execute(
	ctx context.Context,
	input appdto.ConfirmEmailCommand,
) error {
	if err := u.emailValidator.Validate(input.Email); err != nil {
		return err
	}
	email, err := domain.NewEmail(input.Email)
	if err != nil {
		return err
	}
	exists, err := u.userRepo.EmailExists(ctx, email.String())
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%w: такая почта уже существует", application.ErrAlreadyExists)
	}
	encryptedEmail, err := u.emailCrypter.Encrypt(email.String())
	if err != nil {
		return err
	}
	emailMessage := appdto.Email{
		To:    email.String(),
		Token: encryptedEmail,
	}
	go u.emailProvider.SendConfirmEmail(emailMessage)
	return nil
}

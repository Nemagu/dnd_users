package usecase

import (
	"context"

	"github.com/Nemagu/dnd/internal/application"
	"github.com/Nemagu/dnd/internal/application/service"
	"github.com/Nemagu/dnd/internal/domain"
)

type ConfirmEmailUseCase struct {
	emailCryptoService service.EmailCryptoService
	emailProvider      service.EmailProviderService
}

func NewConfirmEmail(
	emailCryptoService service.EmailCryptoService,
	emailProvider service.EmailProviderService,
) *ConfirmEmailUseCase {
	return &ConfirmEmailUseCase{
		emailCryptoService: emailCryptoService,
		emailProvider:      emailProvider,
	}
}

func (u *ConfirmEmailUseCase) Execute(
	ctx context.Context,
	input application.ConfirmEmailCommand,
) error {
	email, err := domain.NewEmail(input.Email)
	if err != nil {
		return err
	}
	encryptedEmail, err := u.emailCryptoService.Encrypt(email)
	if err != nil {
		return err
	}
	emailMessage := application.EmailMessage{
		To:   email.String(),
		Data: encryptedEmail,
	}
	err = u.emailProvider.SendConfirmEmail(ctx, emailMessage)
	if err != nil {
		return err
	}
	return nil
}

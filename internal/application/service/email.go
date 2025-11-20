package service

import (
	"context"

	"github.com/Nemagu/dnd/internal/application"
)

type EmailProviderService interface {
	SendConfirmEmail(
		ctx context.Context,
		message application.EmailMessage,
	) error
	SendChangeEmail(
		ctx context.Context,
		message application.EmailMessage,
	) error
	SendResetPasswordEmail(
		ctx context.Context,
		message application.EmailMessage,
	) error
}

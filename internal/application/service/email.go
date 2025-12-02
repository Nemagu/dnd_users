package service

import (
	appdto "github.com/Nemagu/dnd/internal/application/dto"
)

type EmailProvider interface {
	SendConfirmEmail(
		message appdto.Email,
	)
	SendChangeEmail(
		message appdto.Email,
	)
	SendResetPasswordEmail(
		message appdto.Email,
	)
}

type EmailCrypter interface {
	Encrypt(email string) (string, error)
	Decrypt(token string) (string, error)
}

type EmailValidator interface {
	Validate(email string) error
}

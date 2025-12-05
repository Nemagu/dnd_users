package usecase

import (
	"github.com/Nemagu/dnd/internal/application"
	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/google/uuid"
)

type mockPasswordComparer struct {
	IsCompare bool
}

func (p *mockPasswordComparer) Compare(password, hash string) (bool, error) {
	return p.IsCompare, nil
}

type mockPasswordValidator struct {
	IsValidate bool
}

func (p *mockPasswordValidator) Validate(password, email string) error {
	if p.IsValidate {
		return nil
	}
	return application.ErrValidation
}

type mockPasswordHasher struct{}

func (p *mockPasswordHasher) Hash(password string) (string, error) {
	return "this_is_password_hash", nil
}

type mockEmailProvider struct{}

func (p *mockEmailProvider) SendConfirmEmail(message appdto.Email) {}

func (p *mockEmailProvider) SendChangeEmail(message appdto.Email) {}

func (p *mockEmailProvider) SendResetPasswordEmail(message appdto.Email) {}

type mockEmailCrypter struct{}

func (c *mockEmailCrypter) EncryptEmail(email string) (string, error) {
	return "encrypted_email", nil
}

func (c *mockEmailCrypter) EncryptEmailUserID(email string, userID uuid.UUID) (string, error) {
	return "encrypted_email", nil
}

type mockEmailValidator struct {
	IsValid bool
}

func (v *mockEmailValidator) Validate(email string) error {
	if v.IsValid {
		return nil
	} else {
		return application.ErrValidation
	}
}

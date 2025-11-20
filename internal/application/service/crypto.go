package service

import "github.com/Nemagu/dnd/internal/domain"

type PasswordCryptoService interface {
	Encrypt(password string) (string, error)
	ComparePassword(password, passwordHash string) bool
}

type EmailCryptoService interface {
	Encrypt(email domain.Email) (string, error)
	Decrypt(email domain.Email) (string, error)
}

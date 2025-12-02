package appdto

import (
	"github.com/google/uuid"
)

type User struct {
	UserID       uuid.UUID
	Email        string
	State        string
	Status       string
	PasswordHash string
	Version      uint64
}

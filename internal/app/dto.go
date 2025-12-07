package app

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID
	Email        string
	State        string
	Status       string
	PasswordHash string
	Version      uint
}

type EmailCode struct {
	To   string
	Code string
}

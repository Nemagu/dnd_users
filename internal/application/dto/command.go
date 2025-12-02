package appdto

import "github.com/google/uuid"

type ConfirmEmailCommand struct {
	Email string
}

type RegisterUserCommand struct {
	Token    string
	Password string
}

type ConfirmNewEmailCommand struct {
	UserID   uuid.UUID
	Email    string
	Password string
}

type ResetPasswordCommand struct {
	UserID          uuid.UUID
	CurrentPassword string
}

type AuthenticateCommand struct {
	Email    string
	Password string
}

type ChangeUserCommand struct {
	InitiatorID uuid.UUID
	UserID      uuid.UUID
	Email       string
	State       string
	Status      string
	Password    string
}

type UpdateUserCommand struct {
	InitiatorID uuid.UUID
	UserID      uuid.UUID
	Status      string
	State       string
}

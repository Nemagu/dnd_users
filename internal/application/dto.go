package application

import "github.com/google/uuid"

type ConfirmEmailCommand struct {
	Email string
}

type CreateUserCommand struct {
	Token    string
	Password string
}

type UpdateUserCommand struct {
	InitiatorID uuid.UUID
	UserID      uuid.UUID
	Email       string
	Password    string
	Status      string
	State       string
}

type UserQuery struct {
	InitiatorID uuid.UUID
	UserID      uuid.UUID
}

type UsersQuery struct {
	InitiatorID    uuid.UUID
	SearchByEmail  string
	FilterByState  string
	FilterByStatus string
}

type EmailMessage struct {
	To   string
	Data string
}

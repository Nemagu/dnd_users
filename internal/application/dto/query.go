package appdto

import "github.com/google/uuid"

type UserQuery struct {
	InitiatorID uuid.UUID
	UserID      uuid.UUID
}

type UsersQuery struct {
	InitiatorID    uuid.UUID
	SearchByEmail  string
	FilterByState  []string
	FilterByStatus []string
	Limit          int
	Offset         int
}

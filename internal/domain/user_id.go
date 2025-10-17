package domain

import "github.com/google/uuid"

type UserID struct {
	id uuid.UUID
}

func NewUserID(id uuid.UUID) (UserID, error) {
	return UserID{id: id}, nil
}

func (ui UserID) ID() uuid.UUID {
	return ui.id
}

func (ui UserID) String() string {
	return ui.id.String()
}

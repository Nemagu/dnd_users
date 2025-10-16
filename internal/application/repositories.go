package application

import "github.com/Nemagu/dnd/internal/domain"

type UserRepository interface {
	GetAll() []*domain.User
	GetOfID(id domain.UserID) (*domain.User, error)
	Save(user *domain.User) error
}

type EventRepository interface {
	EmailChanged(id domain.UserID, email domain.Email) error
	StateChanged(id domain.UserID, state domain.UserState) error
	StatusChanged(id domain.UserID, status domain.UserStatus) error
}

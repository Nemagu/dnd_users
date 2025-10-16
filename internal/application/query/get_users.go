package query

import (
	"github.com/Nemagu/dnd/internal/application"
	"github.com/Nemagu/dnd/internal/domain"
)

type UsersQuery struct {
	repository application.UserRepository
}

func NewUsersQuery(repository application.UserRepository) *UsersQuery {
	return &UsersQuery{repository: repository}
}

func (q *UsersQuery) Execute(initiatorID domain.UserID) ([]*domain.User, error) {
	initiator, err := q.repository.GetOfID(initiatorID)
	if err != nil {
		return make([]*domain.User, 0), err
	}
	if !initiator.Status().IsAdmin() {
		return make([]*domain.User, 0), application.NoAccessError("")
	}
	return q.repository.GetAll(), nil
}

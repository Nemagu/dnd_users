package query

import (
	"github.com/Nemagu/dnd/internal/application"
	"github.com/Nemagu/dnd/internal/domain"
)

type UserQuery struct {
	repository application.UserRepository
}

func NewUserQuery(repository application.UserRepository) *UserQuery {
	return &UserQuery{repository: repository}
}

func (q *UserQuery) Execute(initiatorID domain.UserID, userID domain.UserID) (*domain.User, error) {
	initiator, err := q.repository.GetOfID(initiatorID)
	if err != nil {
		return initiator, err
	}
	if initiatorID == userID {
		return initiator, nil
	}
	if !initiator.Status().IsAdmin() {
		return &domain.User{}, application.NoAccessError("")
	}
	return q.repository.GetOfID(userID)
}

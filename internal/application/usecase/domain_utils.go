package usecase

import (
	"errors"
	"fmt"

	"github.com/Nemagu/dnd/internal/application"
	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/Nemagu/dnd/internal/domain"
	"github.com/Nemagu/dnd/internal/domain/duser"
	"github.com/google/uuid"
)

func toModifyAppUser(user *duser.User) *appdto.User {
	return &appdto.User{
		UserID:       user.ID(),
		Email:        user.Email().String(),
		State:        user.State().String(),
		Status:       user.Status().String(),
		PasswordHash: user.PasswordHash(),
		Version:      user.ModifyVersion(),
	}
}

func toAppUser(user *duser.User) *appdto.User {
	return &appdto.User{
		UserID:       user.ID(),
		Email:        user.Email().String(),
		State:        user.State().String(),
		Status:       user.Status().String(),
		PasswordHash: user.PasswordHash(),
		Version:      user.Version(),
	}
}

func restoreDomainUser(user *appdto.User) (*duser.User, error) {
	e, err := toDomainEmail(user.Email)
	if err != nil {
		return nil, err
	}

	se, err := toDomainState(user.State)
	if err != nil {
		return nil, err
	}

	su, err := toDomainStatus(user.Status)
	if err != nil {
		return nil, err
	}

	u, err := duser.Restore(
		user.UserID, e, se, su, user.PasswordHash, user.Version,
	)
	if err != nil {
		return u, handleError(err)
	}

	return u, err
}

func newDomainUser(userID uuid.UUID, email, passwordHash string) (*duser.User, error) {
	e, err := toDomainEmail(email)
	if err != nil {
		return nil, err
	}

	u, err := duser.New(
		userID, e, passwordHash,
	)
	if err != nil {
		return u, handleError(err)
	}

	return u, err
}

func toDomainEmail(email string) (domain.Email, error) {
	e, err := domain.NewEmail(email)
	if err != nil {
		return e, handleError(err)
	}
	return e, err
}

func toDomainState(state string) (duser.State, error) {
	s, err := duser.StateFromString(state)
	if err != nil {
		return s, handleError(err)
	}
	return s, err
}

func toDomainStatus(status string) (duser.Status, error) {
	s, err := duser.StatusFromString(status)
	if err != nil {
		return s, handleError(err)
	}
	return s, err
}

func handleError(err error) error {
	switch {
	case errors.Is(err, domain.ErrIdempotent) || errors.Is(
		err, domain.ErrValidation,
	) || errors.Is(err, domain.ErrInvalidData):
		return fmt.Errorf("%w: %s", application.ErrValidation, err)
	case errors.Is(err, domain.ErrNotAllowed):
		return fmt.Errorf("%w: %s", application.ErrNotAllowed, err)
	default:
		return fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
}

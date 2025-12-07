package app

import (
	"errors"
	"fmt"

	"github.com/Nemagu/dnd_users/internal/domain"
)

func currentUser(u *domain.User) (*User, error) {
	if u == nil {
		return nil, fmt.Errorf(
			"%w: получен nil для преобразования из доменного пользователя в пользователя из приложения",
			ErrInternal,
		)
	}
	return &User{
		ID:           u.ID(),
		Email:        u.Email(),
		State:        u.State().String(),
		Status:       u.Status().String(),
		PasswordHash: u.PasswordHash(),
		Version:      u.Version(),
	}, nil
}

func modifiedUser(u *domain.User) (*User, error) {
	if u == nil {
		return nil, fmt.Errorf(
			"%w: получен nil для преобразования из доменного пользователя в пользователя из приложения",
			ErrInternal,
		)
	}
	return &User{
		ID:           u.ID(),
		Email:        u.Email(),
		State:        u.State().String(),
		Status:       u.Status().String(),
		PasswordHash: u.PasswordHash(),
		Version:      u.ModifiedVersion(),
	}, nil
}

func domainUser(u *User) (*domain.User, error) {
	if u == nil {
		return nil, fmt.Errorf(
			"%w: получен nil для преобразования из пользователя из приложения в доменного пользователя",
			ErrInternal,
		)
	}

	dState, err := domainState(u.State)
	if err != nil {
		return nil, err
	}

	dStatus, err := domainStatus(u.Status)
	if err != nil {
		return nil, err
	}

	user, err := domain.RestoreUser(u.ID, u.Email, u.PasswordHash, dState, dStatus, u.Version)
	if err != nil {
		return nil, handleDomainError(err)
	}

	return user, nil
}

func domainState(s string) (domain.State, error) {
	state, err := domain.NewState(s)
	if err != nil {
		return domain.NilState, handleDomainError(err)
	}
	return state, nil
}

func domainStatus(s string) (domain.Status, error) {
	status, err := domain.NewStatus(s)
	if err != nil {
		return domain.NilStatus, handleDomainError(err)
	}
	return status, nil
}

func handleDomainError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidData):
		return fmt.Errorf("%w: %s", ErrInvalidData, err)
	case errors.Is(err, domain.ErrUserNotActive):
		return fmt.Errorf("%w: %s", ErrUserNotActive, err)
	case errors.Is(err, domain.ErrIdempotent):
		return fmt.Errorf("%w: %s", ErrIdempotent, err)
	default:
		return fmt.Errorf("%w: %s", ErrInternal, err)
	}
}

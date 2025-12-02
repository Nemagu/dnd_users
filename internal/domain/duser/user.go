package duser

import (
	"fmt"

	"github.com/Nemagu/dnd/internal/domain"
	"github.com/google/uuid"
)

type User struct {
	userID       uuid.UUID
	email        domain.Email
	state        State
	status       Status
	passwordHash string
	version      uint64
}

func New(
	userID uuid.UUID,
	email domain.Email,
	passwordHash string,
) (*User, error) {
	return &User{
		userID:       userID,
		email:        email,
		state:        NewActiveState(),
		status:       NewOrdinaryStatus(),
		passwordHash: passwordHash,
		version:      0,
	}, nil
}

func Restore(
	userID uuid.UUID,
	email domain.Email,
	state State,
	status Status,
	passwordHash string,
	version uint64,
) (*User, error) {
	if version < 1 {
		return nil, fmt.Errorf(
			"%w: при восстановлении пользователя версия не может быть меньше либо ровна 0",
			domain.ErrInternal,
		)
	}
	return &User{
		userID:       userID,
		email:        email,
		state:        state,
		status:       status,
		passwordHash: passwordHash,
		version:      version,
	}, nil
}

func (u *User) ID() uuid.UUID {
	return u.userID
}

func (u *User) Email() domain.Email {
	return u.email
}

func (u *User) State() State {
	return u.state
}

func (u *User) Status() Status {
	return u.status
}

func (u *User) PasswordHash() string {
	return u.passwordHash
}

func (u *User) Version() uint64 {
	return u.version
}

func (u *User) ModifyVersion() uint64 {
	return u.version + 1
}

func (u *User) ChangePassword(passwordHash string) error {
	if err := u.assertIsNotActive(); err != nil {
		return err
	}
	u.passwordHash = passwordHash
	return nil
}

func (u *User) ChangeEmail(email domain.Email) error {
	if err := u.assertIsNotActive(); err != nil {
		return err
	}
	u.email = email
	return nil
}

func (u *User) ChangeStatus(status Status) error {
	if err := u.assertIsNotActive(); err != nil {
		return err
	}
	if u.status == status {
		return fmt.Errorf(
			"%w: статус пользователя не изменяется",
			domain.ErrIdempotent,
		)
	}
	u.status = status
	return nil
}

func (u *User) ChangeState(state State) error {
	if u.state == state {
		return fmt.Errorf(
			"%w: состояние пользователя не изменяется",
			domain.ErrIdempotent,
		)
	}
	u.state = state
	return nil
}

func (u *User) assertIsNotActive() (err error) {
	if !u.state.IsActive() {
		err = fmt.Errorf("%w: пользователь не является активным", domain.ErrNotAllowed)
	}
	return
}

func (u *User) String() string {
	return fmt.Sprintf("id: %s\nemail: %s", u.userID, u.email)
}

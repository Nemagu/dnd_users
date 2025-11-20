package duser

import (
	"fmt"

	"github.com/Nemagu/dnd/internal/domain"
	"github.com/google/uuid"
)

type User struct {
	userID       uuid.UUID
	email        domain.Email
	state        UserState
	status       UserStatus
	passwordHash string
	version      uint64
}

func New(
	userID uuid.UUID,
	email domain.Email,
	state UserState,
	status UserStatus,
	passwordHash string,
) (*User, error) {
	return &User{
		userID:       userID,
		email:        email,
		state:        state,
		status:       status,
		passwordHash: passwordHash,
		version:      0,
	}, nil
}

func Restore(
	userID uuid.UUID,
	email domain.Email,
	state UserState,
	status UserStatus,
	passwordHash string,
	version uint64,
) (*User, error) {
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

func (u *User) State() UserState {
	return u.state
}

func (u *User) Status() UserStatus {
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

func (u *User) AppointAdmin() error {
	if err := u.assertIsNotActive(); err != nil {
		return err
	}
	if u.status.IsAdmin() {
		return fmt.Errorf(
			"%w: пользователь уже является администратором",
			domain.ErrIdempotent,
		)
	}
	u.status = NewAdminStatus()
	return nil
}

func (u *User) AppointOrdinary() error {
	if err := u.assertIsNotActive(); err != nil {
		return err
	}
	if u.status.IsOrdinary() {
		return fmt.Errorf(
			"%w: пользователь уже является обычным",
			domain.ErrIdempotent,
		)
	}
	u.status = NewAdminStatus()
	return nil
}

func (u *User) Activate() error {
	if u.state.IsActive() {
		return fmt.Errorf(
			"%w: пользователь уже активен",
			domain.ErrIdempotent,
		)
	}
	u.state = NewActiveState()
	return nil
}

func (u *User) Freeze() error {
	if u.state.IsFrozen() {
		return fmt.Errorf(
			"%w: пользователь уже заморожен",
			domain.ErrIdempotent,
		)
	}
	u.state = NewFrozenState()
	return nil
}

func (u *User) Delete() error {
	if u.state.IsDeleted() {
		return fmt.Errorf("%w: пользователь уже удален", domain.ErrIdempotent)
	}
	u.state = NewDeletedState()
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

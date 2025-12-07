package domain

import (
	"fmt"

	"github.com/google/uuid"
)

type User struct {
	id           uuid.UUID
	email        string
	state        State
	status       Status
	passwordHash string
	version      uint
}

func NewUser(id uuid.UUID, email, passwordHash string) (*User, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("%w: id пользователя не может быть пустым", ErrInvalidData)
	}
	if email == "" {
		return nil, fmt.Errorf("%w: email пользователя не может быть пустым", ErrInvalidData)
	}
	if passwordHash == "" {
		return nil, fmt.Errorf("%w: пароль пользователя не может быть пустым", ErrInvalidData)
	}
	return &User{
		id:           id,
		email:        email,
		state:        newActiveState(),
		status:       newUserStatus(),
		passwordHash: passwordHash,
		version:      0,
	}, nil
}

func RestoreUser(
	id uuid.UUID,
	email, passwordHash string,
	state State,
	status Status,
	version uint,
) (*User, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("%w: id пользователя не может быть пустым", ErrInvalidData)
	}
	if email == "" {
		return nil, fmt.Errorf("%w: email пользователя не может быть пустым", ErrInvalidData)
	}
	if state == nilState {
		return nil, fmt.Errorf("%w: состояние пользователя не может быть пустым", ErrInvalidData)
	}
	if status == nilStatus {
		return nil, fmt.Errorf("%w: статус пользователя не может быть пустым", ErrInvalidData)
	}
	if passwordHash == "" {
		return nil, fmt.Errorf("%w: пароль пользователя не может быть пустым", ErrInvalidData)
	}
	if version == 0 {
		return nil, fmt.Errorf("%w: версия пользователя не может быть равна 0", ErrInvalidData)
	}
	return &User{
		id:           id,
		email:        email,
		state:        state,
		status:       status,
		passwordHash: passwordHash,
		version:      version,
	}, nil
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) Email() string {
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

func (u *User) Version() uint {
	return u.version
}

func (u *User) ModifiedVersion() uint {
	return u.version + 1
}

func (u *User) NewEmail(email string) error {
	if err := u.checkState(); err != nil {
		return err
	}
	if email == "" {
		return fmt.Errorf("%w: email пользователя не может быть пустым", ErrInvalidData)
	}
	if u.email == email {
		return fmt.Errorf("%w: email пользователя уже %s", ErrIdempotent, email)
	}
	u.email = email
	return nil
}

func (u *User) NewState(state State) error {
	if err := u.checkState(); err != nil {
		return err
	}
	if state == nilState {
		return fmt.Errorf("%w: состояние пользователя не может быть пустым", ErrInvalidData)
	}
	if u.state == state {
		return fmt.Errorf("%w: состояние пользователя уже %s", ErrIdempotent, state)
	}
	u.state = state
	return nil
}

func (u *User) NewStatus(status Status) error {
	if err := u.checkState(); err != nil {
		return err
	}
	if status == nilStatus {
		return fmt.Errorf("%w: статус пользователя не может быть пустым", ErrInvalidData)
	}
	if u.status == status {
		return fmt.Errorf("%w: статус пользователя уже %s", ErrIdempotent, status)
	}
	u.status = status
	return nil
}

func (u *User) NewPasswordHash(passwordHash string) error {
	if err := u.checkState(); err != nil {
		return err
	}
	if passwordHash == "" {
		return fmt.Errorf("%w: пароль пользователя не может быть пустым", ErrInvalidData)
	}
	u.passwordHash = passwordHash
	return nil
}

func (u *User) checkState() error {
	if !u.state.IsActive() {
		return fmt.Errorf("%w: id пользователя %s", ErrUserNotActive, u.id)
	}
	return nil
}

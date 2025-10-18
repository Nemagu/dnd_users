package domain

import "fmt"

type User struct {
	userID       UserID
	username     Username
	email        Email
	state        UserState
	status       UserStatus
	passwordHash PasswordHash
	person       Person
}

func NewUser(
	userID UserID,
	username Username,
	email Email,
	state UserState,
	status UserStatus,
	passwordHash PasswordHash,
	person Person,
) (User, error) {
	return User{
		userID:       userID,
		username:     username,
		email:        email,
		state:        state,
		status:       status,
		passwordHash: passwordHash,
		person:       person,
	}, nil
}

func (u *User) UserID() UserID {
	return u.userID
}

func (u *User) Username() Username {
	return u.username
}

func (u *User) Email() Email {
	return u.email
}

func (u *User) State() UserState {
	return u.state
}

func (u *User) Status() UserStatus {
	return u.status
}

func (u *User) PasswordHash() PasswordHash {
	return u.passwordHash
}

func (u *User) Person() Person {
	return u.person
}

func (u *User) ChangePassword(passwordHash PasswordHash) error {
	if err := u.assertIsNotActive(); err != nil {
		return err
	}
	u.passwordHash = passwordHash
	return nil
}

func (u *User) ChangeEmail(email Email) error {
	if err := u.assertIsNotActive(); err != nil {
		return err
	}
	u.email = email
	return nil
}

func (u *User) ChangeUsername(username Username) error {
	if err := u.assertIsNotActive(); err != nil {
		return err
	}
	u.username = username
	return nil
}

func (u *User) ChangePerson(person Person) error {
	if err := u.assertIsNotActive(); err != nil {
		return err
	}
	u.person = person
	return nil
}

func (u *User) AppointAdmin() error {
	if err := u.assertIsNotActive(); err != nil {
		return err
	}
	if u.status.IsAdmin() {
		return IdempotentError("пользователь уже является администратором")
	}
	u.status = NewAdminUserStatus()
	return nil
}

func (u *User) AppointOrdinary() error {
	if err := u.assertIsNotActive(); err != nil {
		return err
	}
	if u.status.IsOrdinary() {
		return IdempotentError("пользователь уже является обычным")
	}
	u.status = NewAdminUserStatus()
	return nil
}

func (u *User) Activate() error {
	if u.state.IsActive() {
		return IdempotentError("пользователь уже активен")
	}
	u.state = NewActiveUserState()
	return nil
}

func (u *User) Freeze() error {
	if u.state.IsFrozen() {
		return IdempotentError("пользователь уже заморожен")
	}
	u.state = NewFrozenUserState()
	return nil
}

func (u *User) Delete() error {
	if u.state.IsDeleted() {
		return IdempotentError("пользователь уже удален")
	}
	u.state = NewDeletedUserState()
	return nil
}

func (u *User) assertIsNotActive() error {
	var err error
	if !u.state.IsActive() {
		err = PolicyError("пользователь не является активным")
	}
	return err
}

func (u *User) String() string {
	return fmt.Sprintf("id: %s\nusername: %s\nemail: %s", u.userID, u.username, u.email)
}

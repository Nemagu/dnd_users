package domain

type User struct {
	userID       UserID
	username     Username
	email        Email
	state        UserState
	passwordHash PasswordHash
	person       Person
}

func NewUser(
	userID UserID,
	username Username,
	email Email,
	state UserState,
	passwordHash PasswordHash,
	person Person,
) (User, error) {
	return User{
		userID:       userID,
		username:     username,
		email:        email,
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

func (u *User) assertIsNotActive() error {
	var err error
	if !u.state.IsActive() {
		err = &DomainError{Message: "Пользователь не является активным"}
	}
	return err
}

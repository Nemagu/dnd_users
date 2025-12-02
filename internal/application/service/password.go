package service

type PasswordHasher interface {
	Hash(password string) (string, error)
	ComparePassword(password, passwordHash string) (bool, error)
}

type PasswordValidator interface {
	Validate(password, email string) error
}

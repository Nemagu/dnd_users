package domain

type PasswordHash struct {
	passwordHash string
}

func NewPasswordHash(passwordHash string) (PasswordHash, error) {
	return PasswordHash{passwordHash: passwordHash}, nil
}

func (ph PasswordHash) PasswordHash() string {
	return ph.passwordHash
}

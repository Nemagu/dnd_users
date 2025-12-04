package password

import (
	"errors"
	"fmt"

	"github.com/Nemagu/dnd/internal/application"
	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordHasher struct {
	cost int
}

func MustNewBcryptPasswordHasher(cost int) *BcryptPasswordHasher {
	if cost < bcrypt.DefaultCost || cost > bcrypt.MaxCost {
		panic(fmt.Sprintf(
			"стоимость bcrypt не попадает в рамки допустимых значений от %d до %d",
			bcrypt.DefaultCost,
			bcrypt.MaxCost,
		))
	} else {
		return &BcryptPasswordHasher{cost: cost}
	}
}

func (s *BcryptPasswordHasher) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	return string(hashedPassword), nil
}

func (s *BcryptPasswordHasher) Compare(password, passwordHash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, fmt.Errorf("%w: %s", application.ErrInternal, err)
	}
	return true, nil
}

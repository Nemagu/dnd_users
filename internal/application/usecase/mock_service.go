package usecase

import "github.com/Nemagu/dnd/internal/application"

type mockPasswordComparer struct {
	IsCompare bool
}

func (p *mockPasswordComparer) Compare(password, hash string) (bool, error) {
	return p.IsCompare, nil
}

type mockPasswordValidator struct {
	IsValidate bool
}

func (p *mockPasswordValidator) Validate(password, email string) error {
	if p.IsValidate {
		return nil
	}
	return application.ErrValidation
}

type mockPasswordHasher struct{}

func (p *mockPasswordHasher) Hash(password string) (string, error) {
	return "this_is_password_hash", nil
}

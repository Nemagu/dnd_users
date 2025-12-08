package app

import (
	"fmt"
	"slices"
)

type mockEmailValidator struct {
	InvalidEmails []string
}

func (m *mockEmailValidator) Validate(email string) error {
	if !slices.Contains(m.InvalidEmails, email) {
		return nil
	} else {
		return fmt.Errorf("%w: invalid email %s", ErrInvalidData, email)
	}
}

type mockCodeGenerator struct {
	Code string
}

func (m *mockCodeGenerator) Generate() string {
	if m.Code == "" {
		return "123456"
	} else {
		return m.Code
	}
}

type mockPasswordValidator struct {
	InvalidPasswords []string
}

func (m *mockPasswordValidator) Validate(password, email string) error {
	if slices.Contains(m.InvalidPasswords, password) {
		return fmt.Errorf("%w: password is invalid", ErrInvalidData)
	} else {
		return nil
	}
}

type mockPasswordHasher struct {
	Err error
}

func (m *mockPasswordHasher) Hash(password string) (string, error) {
	return password, m.Err
}

type mockPasswordComparer struct {
	ValidPassword []string
	Err           error
}

func (m *mockPasswordComparer) Compare(password, hash string) (bool, error) {
	return slices.Contains(m.ValidPassword, password), m.Err
}

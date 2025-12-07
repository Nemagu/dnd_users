package app

import (
	"fmt"
	"slices"
)

type mockEmailValidator struct {
	ValidEmail []string
}

func (m *mockEmailValidator) Validate(email string) error {
	if slices.Contains(m.ValidEmail, email) {
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

package app

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"testing"
)

type mockConfirmEmailRepository struct {
	NotExistsEmails []string
}

func (m *mockConfirmEmailRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	return !slices.Contains(m.NotExistsEmails, email), nil
}

type mockConfirmEmailCodeStore struct {
	Err error
}

func (m *mockConfirmEmailCodeStore) SetCEC(ctx context.Context, key, value string) error {
	return m.Err
}

type mockConfirmEmailProvider struct{}

func (m *mockConfirmEmailProvider) SendConfirmationEmail(data EmailCode) {}

func TestConfirmEmailUseCase_Execute(t *testing.T) {
	validEmail := "test@mail.com"
	cases := []struct {
		TestName string
		Expected error
		UC       *ConfirmEmailUseCase
		Command  *ConfirmEmailCommand
	}{
		{
			TestName: "test_confirm_email_use_case_ok",
			Expected: nil,
			UC: MustConfirmEmailUseCase(
				&mockConfirmEmailRepository{NotExistsEmails: []string{validEmail}},
				&mockConfirmEmailCodeStore{},
				&mockEmailValidator{},
				&mockConfirmEmailProvider{},
				&mockCodeGenerator{},
			),
			Command: &ConfirmEmailCommand{Email: validEmail},
		},
		{
			TestName: "test_confirm_email_use_case_email_exists",
			Expected: ErrAlreadyExists,
			UC: MustConfirmEmailUseCase(
				&mockConfirmEmailRepository{},
				&mockConfirmEmailCodeStore{},
				&mockEmailValidator{},
				&mockConfirmEmailProvider{},
				&mockCodeGenerator{},
			),
			Command: &ConfirmEmailCommand{Email: validEmail},
		},
		{
			TestName: "test_confirm_email_use_case_setting_code_error",
			Expected: ErrInternal,
			UC: MustConfirmEmailUseCase(
				&mockConfirmEmailRepository{NotExistsEmails: []string{validEmail}},
				&mockConfirmEmailCodeStore{Err: fmt.Errorf("%w: internal error", ErrInternal)},
				&mockEmailValidator{},
				&mockConfirmEmailProvider{},
				&mockCodeGenerator{},
			),
			Command: &ConfirmEmailCommand{Email: validEmail},
		},
		{
			TestName: "test_confirm_email_use_case_invalid_email",
			Expected: ErrInvalidData,
			UC: MustConfirmEmailUseCase(
				&mockConfirmEmailRepository{NotExistsEmails: []string{validEmail}},
				&mockConfirmEmailCodeStore{},
				&mockEmailValidator{NotValidEmails: []string{validEmail}},
				&mockConfirmEmailProvider{},
				&mockCodeGenerator{},
			),
			Command: &ConfirmEmailCommand{Email: validEmail},
		},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			err := c.UC.Execute(context.Background(), c.Command)
			if c.Expected == nil {
				if err != nil {
					t.Errorf("expected nil, but got %v", err)
				}
			}
			if !errors.Is(err, c.Expected) {
				t.Errorf("expected %T, but got %v", c.Expected, err)
			}
		})
	}
}

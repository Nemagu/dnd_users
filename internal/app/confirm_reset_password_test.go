package app

import (
	"context"
	"errors"
	"testing"
)

type mockConfirmResetPasswordRepository struct {
	Err    error
	Exists bool
}

func (m *mockConfirmResetPasswordRepository) EmailExists(
	ctx context.Context,
	email string,
) (bool, error) {
	return m.Exists, m.Err
}

type mockConfirmResetPasswordCodeStore struct {
	Err error
}

func (m *mockConfirmResetPasswordCodeStore) SetResetPassword(
	ctx context.Context,
	key, value string,
) error {
	return m.Err
}

type mockConfirmResetPasswordProvider struct{}

func (m *mockConfirmResetPasswordProvider) SendResetPasswordEmail(data EmailCode) {}

func TestConfirmResetPasswordUseCase_Execute(t *testing.T) {
	validEmail := "test@mail.com"
	cases := []struct {
		TestName string
		Expected error
		UC       *ConfirmResetPasswordUseCase
		Command  *ConfirmResetPasswordCommand
	}{
		{
			TestName: "test_confirm_reset_password_use_case_ok",
			Expected: nil,
			UC: MustConfirmResetPasswordUseCase(
				&mockConfirmResetPasswordRepository{Exists: true},
				&mockConfirmResetPasswordCodeStore{},
				&mockConfirmResetPasswordProvider{},
				&mockCodeGenerator{},
			),
			Command: &ConfirmResetPasswordCommand{Email: validEmail},
		},
		{
			TestName: "test_confirm_reset_password_use_case_email_not_exists",
			Expected: ErrNotFound,
			UC: MustConfirmResetPasswordUseCase(
				&mockConfirmResetPasswordRepository{Err: ErrNotFound},
				&mockConfirmResetPasswordCodeStore{},
				&mockConfirmResetPasswordProvider{},
				&mockCodeGenerator{},
			),
			Command: &ConfirmResetPasswordCommand{Email: validEmail},
		},
		{
			TestName: "test_confirm_reset_password_use_case_code_store_error",
			Expected: ErrInternal,
			UC: MustConfirmResetPasswordUseCase(
				&mockConfirmResetPasswordRepository{Exists: true},
				&mockConfirmResetPasswordCodeStore{Err: ErrInternal},
				&mockConfirmResetPasswordProvider{},
				&mockCodeGenerator{},
			),
			Command: &ConfirmResetPasswordCommand{Email: validEmail},
		},
		{
			TestName: "test_confirm_reset_password_use_case_repo_error",
			Expected: ErrInternal,
			UC: MustConfirmResetPasswordUseCase(
				&mockConfirmResetPasswordRepository{Exists: true, Err: ErrInternal},
				&mockConfirmResetPasswordCodeStore{},
				&mockConfirmResetPasswordProvider{},
				&mockCodeGenerator{},
			),
			Command: &ConfirmResetPasswordCommand{Email: validEmail},
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

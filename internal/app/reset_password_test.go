package app

import (
	"context"
	"errors"
	"testing"

	"github.com/Nemagu/dnd_users/internal/domain"
	"github.com/google/uuid"
)

type mockResetPasswordRepository struct {
	User       *User
	ErrByEmail error
	ErrSave    error
}

func (m *mockResetPasswordRepository) ByEmail(ctx context.Context, email string) (*User, error) {
	return m.User, m.ErrByEmail
}

func (m *mockResetPasswordRepository) Save(ctx context.Context, user *User) error {
	return m.ErrSave
}

type mockResetPasswordCodeStore struct {
	ErrGet error
	ErrDel error
	Code   string
}

func (m *mockResetPasswordCodeStore) GetResetPassword(
	ctx context.Context,
	key string,
) (string, error) {
	return m.Code, m.ErrGet
}

func (m *mockResetPasswordCodeStore) DelResetPassword(ctx context.Context, key string) error {
	return m.ErrDel
}

func TestResetPasswordUseCase_Execute(t *testing.T) {
	activeUser := &User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		State:        domain.ACTIVE,
		Status:       domain.USER,
		PasswordHash: "test",
		Version:      1,
	}
	notActiveUser := &User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		State:        domain.FROZEN,
		Status:       domain.USER,
		PasswordHash: "test",
		Version:      1,
	}
	newPassword := "new_password"
	validCode := "valid_code"
	cases := []struct {
		TestName string
		Expected error
		UC       *ResetPasswordUseCase
		Command  *ResetPasswordCommand
	}{
		{
			TestName: "test_reset_password_use_case_ok",
			Expected: nil,
			UC: MustResetPasswordUseCase(
				&mockResetPasswordRepository{User: activeUser},
				&mockResetPasswordCodeStore{Code: validCode},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &ResetPasswordCommand{
				NewPassword: newPassword,
				Code:        validCode,
				Email:       activeUser.Email,
			},
		},
		{
			TestName: "test_reset_password_use_case_email_not_exists",
			Expected: ErrNotFound,
			UC: MustResetPasswordUseCase(
				&mockResetPasswordRepository{User: activeUser, ErrByEmail: ErrNotFound},
				&mockResetPasswordCodeStore{Code: validCode},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &ResetPasswordCommand{
				NewPassword: newPassword,
				Code:        validCode,
				Email:       activeUser.Email,
			},
		},
		{
			TestName: "test_reset_password_use_case_saving_error",
			Expected: ErrInternal,
			UC: MustResetPasswordUseCase(
				&mockResetPasswordRepository{User: activeUser, ErrSave: ErrInternal},
				&mockResetPasswordCodeStore{Code: validCode},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &ResetPasswordCommand{
				NewPassword: newPassword,
				Code:        validCode,
				Email:       activeUser.Email,
			},
		},
		{
			TestName: "test_reset_password_use_case_invalid_code",
			Expected: ErrInvalidData,
			UC: MustResetPasswordUseCase(
				&mockResetPasswordRepository{User: activeUser},
				&mockResetPasswordCodeStore{Code: "validCode"},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &ResetPasswordCommand{
				NewPassword: newPassword,
				Code:        validCode,
				Email:       activeUser.Email,
			},
		},
		{
			TestName: "test_reset_password_use_case_code_get_error",
			Expected: ErrInternal,
			UC: MustResetPasswordUseCase(
				&mockResetPasswordRepository{User: activeUser},
				&mockResetPasswordCodeStore{Code: validCode, ErrGet: ErrInternal},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &ResetPasswordCommand{
				NewPassword: newPassword,
				Code:        validCode,
				Email:       activeUser.Email,
			},
		},
		{
			TestName: "test_reset_password_use_case_code_delete_error",
			Expected: ErrInternal,
			UC: MustResetPasswordUseCase(
				&mockResetPasswordRepository{User: activeUser},
				&mockResetPasswordCodeStore{Code: validCode, ErrDel: ErrInternal},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &ResetPasswordCommand{
				NewPassword: newPassword,
				Code:        validCode,
				Email:       activeUser.Email,
			},
		},
		{
			TestName: "test_reset_password_use_case_not_active_user",
			Expected: ErrUserNotActive,
			UC: MustResetPasswordUseCase(
				&mockResetPasswordRepository{User: notActiveUser},
				&mockResetPasswordCodeStore{Code: validCode},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &ResetPasswordCommand{
				NewPassword: newPassword,
				Code:        validCode,
				Email:       notActiveUser.Email,
			},
		},
		{
			TestName: "test_reset_password_use_case_invalid_password",
			Expected: ErrInvalidData,
			UC: MustResetPasswordUseCase(
				&mockResetPasswordRepository{User: activeUser},
				&mockResetPasswordCodeStore{Code: validCode},
				&mockPasswordValidator{InvalidPasswords: []string{newPassword}},
				&mockPasswordHasher{},
			),
			Command: &ResetPasswordCommand{
				NewPassword: newPassword,
				Code:        validCode,
				Email:       activeUser.Email,
			},
		},
		{
			TestName: "test_reset_password_use_case_hashing_error",
			Expected: ErrInternal,
			UC: MustResetPasswordUseCase(
				&mockResetPasswordRepository{User: activeUser},
				&mockResetPasswordCodeStore{Code: validCode},
				&mockPasswordValidator{},
				&mockPasswordHasher{Err: ErrInternal},
			),
			Command: &ResetPasswordCommand{
				NewPassword: newPassword,
				Code:        validCode,
				Email:       activeUser.Email,
			},
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

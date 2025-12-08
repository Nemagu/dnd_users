package app

import (
	"context"
	"errors"
	"testing"

	"github.com/Nemagu/dnd_users/internal/domain"
	"github.com/google/uuid"
)

type mockNewPasswordRepository struct {
	User    *User
	ErrByID error
	ErrSave error
}

func (m *mockNewPasswordRepository) ByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return m.User, m.ErrByID
}

func (m *mockNewPasswordRepository) Save(ctx context.Context, user *User) error {
	return m.ErrSave
}

type mockNewPasswordCodeStore struct {
	Code   string
	ErrGet error
	ErrDel error
}

func (m *mockNewPasswordCodeStore) GetNewPassword(ctx context.Context, key string) (string, error) {
	return m.Code, m.ErrGet
}

func (m *mockNewPasswordCodeStore) DelNewPassword(ctx context.Context, key string) error {
	return m.ErrDel
}

func TestNewPasswordUseCase_Execute(t *testing.T) {
	activeUser := &User{
		ID:           uuid.New(),
		Email:        "test@mail.com",
		State:        domain.ACTIVE,
		Status:       domain.USER,
		PasswordHash: "test",
		Version:      1,
	}
	notActiveUser := &User{
		ID:           uuid.New(),
		Email:        "test@mail.com",
		State:        domain.FROZEN,
		Status:       domain.USER,
		PasswordHash: "test",
		Version:      1,
	}
	validCode := "valid_code"
	invalidCode := "invalid_code"
	cases := []struct {
		TestName string
		Expected error
		UC       *NewPasswordUseCase
		Command  *NewPasswordCommand
	}{
		{
			TestName: "test_new_password_use_case_ok",
			Expected: nil,
			UC: MustNewPasswordUseCase(
				&mockNewPasswordRepository{User: activeUser},
				&mockNewPasswordCodeStore{Code: validCode},
				&mockPasswordComparer{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &NewPasswordCommand{
				InitiatorID: activeUser.ID,
				UserID:      activeUser.ID,
				OldPassword: "old_password",
				NewPassword: "new_password",
				Code:        validCode,
			},
		},
		{
			TestName: "test_new_password_use_case_invalid_code",
			Expected: ErrInvalidData,
			UC: MustNewPasswordUseCase(
				&mockNewPasswordRepository{User: activeUser},
				&mockNewPasswordCodeStore{Code: validCode},
				&mockPasswordComparer{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &NewPasswordCommand{
				InitiatorID: activeUser.ID,
				UserID:      activeUser.ID,
				OldPassword: "old_password",
				NewPassword: "new_password",
				Code:        invalidCode,
			},
		},
		{
			TestName: "test_new_password_use_case_user_not_active",
			Expected: ErrUserNotActive,
			UC: MustNewPasswordUseCase(
				&mockNewPasswordRepository{User: notActiveUser},
				&mockNewPasswordCodeStore{Code: validCode},
				&mockPasswordComparer{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &NewPasswordCommand{
				InitiatorID: notActiveUser.ID,
				UserID:      notActiveUser.ID,
				OldPassword: "old_password",
				NewPassword: "new_password",
				Code:        validCode,
			},
		},
		{
			TestName: "test_new_password_use_case_by_id_error",
			Expected: ErrInternal,
			UC: MustNewPasswordUseCase(
				&mockNewPasswordRepository{User: activeUser, ErrByID: ErrInternal},
				&mockNewPasswordCodeStore{Code: validCode},
				&mockPasswordComparer{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &NewPasswordCommand{
				InitiatorID: activeUser.ID,
				UserID:      activeUser.ID,
				OldPassword: "old_password",
				NewPassword: "new_password",
				Code:        validCode,
			},
		},
		{
			TestName: "test_new_password_use_case_save_error",
			Expected: ErrInternal,
			UC: MustNewPasswordUseCase(
				&mockNewPasswordRepository{User: activeUser, ErrSave: ErrInternal},
				&mockNewPasswordCodeStore{Code: validCode},
				&mockPasswordComparer{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &NewPasswordCommand{
				InitiatorID: activeUser.ID,
				UserID:      activeUser.ID,
				OldPassword: "old_password",
				NewPassword: "new_password",
				Code:        validCode,
			},
		},
		{
			TestName: "test_new_password_use_case_get_code_error",
			Expected: ErrInternal,
			UC: MustNewPasswordUseCase(
				&mockNewPasswordRepository{User: activeUser},
				&mockNewPasswordCodeStore{Code: validCode, ErrGet: ErrInternal},
				&mockPasswordComparer{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &NewPasswordCommand{
				InitiatorID: activeUser.ID,
				UserID:      activeUser.ID,
				OldPassword: "old_password",
				NewPassword: "new_password",
				Code:        validCode,
			},
		},
		{
			TestName: "test_new_password_use_case_del_code_error",
			Expected: ErrInternal,
			UC: MustNewPasswordUseCase(
				&mockNewPasswordRepository{User: activeUser},
				&mockNewPasswordCodeStore{Code: validCode, ErrDel: ErrInternal},
				&mockPasswordComparer{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &NewPasswordCommand{
				InitiatorID: activeUser.ID,
				UserID:      activeUser.ID,
				OldPassword: "old_password",
				NewPassword: "new_password",
				Code:        validCode,
			},
		},
		{
			TestName: "test_new_password_use_case_password_error",
			Expected: ErrInvalidData,
			UC: MustNewPasswordUseCase(
				&mockNewPasswordRepository{User: activeUser},
				&mockNewPasswordCodeStore{Code: validCode},
				&mockPasswordComparer{InvalidPassword: []string{"old_password"}},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &NewPasswordCommand{
				InitiatorID: activeUser.ID,
				UserID:      activeUser.ID,
				OldPassword: "old_password",
				NewPassword: "new_password",
				Code:        validCode,
			},
		},
		{
			TestName: "test_new_password_use_case_compare_password_error",
			Expected: ErrInternal,
			UC: MustNewPasswordUseCase(
				&mockNewPasswordRepository{User: activeUser},
				&mockNewPasswordCodeStore{Code: validCode},
				&mockPasswordComparer{Err: ErrInternal},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &NewPasswordCommand{
				InitiatorID: activeUser.ID,
				UserID:      activeUser.ID,
				OldPassword: "old_password",
				NewPassword: "new_password",
				Code:        validCode,
			},
		},
		{
			TestName: "test_new_password_use_case_invalid_new_password",
			Expected: ErrInvalidData,
			UC: MustNewPasswordUseCase(
				&mockNewPasswordRepository{User: activeUser},
				&mockNewPasswordCodeStore{Code: validCode},
				&mockPasswordComparer{},
				&mockPasswordValidator{InvalidPasswords: []string{"new_password"}},
				&mockPasswordHasher{},
			),
			Command: &NewPasswordCommand{
				InitiatorID: activeUser.ID,
				UserID:      activeUser.ID,
				OldPassword: "old_password",
				NewPassword: "new_password",
				Code:        validCode,
			},
		},
		{
			TestName: "test_new_password_use_case_hasher_error",
			Expected: ErrInternal,
			UC: MustNewPasswordUseCase(
				&mockNewPasswordRepository{User: activeUser},
				&mockNewPasswordCodeStore{Code: validCode},
				&mockPasswordComparer{},
				&mockPasswordValidator{},
				&mockPasswordHasher{Err: ErrInternal},
			),
			Command: &NewPasswordCommand{
				InitiatorID: activeUser.ID,
				UserID:      activeUser.ID,
				OldPassword: "old_password",
				NewPassword: "new_password",
				Code:        validCode,
			},
		},
		{
			TestName: "test_new_password_use_case_ok",
			Expected: nil,
			UC: MustNewPasswordUseCase(
				&mockNewPasswordRepository{User: activeUser},
				&mockNewPasswordCodeStore{Code: validCode},
				&mockPasswordComparer{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &NewPasswordCommand{
				InitiatorID: activeUser.ID,
				UserID:      activeUser.ID,
				OldPassword: "old_password",
				NewPassword: "new_password",
				Code:        validCode,
			},
		},
		{
			TestName: "test_new_password_use_case_new_password_is_same",
			Expected: ErrInvalidData,
			UC: MustNewPasswordUseCase(
				&mockNewPasswordRepository{User: activeUser},
				&mockNewPasswordCodeStore{Code: validCode},
				&mockPasswordComparer{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &NewPasswordCommand{
				InitiatorID: activeUser.ID,
				UserID:      activeUser.ID,
				OldPassword: "old_password",
				NewPassword: "old_password",
				Code:        validCode,
			},
		},
		{
			TestName: "test_new_password_use_case_initiator_and_user_are_diff",
			Expected: ErrNotAllowed,
			UC: MustNewPasswordUseCase(
				&mockNewPasswordRepository{User: activeUser},
				&mockNewPasswordCodeStore{Code: validCode},
				&mockPasswordComparer{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &NewPasswordCommand{
				InitiatorID: uuid.New(),
				UserID:      activeUser.ID,
				OldPassword: "old_password",
				NewPassword: "new_password",
				Code:        validCode,
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

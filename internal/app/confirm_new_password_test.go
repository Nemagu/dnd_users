package app

import (
	"context"
	"errors"
	"testing"

	"github.com/Nemagu/dnd_users/internal/domain"
	"github.com/google/uuid"
)

type mockConfirmNewPasswordRepository struct {
	User *User
	Err  error
}

func (m *mockConfirmNewPasswordRepository) ByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return m.User, m.Err
}

type mockConfirmNewPasswordCodeStore struct {
	Err error
}

func (m *mockConfirmNewPasswordCodeStore) SetNewPassword(
	ctx context.Context,
	key, value string,
) error {
	return m.Err
}

type mockConfirmNewPasswordEmailProvider struct{}

func (m *mockConfirmNewPasswordEmailProvider) SendConfirmationNewPassword(data EmailCode) {}

func TestConfirmNewPasswordUseCase_Execute(t *testing.T) {
	user := &User{
		ID:           uuid.New(),
		Email:        "test@mail.com",
		State:        domain.ACTIVE,
		Status:       domain.USER,
		PasswordHash: "hash",
		Version:      1,
	}
	cases := []struct {
		TestName string
		Expected error
		UC       *ConfirmNewPasswordUseCase
		Command  *ConfirmNewPasswordCommand
	}{
		{
			TestName: "test_confirm_new_password_use_case_ok",
			Expected: nil,
			UC: MustConfirmNewPasswordUseCase(
				&mockConfirmNewPasswordRepository{User: user},
				&mockConfirmNewPasswordCodeStore{},
				&mockConfirmNewPasswordEmailProvider{},
				&mockCodeGenerator{},
			),
			Command: &ConfirmNewPasswordCommand{
				InitiatorID: user.ID,
				UserID:      user.ID,
			},
		},
		{
			TestName: "test_confirm_new_password_use_case_initiator_and_user_are diff",
			Expected: ErrNotAllowed,
			UC: MustConfirmNewPasswordUseCase(
				&mockConfirmNewPasswordRepository{User: user},
				&mockConfirmNewPasswordCodeStore{},
				&mockConfirmNewPasswordEmailProvider{},
				&mockCodeGenerator{},
			),
			Command: &ConfirmNewPasswordCommand{
				InitiatorID: uuid.New(),
				UserID:      user.ID,
			},
		},
		{
			TestName: "test_confirm_new_password_use_case_by_id_error",
			Expected: ErrInternal,
			UC: MustConfirmNewPasswordUseCase(
				&mockConfirmNewPasswordRepository{User: user, Err: ErrInternal},
				&mockConfirmNewPasswordCodeStore{},
				&mockConfirmNewPasswordEmailProvider{},
				&mockCodeGenerator{},
			),
			Command: &ConfirmNewPasswordCommand{
				InitiatorID: user.ID,
				UserID:      user.ID,
			},
		},
		{
			TestName: "test_confirm_new_password_use_case_setting_code_error",
			Expected: ErrInternal,
			UC: MustConfirmNewPasswordUseCase(
				&mockConfirmNewPasswordRepository{User: user},
				&mockConfirmNewPasswordCodeStore{Err: ErrInternal},
				&mockConfirmNewPasswordEmailProvider{},
				&mockCodeGenerator{},
			),
			Command: &ConfirmNewPasswordCommand{
				InitiatorID: user.ID,
				UserID:      user.ID,
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

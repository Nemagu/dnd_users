package app

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/Nemagu/dnd_users/internal/domain"
	"github.com/google/uuid"
)

type mockConfirmNewEmailRepository struct {
	ExistsEmails []string
	User         *User
	ErrExists    error
	ErrByID      error
}

func (m *mockConfirmNewEmailRepository) ByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return m.User, m.ErrByID
}

func (m *mockConfirmNewEmailRepository) EmailExists(
	ctx context.Context,
	email string,
) (bool, error) {
	return slices.Contains(m.ExistsEmails, email), m.ErrExists
}

type mockConfirmNewEmailCodeStore struct {
	Err error
}

func (m *mockConfirmNewEmailCodeStore) SetCNEC(ctx context.Context, key, value string) error {
	return m.Err
}

type mockConfirmNewEmailProvider struct{}

func (m *mockConfirmNewEmailProvider) SendConfirmationNewEmail(data []EmailCode) {}

func TestConfirmNewEmailUseCase_Execute(t *testing.T) {
	user := &User{
		ID:           uuid.New(),
		Email:        "test@mail.com",
		State:        domain.ACTIVE,
		Status:       domain.USER,
		PasswordHash: "test",
		Version:      1,
	}
	command := &ConfirmNewEmailCommand{
		InitiatorID: user.ID,
		UserID:      user.ID,
		NewEmail:    "new.email@mail.com",
	}
	cases := []struct {
		TestName string
		Expected error
		UC       *ConfirmNewEmailUseCase
		Command  *ConfirmNewEmailCommand
	}{
		{
			TestName: "test_confirm_new_email_use_case_ok",
			Expected: nil,
			UC: MustConfirmNewEmailUseCase(
				&mockConfirmNewEmailRepository{User: user},
				&mockConfirmNewEmailCodeStore{},
				&mockEmailValidator{},
				&mockConfirmNewEmailProvider{},
				&mockCodeGenerator{},
			),
			Command: command,
		},
		{
			TestName: "test_confirm_new_email_use_case_email_exists",
			Expected: ErrAlreadyExists,
			UC: MustConfirmNewEmailUseCase(
				&mockConfirmNewEmailRepository{
					User:         user,
					ExistsEmails: []string{command.NewEmail},
				},
				&mockConfirmNewEmailCodeStore{},
				&mockEmailValidator{},
				&mockConfirmNewEmailProvider{},
				&mockCodeGenerator{},
			),
			Command: command,
		},
		{
			TestName: "test_confirm_new_email_use_case_check_email_exists_error",
			Expected: ErrInternal,
			UC: MustConfirmNewEmailUseCase(
				&mockConfirmNewEmailRepository{
					User:      user,
					ErrExists: ErrInternal,
				},
				&mockConfirmNewEmailCodeStore{},
				&mockEmailValidator{},
				&mockConfirmNewEmailProvider{},
				&mockCodeGenerator{},
			),
			Command: command,
		},
		{
			TestName: "test_confirm_new_email_use_case_get_user_error",
			Expected: ErrInternal,
			UC: MustConfirmNewEmailUseCase(
				&mockConfirmNewEmailRepository{User: user, ErrByID: ErrInternal},
				&mockConfirmNewEmailCodeStore{},
				&mockEmailValidator{},
				&mockConfirmNewEmailProvider{},
				&mockCodeGenerator{},
			),
			Command: command,
		},
		{
			TestName: "test_confirm_new_email_use_case_set_codes_error",
			Expected: ErrInternal,
			UC: MustConfirmNewEmailUseCase(
				&mockConfirmNewEmailRepository{User: user},
				&mockConfirmNewEmailCodeStore{Err: ErrInternal},
				&mockEmailValidator{},
				&mockConfirmNewEmailProvider{},
				&mockCodeGenerator{},
			),
			Command: command,
		},
		{
			TestName: "test_confirm_new_email_use_case_invalid_email",
			Expected: ErrInvalidData,
			UC: MustConfirmNewEmailUseCase(
				&mockConfirmNewEmailRepository{User: user},
				&mockConfirmNewEmailCodeStore{},
				&mockEmailValidator{InvalidEmails: []string{command.NewEmail}},
				&mockConfirmNewEmailProvider{},
				&mockCodeGenerator{},
			),
			Command: command,
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

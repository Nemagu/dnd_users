package app

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/Nemagu/dnd_users/internal/domain"
	"github.com/google/uuid"
)

type mockNewEmailRepository struct {
	User         *User
	ErrByID      error
	ErrSave      error
	ExistsEmails []string
	ErrExists    error
}

func (m *mockNewEmailRepository) ByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return m.User, m.ErrByID
}

func (m *mockNewEmailRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	return slices.Contains(m.ExistsEmails, email), m.ErrExists
}

func (m *mockNewEmailRepository) Save(ctx context.Context, user *User) error {
	return m.ErrSave
}

type mockNewEmailCodeStore struct {
	Err          error
	NewEmailCode string
	NewEmailKey  string
	OldEmailCode string
	OldEmailKey  string
}

func (m *mockNewEmailCodeStore) GetNEC(ctx context.Context, key string) (string, error) {
	switch key {
	case m.NewEmailKey:
		return m.NewEmailCode, m.Err
	case m.OldEmailKey:
		return m.OldEmailCode, m.Err
	default:
		return "", m.Err
	}
}

func TestNewEmailUseCase_Execute(t *testing.T) {
	user := &User{
		ID:           uuid.New(),
		Email:        "test@mail.com",
		State:        domain.ACTIVE,
		Status:       domain.USER,
		PasswordHash: "test",
		Version:      1,
	}
	newEmail := "new.test@mail.com"
	newEmailCode := "new"
	newEmailKey := newEmail
	oldEmailCode := "old"
	oldEmailKey := user.Email
	password := "password"
	cases := []struct {
		TestName string
		Expected error
		UC       *NewEmailUseCase
		Command  *NewEmailCommand
	}{
		{
			TestName: "test_new_email_use_case_ok",
			Expected: nil,
			UC: MustNewEmailUseCase(
				&mockNewEmailRepository{User: user},
				&mockNewEmailCodeStore{
					NewEmailCode: newEmailCode,
					NewEmailKey:  newEmailKey,
					OldEmailCode: oldEmailCode,
					OldEmailKey:  oldEmailKey,
				},
				&mockEmailValidator{},
				&mockPasswordComparer{},
			),
			Command: &NewEmailCommand{
				InitiatorID:  user.ID,
				UserID:       user.ID,
				NewEmail:     newEmail,
				NewEmailCode: newEmailCode,
				OldEmailCode: oldEmailCode,
				Password:     password,
			},
		},
		{
			TestName: "test_new_email_use_case_email_exists",
			Expected: ErrAlreadyExists,
			UC: MustNewEmailUseCase(
				&mockNewEmailRepository{User: user, ExistsEmails: []string{newEmail}},
				&mockNewEmailCodeStore{
					NewEmailCode: newEmailCode,
					NewEmailKey:  newEmailKey,
					OldEmailCode: oldEmailCode,
					OldEmailKey:  oldEmailKey,
				},
				&mockEmailValidator{},
				&mockPasswordComparer{},
			),
			Command: &NewEmailCommand{
				InitiatorID:  user.ID,
				UserID:       user.ID,
				NewEmail:     newEmail,
				NewEmailCode: newEmailCode,
				OldEmailCode: oldEmailCode,
				Password:     password,
			},
		},
		{
			TestName: "test_new_email_use_case_check_email_error",
			Expected: ErrInternal,
			UC: MustNewEmailUseCase(
				&mockNewEmailRepository{User: user, ErrExists: ErrInternal},
				&mockNewEmailCodeStore{
					NewEmailCode: newEmailCode,
					NewEmailKey:  newEmailKey,
					OldEmailCode: oldEmailCode,
					OldEmailKey:  oldEmailKey,
				},
				&mockEmailValidator{},
				&mockPasswordComparer{},
			),
			Command: &NewEmailCommand{
				InitiatorID:  user.ID,
				UserID:       user.ID,
				NewEmail:     newEmail,
				NewEmailCode: newEmailCode,
				OldEmailCode: oldEmailCode,
				Password:     password,
			},
		},
		{
			TestName: "test_new_email_use_case_saving_error",
			Expected: ErrInternal,
			UC: MustNewEmailUseCase(
				&mockNewEmailRepository{User: user, ErrSave: ErrInternal},
				&mockNewEmailCodeStore{
					NewEmailCode: newEmailCode,
					NewEmailKey:  newEmailKey,
					OldEmailCode: oldEmailCode,
					OldEmailKey:  oldEmailKey,
				},
				&mockEmailValidator{},
				&mockPasswordComparer{},
			),
			Command: &NewEmailCommand{
				InitiatorID:  user.ID,
				UserID:       user.ID,
				NewEmail:     newEmail,
				NewEmailCode: newEmailCode,
				OldEmailCode: oldEmailCode,
				Password:     password,
			},
		},
		{
			TestName: "test_new_email_use_case_getting_error",
			Expected: ErrInternal,
			UC: MustNewEmailUseCase(
				&mockNewEmailRepository{User: user, ErrByID: ErrInternal},
				&mockNewEmailCodeStore{
					NewEmailCode: newEmailCode,
					NewEmailKey:  newEmailKey,
					OldEmailCode: oldEmailCode,
					OldEmailKey:  oldEmailKey,
				},
				&mockEmailValidator{},
				&mockPasswordComparer{},
			),
			Command: &NewEmailCommand{
				InitiatorID:  user.ID,
				UserID:       user.ID,
				NewEmail:     newEmail,
				NewEmailCode: newEmailCode,
				OldEmailCode: oldEmailCode,
				Password:     password,
			},
		},
		{
			TestName: "test_new_email_use_case_getting_codes_error",
			Expected: ErrInternal,
			UC: MustNewEmailUseCase(
				&mockNewEmailRepository{User: user},
				&mockNewEmailCodeStore{
					NewEmailCode: newEmailCode,
					NewEmailKey:  newEmailKey,
					OldEmailCode: oldEmailCode,
					OldEmailKey:  oldEmailKey,
					Err:          ErrInternal,
				},
				&mockEmailValidator{},
				&mockPasswordComparer{},
			),
			Command: &NewEmailCommand{
				InitiatorID:  user.ID,
				UserID:       user.ID,
				NewEmail:     newEmail,
				NewEmailCode: newEmailCode,
				OldEmailCode: oldEmailCode,
				Password:     password,
			},
		},
		{
			TestName: "test_new_email_use_case_new_email_code_invalid",
			Expected: ErrInvalidData,
			UC: MustNewEmailUseCase(
				&mockNewEmailRepository{User: user},
				&mockNewEmailCodeStore{
					NewEmailCode: newEmailCode,
					NewEmailKey:  newEmailKey,
					OldEmailCode: oldEmailCode,
					OldEmailKey:  oldEmailKey,
				},
				&mockEmailValidator{},
				&mockPasswordComparer{},
			),
			Command: &NewEmailCommand{
				InitiatorID:  user.ID,
				UserID:       user.ID,
				NewEmail:     newEmail,
				NewEmailCode: "newEmailCode",
				OldEmailCode: oldEmailCode,
				Password:     password,
			},
		},
		{
			TestName: "test_new_email_use_case_old_email_code_invalid",
			Expected: ErrInvalidData,
			UC: MustNewEmailUseCase(
				&mockNewEmailRepository{User: user},
				&mockNewEmailCodeStore{
					NewEmailCode: newEmailCode,
					NewEmailKey:  newEmailKey,
					OldEmailCode: oldEmailCode,
					OldEmailKey:  oldEmailKey,
				},
				&mockEmailValidator{},
				&mockPasswordComparer{},
			),
			Command: &NewEmailCommand{
				InitiatorID:  user.ID,
				UserID:       user.ID,
				NewEmail:     newEmail,
				NewEmailCode: newEmailCode,
				OldEmailCode: "oldEmailCode",
				Password:     password,
			},
		},
		{
			TestName: "test_new_email_use_case_email_invalid",
			Expected: ErrInvalidData,
			UC: MustNewEmailUseCase(
				&mockNewEmailRepository{User: user},
				&mockNewEmailCodeStore{
					NewEmailCode: newEmailCode,
					NewEmailKey:  newEmailKey,
					OldEmailCode: oldEmailCode,
					OldEmailKey:  oldEmailKey,
				},
				&mockEmailValidator{InvalidEmails: []string{newEmail}},
				&mockPasswordComparer{},
			),
			Command: &NewEmailCommand{
				InitiatorID:  user.ID,
				UserID:       user.ID,
				NewEmail:     newEmail,
				NewEmailCode: newEmailCode,
				OldEmailCode: oldEmailCode,
				Password:     password,
			},
		},
		{
			TestName: "test_new_email_use_case_password_invalid",
			Expected: ErrInvalidData,
			UC: MustNewEmailUseCase(
				&mockNewEmailRepository{User: user},
				&mockNewEmailCodeStore{
					NewEmailCode: newEmailCode,
					NewEmailKey:  newEmailKey,
					OldEmailCode: oldEmailCode,
					OldEmailKey:  oldEmailKey,
				},
				&mockEmailValidator{},
				&mockPasswordComparer{InvalidPassword: []string{password}},
			),
			Command: &NewEmailCommand{
				InitiatorID:  user.ID,
				UserID:       user.ID,
				NewEmail:     newEmail,
				NewEmailCode: newEmailCode,
				OldEmailCode: oldEmailCode,
				Password:     password,
			},
		},
		{
			TestName: "test_new_email_use_case_password_check_error",
			Expected: ErrInternal,
			UC: MustNewEmailUseCase(
				&mockNewEmailRepository{User: user},
				&mockNewEmailCodeStore{
					NewEmailCode: newEmailCode,
					NewEmailKey:  newEmailKey,
					OldEmailCode: oldEmailCode,
					OldEmailKey:  oldEmailKey,
				},
				&mockEmailValidator{},
				&mockPasswordComparer{Err: ErrInternal},
			),
			Command: &NewEmailCommand{
				InitiatorID:  user.ID,
				UserID:       user.ID,
				NewEmail:     newEmail,
				NewEmailCode: newEmailCode,
				OldEmailCode: oldEmailCode,
				Password:     password,
			},
		},
		{
			TestName: "test_new_email_use_case_initiator_and_user_diff",
			Expected: ErrNotAllowed,
			UC: MustNewEmailUseCase(
				&mockNewEmailRepository{User: user},
				&mockNewEmailCodeStore{
					NewEmailCode: newEmailCode,
					NewEmailKey:  newEmailKey,
					OldEmailCode: oldEmailCode,
					OldEmailKey:  oldEmailKey,
				},
				&mockEmailValidator{},
				&mockPasswordComparer{},
			),
			Command: &NewEmailCommand{
				InitiatorID:  uuid.New(),
				UserID:       user.ID,
				NewEmail:     newEmail,
				NewEmailCode: newEmailCode,
				OldEmailCode: oldEmailCode,
				Password:     password,
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

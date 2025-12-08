package app

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/google/uuid"
)

type mockRegistrationRepository struct {
	ExistsEmails []string
	ErrNextID    error
	ErrExists    error
	ErrSave      error
}

func (m *mockRegistrationRepository) NextID(ctx context.Context) (uuid.UUID, error) {
	return uuid.New(), m.ErrNextID
}

func (m *mockRegistrationRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	return slices.Contains(m.ExistsEmails, email), m.ErrExists
}

func (m *mockRegistrationRepository) Save(ctx context.Context, user *User) error {
	return m.ErrSave
}

type mockRegistrationCodeStore struct {
	Value       string
	ErrGetting  error
	ErrDeleting error
}

func (m *mockRegistrationCodeStore) GetConfirmEmail(ctx context.Context, key string) (string, error) {
	return m.Value, m.ErrGetting
}

func (m *mockRegistrationCodeStore) DelConfirmEmail(ctx context.Context, key string) error {
	return m.ErrDeleting
}

func TestRegistrationUseCase_Execute(t *testing.T) {
	validEmail := "test@mail.com"
	validPassword := "password"
	validCode := "123456"
	cases := []struct {
		TestName string
		Expected error
		UC       *RegistrationUseCase
		Command  *RegistrationCommand
	}{
		{
			TestName: "test_registration_use_case_ok",
			Expected: nil,
			UC: MustRegistrationUseCase(
				&mockRegistrationRepository{},
				&mockRegistrationCodeStore{Value: validCode},
				&mockEmailValidator{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &RegistrationCommand{
				Email:    validEmail,
				Password: validPassword,
				Code:     validCode,
			},
		},
		{
			TestName: "test_registration_use_case_next_id_error",
			Expected: ErrInternal,
			UC: MustRegistrationUseCase(
				&mockRegistrationRepository{ErrNextID: ErrInternal},
				&mockRegistrationCodeStore{Value: validCode},
				&mockEmailValidator{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &RegistrationCommand{
				Email:    validEmail,
				Password: validPassword,
				Code:     validCode,
			},
		},
		{
			TestName: "test_registration_use_case_email_exist_error",
			Expected: ErrInternal,
			UC: MustRegistrationUseCase(
				&mockRegistrationRepository{ErrExists: ErrInternal},
				&mockRegistrationCodeStore{Value: validCode},
				&mockEmailValidator{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &RegistrationCommand{
				Email:    validEmail,
				Password: validPassword,
				Code:     validCode,
			},
		},
		{
			TestName: "test_registration_use_case_save_error",
			Expected: ErrInternal,
			UC: MustRegistrationUseCase(
				&mockRegistrationRepository{ErrSave: ErrInternal},
				&mockRegistrationCodeStore{Value: validCode},
				&mockEmailValidator{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &RegistrationCommand{
				Email:    validEmail,
				Password: validPassword,
				Code:     validCode,
			},
		},
		{
			TestName: "test_registration_use_case_email_exists",
			Expected: ErrAlreadyExists,
			UC: MustRegistrationUseCase(
				&mockRegistrationRepository{
					ExistsEmails: []string{validEmail},
					ErrExists:    ErrAlreadyExists,
				},
				&mockRegistrationCodeStore{Value: validCode},
				&mockEmailValidator{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &RegistrationCommand{
				Email:    validEmail,
				Password: validPassword,
				Code:     validCode,
			},
		},
		{
			TestName: "test_registration_use_case_getting_code_error",
			Expected: ErrInternal,
			UC: MustRegistrationUseCase(
				&mockRegistrationRepository{},
				&mockRegistrationCodeStore{Value: validCode, ErrGetting: ErrInternal},
				&mockEmailValidator{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &RegistrationCommand{
				Email:    validEmail,
				Password: validPassword,
				Code:     validCode,
			},
		},
		{
			TestName: "test_registration_use_case_invalid_code",
			Expected: ErrInvalidData,
			UC: MustRegistrationUseCase(
				&mockRegistrationRepository{},
				&mockRegistrationCodeStore{Value: "validCode"},
				&mockEmailValidator{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &RegistrationCommand{
				Email:    validEmail,
				Password: validPassword,
				Code:     validCode,
			},
		},
		{
			TestName: "test_registration_use_case_invalid_email",
			Expected: ErrInvalidData,
			UC: MustRegistrationUseCase(
				&mockRegistrationRepository{},
				&mockRegistrationCodeStore{Value: validCode},
				&mockEmailValidator{InvalidEmails: []string{validEmail}},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &RegistrationCommand{
				Email:    validEmail,
				Password: validPassword,
				Code:     validCode,
			},
		},
		{
			TestName: "test_registration_use_case_invalid_password",
			Expected: ErrInvalidData,
			UC: MustRegistrationUseCase(
				&mockRegistrationRepository{},
				&mockRegistrationCodeStore{Value: validCode},
				&mockEmailValidator{},
				&mockPasswordValidator{InvalidPasswords: []string{validPassword}},
				&mockPasswordHasher{},
			),
			Command: &RegistrationCommand{
				Email:    validEmail,
				Password: validPassword,
				Code:     validCode,
			},
		},
		{
			TestName: "test_registration_use_case_hasher_error",
			Expected: ErrInternal,
			UC: MustRegistrationUseCase(
				&mockRegistrationRepository{},
				&mockRegistrationCodeStore{Value: validCode},
				&mockEmailValidator{},
				&mockPasswordValidator{},
				&mockPasswordHasher{Err: ErrInternal},
			),
			Command: &RegistrationCommand{
				Email:    validEmail,
				Password: validPassword,
				Code:     validCode,
			},
		},
		{
			TestName: "test_registration_use_case_code_delete_error",
			Expected: ErrInternal,
			UC: MustRegistrationUseCase(
				&mockRegistrationRepository{},
				&mockRegistrationCodeStore{Value: validCode, ErrDeleting: ErrInternal},
				&mockEmailValidator{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
			),
			Command: &RegistrationCommand{
				Email:    validEmail,
				Password: validPassword,
				Code:     validCode,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			_, err := c.UC.Execute(context.Background(), c.Command)
			if c.Expected == nil {
				if err != nil {
					t.Errorf("expected nil, but got %v", err)
				}
			}
			if !errors.Is(err, c.Expected) {
				t.Errorf("got not expected error: %v", err)
			}
		})
	}
}

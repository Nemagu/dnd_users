package usecase

import (
	"context"
	"testing"

	"github.com/Nemagu/dnd/internal/application"
	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/Nemagu/dnd/internal/domain/duser"
	"github.com/google/uuid"
)

type mockConfirmNewEmailUserRepository struct {
	UserID uuid.UUID
	Exists bool
}

func (m *mockConfirmNewEmailUserRepository) ByID(
	ctx context.Context,
	id uuid.UUID,
) (*appdto.User, error) {
	switch id {
	case m.UserID:
		return &appdto.User{
			UserID:       id,
			State:        duser.ACTIVE,
			Status:       duser.ORDINARY,
			Email:        "email@mail.com",
			PasswordHash: "password_hash",
			Version:      1,
		}, nil
	default:
		return nil, application.ErrNotFound
	}
}

func (m *mockConfirmNewEmailUserRepository) EmailExists(
	ctx context.Context,
	email string,
) (bool, error) {
	return m.Exists, nil
}

func TestConfirmNewEmailUseCase_Execute_Success(t *testing.T) {
	userID := uuid.New()
	newEmail := "new@mail.com"
	uc := MustNewConfirmNewEmailUseCase(
		&mockConfirmNewEmailUserRepository{UserID: userID, Exists: false},
		&mockPasswordComparer{IsCompare: true},
		&mockEmailCrypter{},
		&mockEmailValidator{IsValid: true},
		&mockEmailProvider{},
	)
	if err := uc.Execute(context.Background(), &appdto.ConfirmNewEmailCommand{
		UserID:   userID,
		Email:    newEmail,
		Password: "this_is_password",
	}); err != nil {
		t.Errorf("got error in during execute: %s", err)
	}
}

func TestConfirmNewEmailUseCase_Execute_Fail(t *testing.T) {
	userID := uuid.New()
	newEmail := "new@mail.com"
	password := "this_is_password"
	cases := []struct {
		TestName string
		UC       *ConfirmNewEmailUseCase
		Command  *appdto.ConfirmNewEmailCommand
	}{
		{
			TestName: "test_email_exists",
			UC: MustNewConfirmNewEmailUseCase(
				&mockConfirmNewEmailUserRepository{UserID: userID, Exists: true},
				&mockPasswordComparer{IsCompare: true},
				&mockEmailCrypter{},
				&mockEmailValidator{IsValid: true},
				&mockEmailProvider{},
			),
			Command: &appdto.ConfirmNewEmailCommand{
				UserID:   userID,
				Email:    newEmail,
				Password: password,
			},
		},
		{
			TestName: "test_password_not_valid",
			UC: MustNewConfirmNewEmailUseCase(
				&mockConfirmNewEmailUserRepository{UserID: userID, Exists: false},
				&mockPasswordComparer{IsCompare: false},
				&mockEmailCrypter{},
				&mockEmailValidator{IsValid: true},
				&mockEmailProvider{},
			),
			Command: &appdto.ConfirmNewEmailCommand{
				UserID:   userID,
				Email:    newEmail,
				Password: password,
			},
		},
		{
			TestName: "test_email_not_valid",
			UC: MustNewConfirmNewEmailUseCase(
				&mockConfirmNewEmailUserRepository{UserID: userID, Exists: false},
				&mockPasswordComparer{IsCompare: true},
				&mockEmailCrypter{},
				&mockEmailValidator{IsValid: false},
				&mockEmailProvider{},
			),
			Command: &appdto.ConfirmNewEmailCommand{
				UserID:   userID,
				Email:    newEmail,
				Password: password,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			if err := c.UC.Execute(context.Background(), c.Command); err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

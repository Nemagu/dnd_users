package usecase

import (
	"context"
	"testing"

	appdto "github.com/Nemagu/dnd/internal/application/dto"
)

type mockConfirmEmailUserRepository struct {
	Exists bool
}

func (m *mockConfirmEmailUserRepository) EmailExists(
	ctx context.Context,
	email string,
) (bool, error) {
	return m.Exists, nil
}

func TestConfirmEmailUseCase_Execute_Success(t *testing.T) {
	email := "email@mail.com"
	uc := MustNewConfirmEmailUseCase(
		&mockConfirmEmailUserRepository{Exists: false},
		&mockEmailCrypter{},
		&mockEmailProvider{},
		&mockEmailValidator{IsValid: true},
	)
	if err := uc.Execute(context.Background(), &appdto.ConfirmEmailCommand{
		Email: email,
	}); err != nil {
		t.Errorf("got error in during execute: %s", err)
	}
}

func TestConfirmEmailUseCase_Execute_Fail(t *testing.T) {
	cases := []struct {
		TestName string
		UC       *ConfirmEmailUseCase
		Command  *appdto.ConfirmEmailCommand
	}{
		{
			TestName: "test_email_exists",
			UC: MustNewConfirmEmailUseCase(
				&mockConfirmEmailUserRepository{Exists: true},
				&mockEmailCrypter{},
				&mockEmailProvider{},
				&mockEmailValidator{IsValid: true},
			),
			Command: &appdto.ConfirmEmailCommand{Email: "email@mail.com"},
		},
		{
			TestName: "test_email_is_not_valid",
			UC: MustNewConfirmEmailUseCase(
				&mockConfirmEmailUserRepository{Exists: false},
				&mockEmailCrypter{},
				&mockEmailProvider{},
				&mockEmailValidator{IsValid: false},
			),
			Command: &appdto.ConfirmEmailCommand{Email: "email@mail.com"},
		},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			if err := c.UC.Execute(context.Background(), c.Command); err == nil {
				t.Errorf("expected error, got nil")
			}
		})
	}
}

package usecase

import (
	"context"
	"testing"

	"github.com/Nemagu/dnd/internal/application"
	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/Nemagu/dnd/internal/domain/duser"
	"github.com/google/uuid"
)

type mockAuthRepo struct {
	Email string
}

func (r *mockAuthRepo) ByEmail(ctx context.Context, email string) (*appdto.User, error) {
	if r.Email != email {
		return nil, application.ErrNotFound
	}
	return &appdto.User{
		UserID:       uuid.New(),
		Email:        email,
		State:        duser.ACTIVE,
		Status:       duser.ORDINARY,
		PasswordHash: "this_is_hash_of_password",
		Version:      1,
	}, nil
}

func TestAuthenticate_Execute_Success(t *testing.T) {
	email := "test@example.com"
	repo := &mockAuthRepo{
		Email: email,
	}
	passComparer := &mockPasswordComparer{
		IsCompare: true,
	}
	uc := MustNewAuthenticateUseCase(repo, passComparer)
	_, err := uc.Execute(t.Context(), &appdto.AuthenticateCommand{
		Email: email,
	})
	if err != nil {
		t.Errorf("got error in during execute: %s", err)
	}
}

func TestAuthenticate_Execute_Fail(t *testing.T) {
	cases := []struct {
		UC      *AuthenticateUseCase
		Command *appdto.AuthenticateCommand
	}{
		{
			UC: MustNewAuthenticateUseCase(
				&mockAuthRepo{
					Email: "email",
				},
				&mockPasswordComparer{
					IsCompare: true,
				},
			),
			Command: &appdto.AuthenticateCommand{
				Email: "email1",
			},
		},
		{
			UC: MustNewAuthenticateUseCase(
				&mockAuthRepo{
					Email: "email",
				},
				&mockPasswordComparer{
					IsCompare: false,
				},
			),
			Command: &appdto.AuthenticateCommand{
				Email: "email",
			},
		},
	}
	for _, c := range cases {
		t.Run("test_fail_auth_with_email_"+c.Command.Email, func(t *testing.T) {
			if _, err := c.UC.Execute(t.Context(), c.Command); err == nil {
				t.Errorf("expected error, got nil")
			}
		})
	}
}

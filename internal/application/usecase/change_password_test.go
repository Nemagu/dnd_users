package usecase

import (
	"context"
	"testing"

	"github.com/Nemagu/dnd/internal/application"
	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/Nemagu/dnd/internal/domain/duser"
	"github.com/google/uuid"
)

type mockChangePasswordUserRepository struct {
	UserID   uuid.UUID
	IsSaving bool
}

func (m *mockChangePasswordUserRepository) ByID(
	ctx context.Context,
	id uuid.UUID,
) (*appdto.User, error) {
	if m.UserID != id {
		return nil, application.ErrNotFound
	}
	return &appdto.User{
		UserID:       id,
		State:        duser.ACTIVE,
		Status:       duser.ORDINARY,
		Email:        "email@mail.com",
		PasswordHash: "password_hash",
		Version:      1,
	}, nil
}

func (m *mockChangePasswordUserRepository) Save(ctx context.Context, user *appdto.User) error {
	if m.IsSaving {
		return nil
	}
	return application.ErrInternal
}

func TestChangePasswordUseCase_Execute_Success(t *testing.T) {
	userID := uuid.New()
	uc := MustNewChangePasswordUseCase(
		&mockChangePasswordUserRepository{UserID: userID, IsSaving: true},
		&mockPasswordValidator{IsValidate: true},
		&mockPasswordComparer{IsCompare: true},
		&mockPasswordHasher{},
	)
	if err := uc.Execute(context.Background(), &appdto.ChangePasswordCommand{
		UserID:      userID,
		OldPassword: "old_password",
		NewPassword: "new_password",
	}); err != nil {
		t.Errorf("got error in during execute: %s", err)
	}
}

func TestChangePasswordUseCase_Execute_Fail(t *testing.T) {
	userID := uuid.New()
	badUserID := uuid.New()
	baseTestName := "test_fail_change_password_"
	cases := []struct {
		TestName string
		UC       *ChangePasswordUseCase
		Command  *appdto.ChangePasswordCommand
	}{
		{
			TestName: baseTestName + "diff_user_id",
			UC: MustNewChangePasswordUseCase(
				&mockChangePasswordUserRepository{UserID: userID},
				&mockPasswordValidator{IsValidate: true},
				&mockPasswordComparer{IsCompare: true},
				&mockPasswordHasher{},
			),
			Command: &appdto.ChangePasswordCommand{
				UserID:      badUserID,
				OldPassword: "old_password",
				NewPassword: "new_password",
			},
		},
		{
			TestName: baseTestName + "pass_is_not_valid",
			UC: MustNewChangePasswordUseCase(
				&mockChangePasswordUserRepository{UserID: userID},
				&mockPasswordValidator{IsValidate: false},
				&mockPasswordComparer{IsCompare: true},
				&mockPasswordHasher{},
			),
			Command: &appdto.ChangePasswordCommand{
				UserID:      userID,
				OldPassword: "old_password",
				NewPassword: "new_password",
			},
		},
		{
			TestName: baseTestName + "pass_not_compare",
			UC: MustNewChangePasswordUseCase(
				&mockChangePasswordUserRepository{UserID: userID},
				&mockPasswordValidator{IsValidate: true},
				&mockPasswordComparer{IsCompare: false},
				&mockPasswordHasher{},
			),
			Command: &appdto.ChangePasswordCommand{
				UserID:      userID,
				OldPassword: "old_password",
				NewPassword: "new_password",
			},
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

package usecase

import (
	"context"
	"testing"

	"github.com/Nemagu/dnd/internal/application"
	appdto "github.com/Nemagu/dnd/internal/application/dto"
	"github.com/Nemagu/dnd/internal/domain/duser"
	"github.com/google/uuid"
)

type mockChangeUserRepository struct {
	UserID      uuid.UUID
	InitiatorID uuid.UUID
	Email       string
	IsSaving    bool
}

func (m *mockChangeUserRepository) ByID(
	ctx context.Context,
	id uuid.UUID,
) (*appdto.User, error) {
	switch {
	case id == m.InitiatorID:
		return &appdto.User{
			UserID:       id,
			State:        duser.ACTIVE,
			Status:       duser.ADMIN,
			Email:        "email@mail.com",
			PasswordHash: "password_hash",
			Version:      10,
		}, nil
	case id == m.UserID:
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

func (m *mockChangeUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	if m.Email == email {
		return true, nil
	} else {
		return false, nil
	}
}

func (m *mockChangeUserRepository) Save(ctx context.Context, user *appdto.User) error {
	if m.IsSaving {
		return nil
	}
	return application.ErrInternal
}

func TestChangeUserUseCase_Execute_Success(t *testing.T) {
	email := "email@mail.com"
	userID := uuid.New()
	initiatorID := uuid.New()
	uc := MustNewChangeUserUseCase(
		&mockChangeUserRepository{
			UserID:      userID,
			InitiatorID: initiatorID,
			Email:       email,
			IsSaving:    true,
		},
		duser.MustNewPolicyService(),
		&mockPasswordHasher{},
		&mockPasswordValidator{IsValidate: true},
	)
	if err := uc.Execute(context.Background(), &appdto.ChangeUserCommand{
		InitiatorID: initiatorID,
		UserID:      userID,
		Email:       "email@gmail.com",
		Status:      duser.ADMIN,
		Password:    "new_password",
	}); err != nil {
		t.Errorf("got error in during execute: %s", err)
	}
}

func TestChangeUserUseCase_Execute_Fail(t *testing.T) {
	userID := uuid.New()
	initiatorID := uuid.New()
	badUserID := uuid.New()
	badInitiatorID := uuid.New()
	email := "email@mail.com"
	newEmail := "email@gmail.com"
	cases := []struct {
		TestName string
		UC       *ChangeUserUseCase
		Command  *appdto.ChangeUserCommand
	}{
		{
			TestName: "test_not_succeed_saving",
			UC: MustNewChangeUserUseCase(
				&mockChangeUserRepository{
					InitiatorID: initiatorID,
					UserID:      userID,
					Email:       email,
					IsSaving:    false,
				},
				duser.MustNewPolicyService(),
				&mockPasswordHasher{},
				&mockPasswordValidator{IsValidate: true},
			),
			Command: &appdto.ChangeUserCommand{
				InitiatorID: initiatorID,
				UserID:      userID,
				Email:       newEmail,
				Status:      duser.ADMIN,
				Password:    "new_password",
			},
		},
		{
			TestName: "test_bad_initiator_id",
			UC: MustNewChangeUserUseCase(
				&mockChangeUserRepository{
					InitiatorID: initiatorID,
					UserID:      userID,
					Email:       email,
					IsSaving:    true,
				},
				duser.MustNewPolicyService(),
				&mockPasswordHasher{},
				&mockPasswordValidator{IsValidate: true},
			),
			Command: &appdto.ChangeUserCommand{
				InitiatorID: badInitiatorID,
				UserID:      userID,
				Email:       newEmail,
				Status:      duser.ADMIN,
				Password:    "new_password",
			},
		},
		{
			TestName: "test_bad_user_id",
			UC: MustNewChangeUserUseCase(
				&mockChangeUserRepository{
					InitiatorID: initiatorID,
					UserID:      userID,
					Email:       email,
					IsSaving:    true,
				},
				duser.MustNewPolicyService(),
				&mockPasswordHasher{},
				&mockPasswordValidator{IsValidate: true},
			),
			Command: &appdto.ChangeUserCommand{
				InitiatorID: initiatorID,
				UserID:      badUserID,
				Email:       newEmail,
				Status:      duser.ADMIN,
				Password:    "new_password",
			},
		},
		{
			TestName: "test_no_access",
			UC: MustNewChangeUserUseCase(
				&mockChangeUserRepository{
					InitiatorID: initiatorID,
					UserID:      userID,
					Email:       email,
					IsSaving:    false,
				},
				duser.MustNewPolicyService(),
				&mockPasswordHasher{},
				&mockPasswordValidator{IsValidate: true},
			),
			Command: &appdto.ChangeUserCommand{
				InitiatorID: userID,
				UserID:      initiatorID,
				Email:       newEmail,
				Status:      duser.ADMIN,
				Password:    "new_password",
			},
		},
		{
			TestName: "test_password_is_not_valid",
			UC: MustNewChangeUserUseCase(
				&mockChangeUserRepository{
					InitiatorID: initiatorID,
					UserID:      userID,
					Email:       email,
					IsSaving:    false,
				},
				duser.MustNewPolicyService(),
				&mockPasswordHasher{},
				&mockPasswordValidator{IsValidate: false},
			),
			Command: &appdto.ChangeUserCommand{
				InitiatorID: initiatorID,
				UserID:      userID,
				Email:       newEmail,
				Status:      duser.ADMIN,
				Password:    "new_password",
			},
		},
		{
			TestName: "test_email_exists",
			UC: MustNewChangeUserUseCase(
				&mockChangeUserRepository{
					InitiatorID: initiatorID,
					UserID:      userID,
					Email:       email,
					IsSaving:    false,
				},
				duser.MustNewPolicyService(),
				&mockPasswordHasher{},
				&mockPasswordValidator{IsValidate: true},
			),
			Command: &appdto.ChangeUserCommand{
				InitiatorID: initiatorID,
				UserID:      userID,
				Email:       email,
				Status:      duser.ADMIN,
				Password:    "new_password",
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

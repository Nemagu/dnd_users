package app

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/Nemagu/dnd_users/internal/domain"
	"github.com/google/uuid"
)

type mockChangeUserRepository struct {
	ExistsEmails  []string
	NotExistsIDs  []uuid.UUID
	InitiatorUser *User
	User          *User
	InitiatorID   uuid.UUID
	UserID        uuid.UUID
	ErrEmail      error
	ErrID         error
	ErrSave       error
	ErrByID       error
}

func (m *mockChangeUserRepository) IDExists(ctx context.Context, id uuid.UUID) (bool, error) {
	return !slices.Contains(m.NotExistsIDs, id), m.ErrID
}

func (m *mockChangeUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	return slices.Contains(m.ExistsEmails, email), m.ErrEmail
}

func (m *mockChangeUserRepository) ByID(ctx context.Context, id uuid.UUID) (*User, error) {
	switch id {
	case m.InitiatorID:
		return m.InitiatorUser, m.ErrByID
	case m.UserID:
		return m.User, m.ErrByID
	default:
		return nil, m.ErrByID
	}
}

func (m *mockChangeUserRepository) Save(ctx context.Context, user *User) error {
	return m.ErrSave
}

func TestChangeUserUseCase_Execute(t *testing.T) {
	adminUser := &User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		State:        domain.ACTIVE,
		Status:       domain.ADMIN,
		PasswordHash: "password_hash",
		Version:      3,
	}
	ordinaryUser := &User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		State:        domain.ACTIVE,
		Status:       domain.USER,
		PasswordHash: "password_hash",
		Version:      3,
	}
	notActiveUser := &User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		State:        domain.FROZEN,
		Status:       domain.ADMIN,
		PasswordHash: "password_hash",
		Version:      3,
	}
	cases := []struct {
		TestName string
		Expected error
		UC       *ChangeUserUseCase
		Command  *ChangeUserCommand
	}{
		{
			TestName: "test_change_user_without_state_use_case_ok",
			Expected: nil,
			UC: MustChangeUserUseCase(
				&mockChangeUserRepository{
					InitiatorUser: adminUser,
					User:          ordinaryUser,
					InitiatorID:   adminUser.ID,
					UserID:        ordinaryUser.ID,
				},
				&mockEmailValidator{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
				domain.MustPolicyService(),
			),
			Command: &ChangeUserCommand{
				InitiatorID: adminUser.ID,
				UserID:      ordinaryUser.ID,
				Email:       "new_email@example.com",
				Status:      domain.ADMIN,
				Password:    "new_password",
			},
		},
		{
			TestName: "test_change_user_with_state_use_case_ok",
			Expected: nil,
			UC: MustChangeUserUseCase(
				&mockChangeUserRepository{
					InitiatorUser: adminUser,
					User:          ordinaryUser,
					InitiatorID:   adminUser.ID,
					UserID:        ordinaryUser.ID,
				},
				&mockEmailValidator{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
				domain.MustPolicyService(),
			),
			Command: &ChangeUserCommand{
				InitiatorID: adminUser.ID,
				UserID:      ordinaryUser.ID,
				State:       domain.FROZEN,
			},
		},
		{
			TestName: "test_change_user_use_case_not_admin",
			Expected: ErrNotAllowed,
			UC: MustChangeUserUseCase(
				&mockChangeUserRepository{
					InitiatorUser: ordinaryUser,
					User:          ordinaryUser,
					InitiatorID:   ordinaryUser.ID,
					UserID:        ordinaryUser.ID,
				},
				&mockEmailValidator{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
				domain.MustPolicyService(),
			),
			Command: &ChangeUserCommand{
				InitiatorID: ordinaryUser.ID,
				UserID:      ordinaryUser.ID,
				State:       domain.FROZEN,
			},
		},
		{
			TestName: "test_change_user_use_case_not_active",
			Expected: ErrUserNotActive,
			UC: MustChangeUserUseCase(
				&mockChangeUserRepository{
					InitiatorUser: adminUser,
					User:          notActiveUser,
					InitiatorID:   adminUser.ID,
					UserID:        notActiveUser.ID,
				},
				&mockEmailValidator{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
				domain.MustPolicyService(),
			),
			Command: &ChangeUserCommand{
				InitiatorID: adminUser.ID,
				UserID:      notActiveUser.ID,
				Email:       "new_email@example.com",
			},
		},
		{
			TestName: "test_change_user_use_case_email_exists",
			Expected: ErrInvalidData,
			UC: MustChangeUserUseCase(
				&mockChangeUserRepository{
					InitiatorUser: adminUser,
					User:          ordinaryUser,
					InitiatorID:   adminUser.ID,
					UserID:        ordinaryUser.ID,
					ExistsEmails:  []string{"new_email@example.com"},
				},
				&mockEmailValidator{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
				domain.MustPolicyService(),
			),
			Command: &ChangeUserCommand{
				InitiatorID: adminUser.ID,
				UserID:      ordinaryUser.ID,
				Email:       "new_email@example.com",
			},
		},
		{
			TestName: "test_change_user_use_case_user_id_not_exists",
			Expected: ErrNotFound,
			UC: MustChangeUserUseCase(
				&mockChangeUserRepository{
					InitiatorUser: adminUser,
					User:          ordinaryUser,
					InitiatorID:   adminUser.ID,
					UserID:        ordinaryUser.ID,
					NotExistsIDs:  []uuid.UUID{ordinaryUser.ID},
				},
				&mockEmailValidator{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
				domain.MustPolicyService(),
			),
			Command: &ChangeUserCommand{
				InitiatorID: adminUser.ID,
				UserID:      ordinaryUser.ID,
				Email:       "new_email@example.com",
			},
		},
		{
			TestName: "test_change_user_use_case_email_internal_error",
			Expected: ErrInternal,
			UC: MustChangeUserUseCase(
				&mockChangeUserRepository{
					InitiatorUser: adminUser,
					User:          ordinaryUser,
					InitiatorID:   adminUser.ID,
					UserID:        ordinaryUser.ID,
					ErrEmail:      ErrInternal,
				},
				&mockEmailValidator{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
				domain.MustPolicyService(),
			),
			Command: &ChangeUserCommand{
				InitiatorID: adminUser.ID,
				UserID:      ordinaryUser.ID,
				Email:       "new_email@example.com",
			},
		},
		{
			TestName: "test_change_user_use_case_id_exists_internal_error",
			Expected: ErrInternal,
			UC: MustChangeUserUseCase(
				&mockChangeUserRepository{
					InitiatorUser: adminUser,
					User:          ordinaryUser,
					InitiatorID:   adminUser.ID,
					UserID:        ordinaryUser.ID,
					ErrID:         ErrInternal,
				},
				&mockEmailValidator{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
				domain.MustPolicyService(),
			),
			Command: &ChangeUserCommand{
				InitiatorID: adminUser.ID,
				UserID:      ordinaryUser.ID,
				Email:       "new_email@example.com",
			},
		},
		{
			TestName: "test_change_user_use_case_save_internal_error",
			Expected: ErrInternal,
			UC: MustChangeUserUseCase(
				&mockChangeUserRepository{
					InitiatorUser: adminUser,
					User:          ordinaryUser,
					InitiatorID:   adminUser.ID,
					UserID:        ordinaryUser.ID,
					ErrSave:       ErrInternal,
				},
				&mockEmailValidator{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
				domain.MustPolicyService(),
			),
			Command: &ChangeUserCommand{
				InitiatorID: adminUser.ID,
				UserID:      ordinaryUser.ID,
				Email:       "new_email@example.com",
			},
		},
		{
			TestName: "test_change_user_use_case_getting_by_id_internal_error",
			Expected: ErrInternal,
			UC: MustChangeUserUseCase(
				&mockChangeUserRepository{
					InitiatorUser: adminUser,
					User:          ordinaryUser,
					InitiatorID:   adminUser.ID,
					UserID:        ordinaryUser.ID,
					ErrByID:       ErrInternal,
				},
				&mockEmailValidator{},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
				domain.MustPolicyService(),
			),
			Command: &ChangeUserCommand{
				InitiatorID: adminUser.ID,
				UserID:      ordinaryUser.ID,
				Email:       "new_email@example.com",
			},
		},
		{
			TestName: "test_change_user_use_case_getting_by_id_internal_error",
			Expected: ErrInvalidData,
			UC: MustChangeUserUseCase(
				&mockChangeUserRepository{
					InitiatorUser: adminUser,
					User:          ordinaryUser,
					InitiatorID:   adminUser.ID,
					UserID:        ordinaryUser.ID,
				},
				&mockEmailValidator{
					InvalidEmails: []string{"new_email@example.com"},
				},
				&mockPasswordValidator{},
				&mockPasswordHasher{},
				domain.MustPolicyService(),
			),
			Command: &ChangeUserCommand{
				InitiatorID: adminUser.ID,
				UserID:      ordinaryUser.ID,
				Email:       "new_email@example.com",
			},
		},
		{
			TestName: "test_change_user_use_case_hashing_internal_error",
			Expected: ErrInternal,
			UC: MustChangeUserUseCase(
				&mockChangeUserRepository{
					InitiatorUser: adminUser,
					User:          ordinaryUser,
					InitiatorID:   adminUser.ID,
					UserID:        ordinaryUser.ID,
				},
				&mockEmailValidator{},
				&mockPasswordValidator{},
				&mockPasswordHasher{
					Err: ErrInternal,
				},
				domain.MustPolicyService(),
			),
			Command: &ChangeUserCommand{
				InitiatorID: adminUser.ID,
				UserID:      ordinaryUser.ID,
				Email:       "new_email@example.com",
				Password:    "new_password",
			},
		},
		{
			TestName: "test_change_user_use_case_invalid_password",
			Expected: ErrInvalidData,
			UC: MustChangeUserUseCase(
				&mockChangeUserRepository{
					InitiatorUser: adminUser,
					User:          ordinaryUser,
					InitiatorID:   adminUser.ID,
					UserID:        ordinaryUser.ID,
				},
				&mockEmailValidator{},
				&mockPasswordValidator{
					InvalidPasswords: []string{"new_password"},
				},
				&mockPasswordHasher{},
				domain.MustPolicyService(),
			),
			Command: &ChangeUserCommand{
				InitiatorID: adminUser.ID,
				UserID:      ordinaryUser.ID,
				Email:       "new_email@example.com",
				Password:    "new_password",
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

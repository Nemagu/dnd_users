package domain

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestUser_NewUser(t *testing.T) {
	cases := []struct {
		TestName     string
		Expected     error
		ID           uuid.UUID
		Email        string
		PasswordHash string
	}{
		{
			TestName:     "test_new_user_ok",
			Expected:     nil,
			ID:           uuid.New(),
			Email:        "test@test.ru",
			PasswordHash: "test",
		},
		{
			TestName:     "test_new_user_id_is_empty",
			Expected:     ErrInvalidData,
			ID:           uuid.Nil,
			Email:        "test@test.ru",
			PasswordHash: "test",
		},
		{
			TestName:     "test_new_user_email_is_empty",
			Expected:     ErrInvalidData,
			ID:           uuid.New(),
			Email:        "",
			PasswordHash: "test",
		},
		{
			TestName:     "test_new_user_password_hash_is_empty",
			Expected:     ErrInvalidData,
			ID:           uuid.New(),
			Email:        "test@test.ru",
			PasswordHash: "",
		},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			_, err := NewUser(c.ID, c.Email, c.PasswordHash)
			if c.Expected == nil {
				if err != nil {
					t.Errorf("expected %T, but got nil", c.Expected)
				}
			}
			if !errors.Is(err, c.Expected) {
				t.Errorf("expected %T, but got %T", c.Expected, err)
			}
		})
	}
}

func TestUser_RestoreUser(t *testing.T) {
	cases := []struct {
		TestName     string
		Expected     error
		ID           uuid.UUID
		Email        string
		State        State
		Status       Status
		PasswordHash string
		Version      uint
	}{
		{
			TestName:     "test_restore_user_ok",
			Expected:     nil,
			ID:           uuid.New(),
			Email:        "test@test.ru",
			State:        newActiveState(),
			Status:       newUserStatus(),
			PasswordHash: "test",
			Version:      1,
		},
		{
			TestName:     "test_restore_user_id_is_empty",
			Expected:     ErrInvalidData,
			ID:           uuid.Nil,
			Email:        "test@test.ru",
			State:        newActiveState(),
			Status:       newUserStatus(),
			PasswordHash: "test",
			Version:      1,
		},
		{
			TestName:     "test_restore_user_email_is_empty",
			Expected:     ErrInvalidData,
			ID:           uuid.New(),
			Email:        "",
			State:        newActiveState(),
			Status:       newUserStatus(),
			PasswordHash: "test",
			Version:      1,
		},
		{
			TestName:     "test_restore_user_state_is_empty",
			Expected:     ErrInvalidData,
			ID:           uuid.New(),
			Email:        "test@test.ru",
			State:        NilState,
			Status:       newUserStatus(),
			PasswordHash: "test",
			Version:      1,
		},
		{
			TestName:     "test_restore_user_status_is_empty",
			Expected:     ErrInvalidData,
			ID:           uuid.New(),
			Email:        "test@test.ru",
			State:        newActiveState(),
			Status:       NilStatus,
			PasswordHash: "test",
			Version:      1,
		},
		{
			TestName:     "test_restore_user_password_hash_id_empty",
			Expected:     ErrInvalidData,
			ID:           uuid.New(),
			Email:        "test@test.ru",
			State:        newActiveState(),
			Status:       newUserStatus(),
			PasswordHash: "",
			Version:      1,
		},
		{
			TestName:     "test_restore_user_version_is_zero",
			Expected:     ErrInvalidData,
			ID:           uuid.New(),
			Email:        "test@test.ru",
			State:        newActiveState(),
			Status:       newUserStatus(),
			PasswordHash: "test",
			Version:      0,
		},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			_, err := RestoreUser(c.ID, c.Email, c.PasswordHash, c.State, c.Status, c.Version)
			if c.Expected == nil {
				if err != nil {
					t.Errorf("expected %T, but got nil", c.Expected)
				}
			}
			if !errors.Is(err, c.Expected) {
				t.Errorf("expected %T, but got %T", c.Expected, err)
			}
		})
	}
}

func TestUser_NewEmail(t *testing.T) {
	cases := []struct {
		TestName string
		Expected error
		User     *User
		NewEmail string
	}{
		{
			TestName: "test_user_new_email_ok",
			Expected: nil,
			User:     activeUser(),
			NewEmail: "new.email@test.ru",
		},
		{
			TestName: "test_user_new_email_it_is_empty",
			Expected: ErrInvalidData,
			User:     activeUser(),
			NewEmail: "",
		},
		{
			TestName: "test_user_new_email_it_is_same",
			Expected: ErrIdempotent,
			User:     activeUser(),
			NewEmail: activeUser().Email(),
		},
		{
			TestName: "test_user_new_email_he_is_frozen",
			Expected: ErrUserNotActive,
			User:     frozenUser(),
			NewEmail: "new.email@test.ru",
		},
		{
			TestName: "test_user_new_email_he_is_deleted",
			Expected: ErrUserNotActive,
			User:     deletedUser(),
			NewEmail: "new.email@test.ru",
		},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			err := c.User.NewEmail(c.NewEmail)
			if c.Expected == nil {
				if err != nil {
					t.Errorf("expected %T, but got nil", c.Expected)
				}
			}
			if !errors.Is(err, c.Expected) {
				t.Errorf("expected %T, but got %T", c.Expected, err)
			}
		})
	}
}

func TestUser_NewState(t *testing.T) {
	cases := []struct {
		TestName string
		Expected error
		User     *User
		NewState State
	}{
		{TestName: "test_user_new_state_ok", Expected: nil, User: activeUser(), NewState: FROZEN},
		{
			TestName: "test_user_new_state_it_is_empty",
			Expected: ErrInvalidData,
			User:     activeUser(),
			NewState: NilState,
		},
		{
			TestName: "test_user_new_state_it_is_same",
			Expected: ErrIdempotent,
			User:     activeUser(),
			NewState: activeUser().State(),
		},
		{
			TestName: "test_user_new_state_he_is_frozen",
			Expected: nil,
			User:     frozenUser(),
			NewState: State(ACTIVE),
		},
		{
			TestName: "test_user_new_state_he_is_deleted",
			Expected: nil,
			User:     deletedUser(),
			NewState: State(ACTIVE),
		},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			err := c.User.NewState(c.NewState)
			if c.Expected == nil {
				if err != nil {
					t.Errorf("expected %T, but got nil", c.Expected)
				}
			}
			if !errors.Is(err, c.Expected) {
				t.Errorf("expected %T, but got %T", c.Expected, err)
			}
		})
	}
}

func TestUser_NewStatus(t *testing.T) {
	cases := []struct {
		TestName  string
		Expected  error
		User      *User
		NewStatus Status
	}{
		{TestName: "test_user_new_status_ok", Expected: nil, User: activeUser(), NewStatus: ADMIN},
		{
			TestName:  "test_user_new_status_it_is_empty",
			Expected:  ErrInvalidData,
			User:      activeUser(),
			NewStatus: NilStatus,
		},
		{
			TestName:  "test_user_new_status_it_is_same",
			Expected:  ErrIdempotent,
			User:      activeUser(),
			NewStatus: activeUser().Status(),
		},
		{
			TestName:  "test_user_new_status_he_is_frozen",
			Expected:  ErrUserNotActive,
			User:      frozenUser(),
			NewStatus: Status(ADMIN),
		},
		{
			TestName:  "test_user_new_status_he_is_deleted",
			Expected:  ErrUserNotActive,
			User:      deletedUser(),
			NewStatus: Status(ADMIN),
		},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			err := c.User.NewStatus(c.NewStatus)
			if c.Expected == nil {
				if err != nil {
					t.Errorf("expected %T, but got nil", c.Expected)
				}
			}
			if !errors.Is(err, c.Expected) {
				t.Errorf("expected %T, but got %T", c.Expected, err)
			}
		})
	}
}

func TestUser_NewPasswordHash(t *testing.T) {
	cases := []struct {
		TestName string
		Expected error
		User     *User
		NewHash  string
	}{
		{
			TestName: "test_user_new_password_hash_ok",
			Expected: nil,
			User:     activeUser(),
			NewHash:  "newPasswordHash",
		},
		{
			TestName: "test_user_new_password_hash_it_is_empty",
			Expected: ErrInvalidData,
			User:     activeUser(),
			NewHash:  "",
		},
		{
			TestName: "test_user_new_password_hash_he_is_frozen",
			Expected: ErrUserNotActive,
			User:     frozenUser(),
			NewHash:  "newPasswordHash",
		},
		{
			TestName: "test_user_new_password_hash_he_is_deleted",
			Expected: ErrUserNotActive,
			User:     deletedUser(),
			NewHash:  "newPasswordHash",
		},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			err := c.User.NewPasswordHash(c.NewHash)
			if c.Expected == nil {
				if err != nil {
					t.Errorf("expected %T, but got nil", c.Expected)
				}
			}
			if !errors.Is(err, c.Expected) {
				t.Errorf("expected %T, but got %T", c.Expected, err)
			}
		})
	}
}

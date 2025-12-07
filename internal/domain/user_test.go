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
			State:        nilState,
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
			Status:       nilStatus,
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

func activeUser() *User {
	return &User{
		id:           uuid.New(),
		email:        "test@test.ru",
		state:        newActiveState(),
		status:       newUserStatus(),
		passwordHash: "test",
		version:      1,
	}
}

func TestUser_NewEmail(t *testing.T) {
	cases := []struct {
		TestName string
		Expected error
		NewEmail string
	}{
		{TestName: "test_user_new_email_ok", Expected: nil, NewEmail: "new.email@test.ru"},
		{TestName: "test_user_new_email_it_is_empty", Expected: ErrInvalidData, NewEmail: ""},
		{
			TestName: "test_user_new_email_it_is_same",
			Expected: ErrIdempotent,
			NewEmail: activeUser().Email(),
		},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			err := activeUser().NewEmail(c.NewEmail)
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
		NewState State
	}{
		{TestName: "test_user_new_state_ok", Expected: nil, NewState: FROZEN},
		{TestName: "test_user_new_state_it_is_empty", Expected: ErrInvalidData, NewState: nilState},
		{
			TestName: "test_user_new_state_it_is_same",
			Expected: ErrIdempotent,
			NewState: activeUser().State(),
		},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			err := activeUser().NewState(c.NewState)
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
		NewStatus Status
	}{
		{TestName: "test_user_new_status_ok", Expected: nil, NewStatus: ADMIN},
		{
			TestName:  "test_user_new_status_it_is_empty",
			Expected:  ErrInvalidData,
			NewStatus: nilStatus,
		},
		{
			TestName:  "test_user_new_status_it_is_same",
			Expected:  ErrIdempotent,
			NewStatus: activeUser().Status(),
		},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			err := activeUser().NewStatus(c.NewStatus)
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
		NewHash  string
	}{
		{TestName: "test_user_new_password_hash_ok", Expected: nil, NewHash: "newPasswordHash"},
		{
			TestName: "test_user_new_password_hash_it_is_empty",
			Expected: ErrInvalidData,
			NewHash:  "",
		},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			err := activeUser().NewPasswordHash(c.NewHash)
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

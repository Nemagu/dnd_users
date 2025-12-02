package duser

import (
	"fmt"
	"testing"

	"github.com/Nemagu/dnd/internal/domain"
	"github.com/google/uuid"
)

func mustCreateEmail(email string) domain.Email {
	if e, err := domain.NewEmail(email); err != nil {
		panic(err)
	} else {
		return e
	}
}

func TestUser_New_Success(t *testing.T) {
	cases := []struct {
		UserID       uuid.UUID
		Email        domain.Email
		State        State
		Status       Status
		PasswordHash string
	}{
		{
			UserID:       uuid.New(),
			Email:        mustCreateEmail("test@example.com"),
			PasswordHash: "hashed_password",
		},
		{
			UserID:       uuid.New(),
			Email:        mustCreateEmail("test@example.com"),
			PasswordHash: "hashed_password",
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("test_create_user_%d", i), func(t *testing.T) {
			if _, err := New(
				c.UserID,
				c.Email,
				c.PasswordHash,
			); err != nil {
				t.Error(err)
			}
		})
	}
}

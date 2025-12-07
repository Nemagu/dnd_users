package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestPolicyService_CanEditOthers(t *testing.T) {
	cases := []struct {
		TestName string
		Expected bool
		User     *User
	}{
		{
			TestName: "test_policy_service_can_edit_active_admin",
			Expected: true,
			User:     activeAdmin(),
		},
		{
			TestName: "test_policy_service_can_edit_frozen_admin",
			Expected: false,
			User:     frozenAdmin(),
		},
		{
			TestName: "test_policy_service_can_edit_deleted_admin",
			Expected: false,
			User:     deletedAdmin(),
		},
		{TestName: "test_policy_service_can_edit_active_user", Expected: false, User: activeUser()},
		{TestName: "test_policy_service_can_edit_frozen_user", Expected: false, User: frozenUser()},
		{
			TestName: "test_policy_service_can_edit_deleted_user",
			Expected: false,
			User:     deletedUser(),
		},
	}
	service := MustPolicyService()
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			r := service.CanEditOthers(c.User)
			if r != c.Expected {
				t.Errorf("expected %v, but got %v", c.Expected, r)
			}
		})
	}
}

func TestPolicyService_CanReadOthers(t *testing.T) {
	cases := []struct {
		TestName string
		Expected bool
		User     *User
	}{
		{
			TestName: "test_policy_service_can_read_active_admin",
			Expected: true,
			User:     activeAdmin(),
		},
		{
			TestName: "test_policy_service_can_read_frozen_admin",
			Expected: false,
			User:     frozenAdmin(),
		},
		{
			TestName: "test_policy_service_can_read_deleted_admin",
			Expected: false,
			User:     deletedAdmin(),
		},
		{TestName: "test_policy_service_can_read_active_user", Expected: false, User: activeUser()},
		{TestName: "test_policy_service_can_read_frozen_user", Expected: false, User: frozenUser()},
		{
			TestName: "test_policy_service_can_read_deleted_user",
			Expected: false,
			User:     deletedUser(),
		},
	}
	service := MustPolicyService()
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			r := service.CanReadOthers(c.User)
			if r != c.Expected {
				t.Errorf("expected %v, but got %v", c.Expected, r)
			}
		})
	}
}

func activeAdmin() *User {
	return &User{
		id:           uuid.New(),
		email:        "test@test.ru",
		state:        State(ACTIVE),
		status:       Status(ADMIN),
		passwordHash: "test",
		version:      2,
	}
}

func frozenAdmin() *User {
	return &User{
		id:           uuid.New(),
		email:        "test@test.ru",
		state:        State(FROZEN),
		status:       Status(ADMIN),
		passwordHash: "test",
		version:      2,
	}
}

func deletedAdmin() *User {
	return &User{
		id:           uuid.New(),
		email:        "test@test.ru",
		state:        State(DELETED),
		status:       Status(ADMIN),
		passwordHash: "test",
		version:      2,
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

func frozenUser() *User {
	return &User{
		id:           uuid.New(),
		email:        "test@test.ru",
		state:        State(FROZEN),
		status:       Status(USER),
		passwordHash: "test",
		version:      2,
	}
}

func deletedUser() *User {
	return &User{
		id:           uuid.New(),
		email:        "test@test.ru",
		state:        State(DELETED),
		status:       Status(USER),
		passwordHash: "test",
		version:      2,
	}
}

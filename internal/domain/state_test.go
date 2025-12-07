package domain

import (
	"errors"
	"testing"
)

func TestState_NewState(t *testing.T) {
	cases := []struct {
		TestName  string
		StateName string
		Expected  error
	}{
		{TestName: "test_new_active_state", StateName: ACTIVE, Expected: nil},
		{TestName: "test_new_frozen_state", StateName: FROZEN, Expected: nil},
		{TestName: "test_new_deleted_state", StateName: DELETED, Expected: nil},
		{TestName: "test_new_other_state", StateName: "other", Expected: ErrInvalidData},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			_, err := NewState(c.StateName)
			if !errors.Is(err, c.Expected) {
				t.Errorf("expected %T, but got %T", c.Expected, err)
			}
		})
	}
}

package domain

import (
	"errors"
	"testing"
)

func TestStatus_NewStatus(t *testing.T) {
	cases := []struct {
		TestName   string
		StatusName string
		Expected   error
	}{
		{TestName: "test_new_admin_status", StatusName: ADMIN, Expected: nil},
		{TestName: "test_new_user_status", StatusName: USER, Expected: nil},
		{TestName: "test_new_other_status", StatusName: "other", Expected: ErrInvalidData},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			_, err := NewStatus(c.StatusName)
			if !errors.Is(err, c.Expected) {
				t.Errorf("expected %T, but got %T", c.Expected, err)
			}
		})
	}
}

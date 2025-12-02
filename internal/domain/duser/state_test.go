package duser

import (
	"fmt"
	"testing"
)

func TestState_StateFromString_Success(t *testing.T) {
	cases := []string{
		ACTIVE,
		FROZEN,
		DELETED,
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("test_%s_state", c), func(t *testing.T) {
			if _, err := StateFromString(c); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestState_StateFromString_Bad(t *testing.T) {
	cases := []string{
		"activated",
		"freeze",
		"delete",
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("test_%s_state", c), func(t *testing.T) {
			if _, err := StateFromString(c); err == nil {
				t.Errorf("state from string %s was not raise error", c)
			}
		})
	}
}

package duser

import (
	"fmt"
	"testing"
)

func TestStatus_StatusFromString_Success(t *testing.T) {
	cases := []string{
		ADMIN,
		ORDINARY,
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("test_%s_status", c), func(t *testing.T) {
			if _, err := StatusFromString(c); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestStatus_StatusFromString_Bad(t *testing.T) {
	cases := []string{
		"normal",
		"administration",
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("test_%s_status", c), func(t *testing.T) {
			if _, err := StatusFromString(c); err == nil {
				t.Errorf("status from string %s was not raise error", c)
			}
		})
	}
}

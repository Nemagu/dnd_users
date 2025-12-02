package domain

import (
	"fmt"
	"testing"
)

func TestEmail_NewEmail_Success(t *testing.T) {
	cases := []string{"my.email@mail.com", "ok@mail.com", "kek@lol.com"}
	for _, c := range cases {
		t.Run(fmt.Sprintf("test_good_email_%s", c), func(t *testing.T) {
			if _, err := NewEmail(c); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestEmail_NewEmail_Bad(t *testing.T) {
	cases := []string{"", "@ail", "kek"}
	for _, c := range cases {
		t.Run(fmt.Sprintf("test_bad_email_%s", c), func(t *testing.T) {
			if _, err := NewEmail(c); err == nil {
				t.Errorf("email %s was not raise err", c)
			}
		})
	}
}

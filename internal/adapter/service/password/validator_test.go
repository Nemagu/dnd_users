package password

import (
	"fmt"
	"testing"
)

func initValidator() *PasswordValidator {
	return MustNewPasswordValidator(
		8, 72, true, true, true, true,
	)
}

func TestValidator_Validate_Success(t *testing.T) {
	cases := []struct {
		Password string
		Email    string
	}{
		{
			Password: "pDK&dkg18dN",
			Email:    "email@mail.com",
		},
		{
			Password: "jJ#3sdkfju",
			Email:    "email@mail.com",
		},
		{
			Password: "dsakfljk;DJads;gl*h977J:m;adjj",
			Email:    "email@mail.com",
		},
	}
	validator := initValidator()
	for _, c := range cases {
		t.Run(fmt.Sprintf("test_validate_password_%s", c.Password), func(t *testing.T) {
			if err := validator.Validate(c.Password, c.Email); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestValidator_Validate_Bad(t *testing.T) {
	cases := []struct {
		Password string
		Email    string
	}{
		{
			Password: "",
			Email:    "email@mail.com",
		},
		{
			Password: "hjkl;sadflj",
			Email:    "email@mail.com",
		},
		{
			Password: "email@mail.com",
			Email:    "email@mail.com",
		},
		{
			Password: "dafsjj;lkKLJdkjf",
			Email:    "email@mail.com",
		},
		{
			Password: "kadslfuohenJDK788",
			Email:    "email@mail.com",
		},
		{
			Password: "SKDLJFLDFNDV:KDS:L8",
			Email:    "email@mail.com",
		},
		{
			Password: "password",
			Email:    "email@mail.com",
		},
	}
	validator := initValidator()
	for _, c := range cases {
		t.Run(fmt.Sprintf("test_validate_password_%s", c.Password), func(t *testing.T) {
			if err := validator.Validate(c.Password, c.Email); err == nil {
				t.Errorf("password %s did not raise error", c.Password)
			}
		})
	}
}

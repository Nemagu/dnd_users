package email

import (
	"fmt"
	"regexp"

	"github.com/Nemagu/dnd/internal/application"
)

type EmailValidator struct {
	emailPattern *regexp.Regexp
}

func MustNewEmailValidator() *EmailValidator {
	return &EmailValidator{
		emailPattern: regexp.MustCompile(`^[a-zA-Z0-9]{1}[a-zA-Z0-9._%+-].*@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
	}
}

func (v *EmailValidator) Validate(email string) error {
	if len(email) < 5 || len(email) > 254 {
		return fmt.Errorf(
			"%w: длина email должна находиться в диапазоне от 5 до 254",
			application.ErrValidation,
		)
	}

	matched := v.emailPattern.MatchString(email)
	if !matched {
		return fmt.Errorf(
			"%w: не корректный email",
			application.ErrValidation,
		)
	}

	return nil
}

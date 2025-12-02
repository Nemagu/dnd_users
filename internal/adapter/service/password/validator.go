package password

import (
	"fmt"
	"strings"

	"github.com/Nemagu/dnd/internal/application"
)

type PasswordValidator struct {
	minLength      int
	maxLength      int
	requireUpper   bool
	requireLower   bool
	requireNumber  bool
	requireSpecial bool
}

func NewPasswordValidator(
	minLength, maxLength int,
	requireUpper, requireLower, requireNumber, requireSpecial bool,
) (*PasswordValidator, error) {
	return &PasswordValidator{
		minLength:      minLength,
		maxLength:      maxLength,
		requireUpper:   requireUpper,
		requireLower:   requireLower,
		requireNumber:  requireNumber,
		requireSpecial: requireSpecial,
	}, nil
}

func MustNewPasswordValidator(
	minLength, maxLength int,
	requireUpper, requireLower, requireNumber, requireSpecial bool,
) *PasswordValidator {
	return &PasswordValidator{
		minLength:      minLength,
		maxLength:      maxLength,
		requireUpper:   requireUpper,
		requireLower:   requireLower,
		requireNumber:  requireNumber,
		requireSpecial: requireSpecial,
	}
}

func (v *PasswordValidator) Validate(password, email string) error {
	if err := v.validateLength(password); err != nil {
		return err
	}

	if err := v.validateSimilarToEmail(password, email); err != nil {
		return err
	}

	if err := v.validateCommon(password); err != nil {
		return err
	}

	if err := v.validateCharacters(password); err != nil {
		return err
	}

	return nil
}

func (v *PasswordValidator) validateLength(password string) error {
	if len(password) < v.minLength {
		return fmt.Errorf(
			"%w: пароль слишком короткий",
			application.ErrValidation,
		)
	}

	if len(password) > v.maxLength {
		return fmt.Errorf(
			"%w: пароль слишком длинный",
			application.ErrValidation,
		)
	}

	return nil
}

func (v *PasswordValidator) validateSimilarToEmail(password, email string) error {
	if email == "" {
		return nil
	}

	emailLocalPart := strings.Split(email, "@")[0]
	missingRules := make([]string, 0, 2)

	if strings.EqualFold(emailLocalPart, password) {
		missingRules = append(missingRules, "совпадать с email")
	}

	if len(emailLocalPart) > 4 && strings.Contains(
		strings.ToLower(password),
		strings.ToLower(emailLocalPart),
	) {
		missingRules = append(missingRules, "содержать часть email")
	}

	if len(missingRules) > 0 {
		return fmt.Errorf(
			"%w: пароль не должен %s",
			application.ErrValidation,
			strings.Join(missingRules, ", "),
		)
	}

	return nil
}

func (v *PasswordValidator) validateCommon(password string) error {
	commonPasswords := map[string]bool{
		"password":   true,
		"123456":     true,
		"12345678":   true,
		"qwerty":     true,
		"admin":      true,
		"welcome":    true,
		"password1":  true,
		"123456789":  true,
		"12345":      true,
		"1234567890": true,
	}

	if commonPasswords[strings.ToLower(password)] {
		return fmt.Errorf(
			"%w: пароль слишком простой и распространенный",
			application.ErrValidation,
		)
	}

	return nil
}

func (v *PasswordValidator) validateCharacters(password string) error {
	missingRules := make([]string, 0, 4)
	if v.requireSpecial {
		if !strings.ContainsAny(password, "!@#$%^&*()_+-=[]{}|;':\",.<>/?") {
			missingRules = append(missingRules, "специальные символы")
		}
	}
	if v.requireNumber {
		if !strings.ContainsAny(password, "0123456789") {
			missingRules = append(missingRules, "цифры")
		}
	}
	if v.requireUpper {
		if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
			missingRules = append(missingRules, "заглавные буквы")
		}
	}
	if v.requireLower {
		if !strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") {
			missingRules = append(missingRules, "строчные буквы")
		}
	}
	if len(missingRules) > 0 {
		return fmt.Errorf(
			"%w: пароль должен содержать %s",
			application.ErrValidation,
			strings.Join(missingRules, ", "),
		)
	}
	return nil
}

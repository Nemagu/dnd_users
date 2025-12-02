package email

import "testing"

func createValidator() *EmailValidator {
	return MustNewEmailValidator()
}

func TestEmailValidator_Validate_Success(t *testing.T) {
	cases := []string{
		"normal@mail.com",
		"kek@mail.com",
		"lol@mail.com",
		"ohohoh@mail.com",
	}
	validator := createValidator()
	for _, c := range cases {
		t.Run("test_email_validate_"+c, func(t *testing.T) {
			if err := validator.Validate(c); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestEmailValidator_Validate_Bad(t *testing.T) {
	cases := []string{
		"normalmail.com",
		"kek@mailcom",
		"lol@mail",
		"ohohoh@mail.",
		".kek@mail.com",
	}
	validator := createValidator()
	for _, c := range cases {
		t.Run("test_email_validate_"+c, func(t *testing.T) {
			if err := validator.Validate(c); err == nil {
				t.Errorf("email %s was not raise err", c)
			}
		})
	}
}

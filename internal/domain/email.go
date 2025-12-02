package domain

import "fmt"

type Email struct {
	email string
}

func NewEmail(email string) (Email, error) {
	if len(email) < 5 {
		return Email{}, fmt.Errorf("%w: email слишком короткий", ErrValidation)
	}
	return Email{email: email}, nil
}

func (e Email) String() string {
	return e.email
}

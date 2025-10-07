package domain

type Email struct {
	email string
}

func NewEmail(email string) (Email, error) {
	// TODO: add validation email
	if len(email) < 5 {
		return Email{}, &DomainError{Message: "email слишком короткий"}
	}
	return Email{email: email}, nil
}

func (e Email) Email() string {
	return e.email
}

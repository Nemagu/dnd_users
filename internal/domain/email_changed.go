package domain

import (
	"fmt"
	"time"
)

type EmailChanged struct {
	userID     UserID
	email      Email
	occurredOn time.Time
}

func NewEmailChanged(userID UserID, email Email) *EmailChanged {
	return &EmailChanged{
		userID:     userID,
		email:      email,
		occurredOn: time.Now().UTC(),
	}
}

func NewEmailChangedFromUser(user *User) *EmailChanged {
	return &EmailChanged{
		userID:     user.UserID(),
		email:      user.Email(),
		occurredOn: time.Now().UTC(),
	}
}

func (e *EmailChanged) UserID() UserID {
	return e.userID
}

func (e *EmailChanged) Email() Email {
	return e.email
}

func (e *EmailChanged) OccurredOn() time.Time {
	return e.occurredOn
}

func (e *EmailChanged) String() string {
	return fmt.Sprintf("user id: %s\nemail: %s\noccurred on: %s", e.userID, e.email, e.occurredOn)
}

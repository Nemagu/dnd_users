package domain

import (
	"fmt"
	"time"
)

type StatusChanged struct {
	userID     UserID
	status     UserStatus
	occurredOn time.Time
}

func NewStatusChanged(userID UserID, status UserStatus) *StatusChanged {
	return &StatusChanged{
		userID:     userID,
		status:     status,
		occurredOn: time.Now().UTC(),
	}
}

func NewStatusChangedFromUser(user *User) *StatusChanged {
	return &StatusChanged{
		userID:     user.UserID(),
		status:     user.Status(),
		occurredOn: time.Now().UTC(),
	}
}

func (e *StatusChanged) UserID() UserID {
	return e.userID
}

func (e *StatusChanged) Status() UserStatus {
	return e.status
}

func (e *StatusChanged) OccurredOn() time.Time {
	return e.occurredOn
}

func (e *StatusChanged) String() string {
	return fmt.Sprintf("user id: %s\nstatus: %s\noccurred on: %s", e.userID, e.status, e.occurredOn)
}

package domain

import (
	"fmt"
	"time"
)

type StateChanged struct {
	userID     UserID
	state      UserState
	occurredOn time.Time
}

func NewStateChanged(userID UserID, state UserState) *StateChanged {
	return &StateChanged{
		userID:     userID,
		state:      state,
		occurredOn: time.Now().UTC(),
	}
}

func NewStateChangedFromUser(user *User) *StateChanged {
	return &StateChanged{
		userID:     user.UserID(),
		state:      user.State(),
		occurredOn: time.Now().UTC(),
	}
}

func (e *StateChanged) UserID() UserID {
	return e.userID
}

func (e *StateChanged) State() UserState {
	return e.state
}

func (e *StateChanged) OccurredOn() time.Time {
	return e.occurredOn
}

func (e *StateChanged) String() string {
	return fmt.Sprintf("user id: %s\nstate: %s\noccurred on: %s", e.userID, e.state, e.occurredOn)
}

package domain

import "fmt"

const (
	ACTIVE  = "active"
	FROZEN  = "frozen"
	DELETED = "deleted"
)

var nilState = State("")

type State string

func NewState(state string) (State, error) {
	switch state {
	case ACTIVE:
		return ACTIVE, nil
	case FROZEN:
		return FROZEN, nil
	case DELETED:
		return DELETED, nil
	default:
		return "", fmt.Errorf(
			"%w: состояния пользователя с названием %s не существует",
			ErrInvalidData,
			state,
		)
	}
}

func newActiveState() State {
	return State(ACTIVE)
}

func (s State) IsActive() bool {
	return s == ACTIVE
}

func (s State) IsFrozen() bool {
	return s == FROZEN
}

func (s State) IsDeleted() bool {
	return s == DELETED
}

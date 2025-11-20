package duser

import (
	"fmt"
	"strings"

	"github.com/Nemagu/dnd/internal/domain"
)

const (
	ACTIVE  = "active"
	FROZEN  = "frozen"
	DELETED = "deleted"
)

type UserState string

func NewActiveState() UserState {
	return ACTIVE
}

func NewFrozenState() UserState {
	return FROZEN
}

func NewDeletedState() UserState {
	return DELETED
}

func StateFromString(s string) (UserState, error) {
	switch strings.ToLower(s) {
	case ACTIVE:
		return ACTIVE, nil
	case FROZEN:
		return FROZEN, nil
	case DELETED:
		return DELETED, nil
	default:
		return "", fmt.Errorf("%w: %s", domain.ErrInvalidData, s)
	}
}

func (us UserState) IsActive() bool {
	return us == ACTIVE
}

func (us UserState) IsFrozen() bool {
	return us == FROZEN
}

func (us UserState) IsDeleted() bool {
	return us == DELETED
}

func (us UserState) String() string {
	return string(us)
}

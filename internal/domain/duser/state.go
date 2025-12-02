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

type State string

func NewActiveState() State {
	return ACTIVE
}

func NewFrozenState() State {
	return FROZEN
}

func NewDeletedState() State {
	return DELETED
}

func StateFromString(s string) (State, error) {
	switch strings.ToLower(s) {
	case ACTIVE:
		return ACTIVE, nil
	case FROZEN:
		return FROZEN, nil
	case DELETED:
		return DELETED, nil
	default:
		return "", fmt.Errorf(
			"%w: состояния %s не существует",
			domain.ErrInvalidData,
			s,
		)
	}
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

func (s State) String() string {
	return string(s)
}

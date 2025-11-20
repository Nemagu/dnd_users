package duser

import (
	"fmt"
	"strings"

	"github.com/Nemagu/dnd/internal/domain"
)

const (
	ADMIN    = "admin"
	ORDINARY = "ordinary"
)

type UserStatus string

func NewAdminStatus() UserStatus {
	return ADMIN
}

func NewOrdinaryStatus() UserStatus {
	return ORDINARY
}

func StatusFromString(s string) (UserStatus, error) {
	switch strings.ToLower(s) {
	case ADMIN:
		return ADMIN, nil
	case ORDINARY:
		return ORDINARY, nil
	default:
		return "", fmt.Errorf("%w: %s", domain.ErrInvalidData, s)
	}
}

func (us UserStatus) IsAdmin() bool {
	return us == ADMIN
}

func (us UserStatus) IsOrdinary() bool {
	return us == ORDINARY
}

func (us UserStatus) String() string {
	return string(us)
}

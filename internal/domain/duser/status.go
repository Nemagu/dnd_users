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

type Status string

func NewAdminStatus() Status {
	return ADMIN
}

func NewOrdinaryStatus() Status {
	return ORDINARY
}

func StatusFromString(s string) (Status, error) {
	switch strings.ToLower(s) {
	case ADMIN:
		return ADMIN, nil
	case ORDINARY:
		return ORDINARY, nil
	default:
		return "", fmt.Errorf(
			"%w: статуса %s не существует",
			domain.ErrInvalidData,
			s,
		)
	}
}

func (s Status) IsAdmin() bool {
	return s == ADMIN
}

func (s Status) IsOrdinary() bool {
	return s == ORDINARY
}

func (s Status) String() string {
	return string(s)
}

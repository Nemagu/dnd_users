package domain

import "fmt"

const (
	ADMIN = "admin"
	USER  = "user"
)

var NilStatus = Status("")

type Status string

func NewStatus(status string) (Status, error) {
	switch status {
	case ADMIN:
		return ADMIN, nil
	case USER:
		return USER, nil
	default:
		return "", fmt.Errorf(
			"%w: статуса пользователя с названием %s не существует",
			ErrInvalidData,
			status,
		)
	}
}

func newUserStatus() Status {
	return Status(USER)
}

func (s Status) String() string {
	return string(s)
}

func (s Status) IsAdmin() bool {
	return s == ADMIN
}

func (s Status) IsUser() bool {
	return s == USER
}

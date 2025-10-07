package domain

const (
	ACTIVE = iota
	PENDING
	REMOVED
)

type UserState uint

func NewUserState(state uint) (UserState, error) {
	if state > 2 {
		return UserState(0), &DomainError{Message: "переданного статуса не существует"}
	}
	return UserState(state), nil
}

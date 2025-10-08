package domain

const (
	ACTIVE = iota
	PENDING
	REMOVED
)

type UserState uint

func NewActiveUserState() UserState {
	return UserState(ACTIVE)
}

func NewPendingUserState() UserState {
	return UserState(PENDING)
}

func NewRemovedUserState() UserState {
	return UserState(REMOVED)
}

func (us UserState) IsActive() bool {
	return us == ACTIVE
}

func (us UserState) IsPending() bool {
	return us == PENDING
}

func (us UserState) IsRemoved() bool {
	return us == REMOVED
}

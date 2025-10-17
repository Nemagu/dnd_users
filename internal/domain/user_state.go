package domain

const (
	ACTIVE = iota
	FROZEN
	DELETED
)

type UserState uint

func NewActiveUserState() UserState {
	return ACTIVE
}

func NewFrozenUserState() UserState {
	return FROZEN
}

func NewDeletedUserState() UserState {
	return DELETED
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
	switch us {
	case ACTIVE:
		return "active"
	case FROZEN:
		return "frozen"
	case DELETED:
		return "deleted"
	default:
		return "unknow state"
	}
}

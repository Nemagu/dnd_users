package domain

const (
	ADMIN = iota
	ORDINARY
)

type UserStatus uint

func NewAdminUserStatus() UserStatus {
	return ADMIN
}

func NewOrdinaryUserStatus() UserStatus {
	return ORDINARY
}

func (us UserStatus) IsAdmin() bool {
	return us == ADMIN
}

func (us UserStatus) IsOrdinary() bool {
	return us == ORDINARY
}

func (us UserStatus) String() string {
	switch us {
	case ADMIN:
		return "admin"
	case ORDINARY:
		return "ordinary"
	default:
		return "unknow status"
	}
}

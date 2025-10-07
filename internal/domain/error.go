package domain

type DomainError struct {
	Message string
}

func (de *DomainError) Error() string {
	return de.Message
}

package domain

import (
	"fmt"
)

const (
	NOT_FOUND = iota
	NO_ACCESS
	INVALID_DATA
	IDEMPOTENT
	POLICY
	INTERNAL
)

type ErrorCode uint

func (ec ErrorCode) String() string {
	switch ec {
	case NOT_FOUND:
		return "not found"
	case NO_ACCESS:
		return "no access"
	case INVALID_DATA:
		return "invalid data"
	case IDEMPOTENT:
		return "idempotent error"
	case POLICY:
		return "policy error"
	case INTERNAL:
		return "internal error"
	default:
		return "unknow error"
	}
}

type DomainError struct {
	Code    ErrorCode
	Message string
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("Error: %s\nMessage: %s", e.Code, e.Message)
}

func NotFoundError(message string) *DomainError {
	if len(message) == 0 {
		message = "объект не найден"
	}
	return &DomainError{
		Code:    NOT_FOUND,
		Message: message,
	}
}

func NoAccessError(message string) *DomainError {
	if len(message) == 0 {
		message = "у вас недостаточно прав для выполнения операции"
	}
	return &DomainError{
		Code:    NO_ACCESS,
		Message: message,
	}
}

func InvalidDataError(message string) *DomainError {
	return &DomainError{
		Code:    INVALID_DATA,
		Message: message,
	}
}

func IdempotentError(message string) *DomainError {
	return &DomainError{
		Code:    IDEMPOTENT,
		Message: message,
	}
}

func PolicyError(message string) *DomainError {
	return &DomainError{
		Code:    POLICY,
		Message: message,
	}
}

func InternalError(message string) *DomainError {
	return &DomainError{
		Code:    INTERNAL,
		Message: message,
	}
}

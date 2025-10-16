package application

import (
	"fmt"
)

const (
	NOT_FOUND = iota
	NO_ACCESS
	INVALID_DATA
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
	case INTERNAL:
		return "internal"
	default:
		return "unknow error"
	}
}

type ApplicationError struct {
	Code    ErrorCode
	Message string
}

func (e *ApplicationError) Error() string {
	return fmt.Sprintf("Error: %s\nMessage: %s", e.Code, e.Message)
}

func NotFoundError(message string) *ApplicationError {
	if len(message) == 0 {
		message = "объект не найден"
	}
	return &ApplicationError{
		Code:    NOT_FOUND,
		Message: message,
	}
}

func NoAccessError(message string) *ApplicationError {
	if len(message) == 0 {
		message = "у вас недостаточно прав для выполнения операции"
	}
	return &ApplicationError{
		Code:    NO_ACCESS,
		Message: message,
	}
}

func InvalidDataError(message string) *ApplicationError {
	return &ApplicationError{
		Code:    INVALID_DATA,
		Message: message,
	}
}

func InternalError(message string) *ApplicationError {
	return &ApplicationError{
		Code:    INTERNAL,
		Message: message,
	}
}

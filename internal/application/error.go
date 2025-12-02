package application

import "errors"

var (
	ErrValidation      = errors.New("ошибка валидации")
	ErrAlreadyExists   = errors.New("объект уже существует")
	ErrVersionConflict = errors.New("конфликт версий")
	ErrCredential      = errors.New("ошибка авторизации")
	ErrNotAllowed      = errors.New("действие не разрешено")
	ErrNotFound        = errors.New("объект не найден")
	ErrInternal        = errors.New("внутренняя ошибка")
)

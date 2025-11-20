package application

import "errors"

var (
	ErrNotFound        = errors.New("объект не найден")
	ErrValidation      = errors.New("ошибка валидации")
	ErrInternal        = errors.New("внутренняя ошибка")
	ErrAlreadyExists   = errors.New("объект уже существует")
	ErrNotAllowed      = errors.New("действие не разрешено")
	ErrVersionConflict = errors.New("конфликт версий")
)

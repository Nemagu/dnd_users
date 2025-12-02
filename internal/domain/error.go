package domain

import "errors"

var (
	ErrIdempotent  = errors.New("данные не изменяются")
	ErrValidation  = errors.New("ошибка валидации данных")
	ErrNotAllowed  = errors.New("ошибка политики")
	ErrInvalidData = errors.New("переданы не корректные данные")
	ErrInternal    = errors.New("внутренняя ошибка")
)

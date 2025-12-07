package domain

import "errors"

var (
	ErrInvalidData   = errors.New("не корректные данные")
	ErrIdempotent    = errors.New("попытка изменения пользователя без изменения данных")
	ErrUserNotActive = errors.New("пользователь имеет не активный статус")
)

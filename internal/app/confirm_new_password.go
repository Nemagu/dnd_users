package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ConfirmNewPasswordUseCase struct {
	repo          confirmNewPasswordRepository
	store         confirmNewPasswordCodeStore
	emailProvider confirmNewPasswordEmailProvider
	codeGenerator codeGenerator
}

type ConfirmNewPasswordCommand struct {
	InitiatorID uuid.UUID
	UserID      uuid.UUID
}

type confirmNewPasswordRepository interface {
	ByID(ctx context.Context, id uuid.UUID) (*User, error)
}

type confirmNewPasswordCodeStore interface {
	SetNewPassword(ctx context.Context, key, value string) error
}

type confirmNewPasswordEmailProvider interface {
	SendConfirmationNewPassword(data EmailCode)
}

func MustConfirmNewPasswordUseCase(
	repo confirmNewPasswordRepository,
	store confirmNewPasswordCodeStore,
	emailProvider confirmNewPasswordEmailProvider,
	codeGenerator codeGenerator,
) *ConfirmNewPasswordUseCase {
	if repo == nil {
		panic("confirm new password use case did not get user repository")
	}
	if store == nil {
		panic("confirm new password use case did not get code store")
	}
	if emailProvider == nil {
		panic("confirm new password use case did not get email provider")
	}
	if codeGenerator == nil {
		panic("confirm new password use case did not get code generator")
	}
	return &ConfirmNewPasswordUseCase{
		repo:          repo,
		store:         store,
		emailProvider: emailProvider,
		codeGenerator: codeGenerator,
	}
}

func (u *ConfirmNewPasswordUseCase) Execute(
	ctx context.Context,
	command *ConfirmNewPasswordCommand,
) error {
	if command.InitiatorID != command.UserID {
		return fmt.Errorf("%w: вы не можете изменять пароль другим пользователям", ErrNotAllowed)
	}

	user, err := u.repo.ByID(ctx, command.UserID)
	if err != nil {
		return err
	}

	code := u.codeGenerator.Generate()

	if err = u.store.SetNewPassword(ctx, user.ID.String()+user.Email, code); err != nil {
		return err
	}

	go u.emailProvider.SendConfirmationNewPassword(EmailCode{To: user.Email, Code: code})

	return nil
}

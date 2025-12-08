package app

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type confirmNewEmailRepository interface {
	ByID(ctx context.Context, id uuid.UUID) (*User, error)
	EmailExists(ctx context.Context, email string) (bool, error)
}

type confirmNewEmailCodeStore interface {
	SetNewEmail(ctx context.Context, key, value string) error
}

type confirmNewEmailProvider interface {
	SendConfirmationNewEmail(data []EmailCode)
}

type ConfirmNewEmailCommand struct {
	InitiatorID uuid.UUID
	UserID      uuid.UUID
	NewEmail    string
}

type ConfirmNewEmailUseCase struct {
	repo          confirmNewEmailRepository
	store         confirmNewEmailCodeStore
	validator     emailValidator
	emailProvider confirmNewEmailProvider
	codeGenerator codeGenerator
}

func MustConfirmNewEmailUseCase(
	repo confirmNewEmailRepository,
	store confirmNewEmailCodeStore,
	validator emailValidator,
	emailProvider confirmNewEmailProvider,
	codeGenerator codeGenerator,
) *ConfirmNewEmailUseCase {
	if repo == nil {
		panic("confirm new email use case did not get user repository")
	}
	if store == nil {
		panic("confirm new email use case did not get code store")
	}
	if validator == nil {
		panic("confirm new email use case did not get email validator")
	}
	if emailProvider == nil {
		panic("confirm new email use case did not get email provider")
	}
	if codeGenerator == nil {
		panic("confirm new email use case did not get code generator")
	}
	return &ConfirmNewEmailUseCase{
		repo:          repo,
		store:         store,
		validator:     validator,
		emailProvider: emailProvider,
		codeGenerator: codeGenerator,
	}
}

func (u *ConfirmNewEmailUseCase) Execute(
	ctx context.Context,
	command *ConfirmNewEmailCommand,
) error {
	if command.InitiatorID != command.UserID {
		return fmt.Errorf("%w: вы не можете изменять email другим пользователям", ErrNotAllowed)
	}

	if err := u.validator.Validate(command.NewEmail); err != nil {
		return err
	}

	exists, err := u.repo.EmailExists(ctx, command.NewEmail)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%w: email %s уже существует", ErrAlreadyExists, command.NewEmail)
	}

	user, err := u.repo.ByID(ctx, command.UserID)
	if err != nil {
		return err
	}

	oldEmailCode := u.codeGenerator.Generate()
	newEmailCode := u.codeGenerator.Generate()

	if err = u.store.SetNewEmail(ctx, user.ID.String()+user.Email, oldEmailCode); err != nil {
		return err
	}
	if err = u.store.SetNewEmail(ctx, user.ID.String()+command.NewEmail, newEmailCode); err != nil {
		return err
	}

	go u.emailProvider.SendConfirmationNewEmail(
		[]EmailCode{
			{To: user.Email, Code: oldEmailCode},
			{To: command.NewEmail, Code: newEmailCode},
		},
	)

	return nil
}

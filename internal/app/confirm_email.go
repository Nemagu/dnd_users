package app

import (
	"context"
	"fmt"
)

type ConfirmEmailCommand struct {
	Email string
}

type ConfirmEmailUseCase struct {
	repo          confirmEmailRepository
	store         confirmEmailCodeStore
	validator     emailValidator
	emailProvider confirmEmailProvider
	codeGenerator codeGenerator
}

type confirmEmailRepository interface {
	EmailExists(ctx context.Context, email string) (bool, error)
}

type confirmEmailCodeStore interface {
	SetCEC(ctx context.Context, key, value string) error
}

type confirmEmailProvider interface {
	SendConfirmationEmail(data EmailCode)
}

func MustConfirmEmailUseCase(
	repo confirmEmailRepository,
	store confirmEmailCodeStore,
	validator emailValidator,
	emailProvider confirmEmailProvider,
	codeGenerator codeGenerator,
) *ConfirmEmailUseCase {
	if repo == nil {
		panic("confirm email use case did not get user repository")
	}
	if store == nil {
		panic("confirm email use case did not get key value store")
	}
	if validator == nil {
		panic("confirm email use case did not get email validator")
	}
	if emailProvider == nil {
		panic("confirm email use case did not get email provider")
	}
	if codeGenerator == nil {
		panic("confirm email use case did not get code generator")
	}
	return &ConfirmEmailUseCase{
		repo:          repo,
		store:         store,
		validator:     validator,
		emailProvider: emailProvider,
		codeGenerator: codeGenerator,
	}
}

func (u *ConfirmEmailUseCase) Execute(ctx context.Context, command *ConfirmEmailCommand) error {
	if err := u.validator.Validate(command.Email); err != nil {
		return err
	}

	exists, err := u.repo.EmailExists(ctx, command.Email)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%w: email %s уже существует", ErrAlreadyExists, command.Email)
	}

	code := u.codeGenerator.Generate()
	if err = u.store.SetCEC(ctx, command.Email, code); err != nil {
		return err
	}

	go u.emailProvider.SendConfirmationEmail(EmailCode{To: command.Email, Code: code})

	return nil
}

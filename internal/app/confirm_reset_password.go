package app

import (
	"context"
	"fmt"
)

type ConfirmResetPasswordUseCase struct {
	repo          confirmResetPasswordRepository
	store         confirmResetPasswordCodeStore
	emailProvider confirmResetPasswordProvider
	codeGenerator codeGenerator
}

type ConfirmResetPasswordCommand struct {
	Email string
}

type confirmResetPasswordRepository interface {
	EmailExists(ctx context.Context, email string) (bool, error)
}

type confirmResetPasswordCodeStore interface {
	SetResetPassword(ctx context.Context, key, value string) error
}

type confirmResetPasswordProvider interface {
	SendResetPasswordEmail(data EmailCode)
}

func MustConfirmResetPasswordUseCase(
	repo confirmResetPasswordRepository,
	store confirmResetPasswordCodeStore,
	emailProvider confirmResetPasswordProvider,
	codeGenerator codeGenerator,
) *ConfirmResetPasswordUseCase {
	if repo == nil {
		panic("confirm reset password did not get user repository")
	}
	if store == nil {
		panic("confirm reset password did not get code store")
	}
	if emailProvider == nil {
		panic("confirm reset password did not get email provider")
	}
	if codeGenerator == nil {
		panic("confirm reset password did not get code generator")
	}
	return &ConfirmResetPasswordUseCase{
		repo:          repo,
		store:         store,
		emailProvider: emailProvider,
		codeGenerator: codeGenerator,
	}
}

func (u *ConfirmResetPasswordUseCase) Execute(
	ctx context.Context,
	command *ConfirmResetPasswordCommand,
) error {
	exists, err := u.repo.EmailExists(ctx, command.Email)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%w: такого email не существует", ErrNotFound)
	}

	code := u.codeGenerator.Generate()
	if err = u.store.SetResetPassword(ctx, command.Email, code); err != nil {
		return err
	}

	go u.emailProvider.SendResetPasswordEmail(EmailCode{To: command.Email, Code: code})

	return nil
}
